package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级为 WebSocket
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade:", err)
		return
	}
	defer wsConn.Close()

	// 配置 MediaEngine 和 API
	m := webrtc.MediaEngine{}
	m.RegisterDefaultCodecs()
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	peerConn, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatal("Create PeerConnection:", err)
		return
	}
	defer func() {
		if err := peerConn.Close(); err != nil {
			log.Println("PeerConnection close error:", err)
		}
	}()

	// 接收音频轨道转发器
	_, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio)
	if err != nil {
		log.Fatal("AddTransceiver error:", err)
		return
	}

	// 创建发送轨道 (发送给前端)
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus},
		"audio", "pion",
	)
	if err != nil {
		log.Fatal("Create local audio track error:", err)
		return
	}

	_, err = peerConn.AddTrack(audioTrack)
	if err != nil {
		log.Fatal("AddTrack error:", err)
		return
	}

	// 监听远端音频轨道，接收 RTP 并写入本地轨道转发
	peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Println("✅ Got audio track:", track.Codec().MimeType)

		go func() {
			for {
				pkt, _, err := track.ReadRTP()
				if err != nil {
					log.Println("Track read error:", err)
					return
				}
				log.Printf("📦 Received RTP packet: SSRC=%d Seq=%d TS=%d Size=%d\n",
					pkt.SSRC, pkt.SequenceNumber, pkt.Timestamp, len(pkt.Payload))

				// 直接写 RTP 包转发给前端
				if err := audioTrack.WriteRTP(pkt); err != nil {
					log.Println("WriteRTP error:", err)
					return
				}
			}
		}()
	})

	// 发送 ICE candidate 给客户端
	peerConn.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		cand := map[string]interface{}{"candidate": c.ToJSON()}
		msg, err := json.Marshal(cand)
		if err != nil {
			log.Println("ICE candidate marshal error:", err)
			return
		}
		if err = wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("WriteMessage ICE candidate error:", err)
		}
	})

	// 读取客户端 Offer
	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		log.Println("ReadMessage error:", err)
		return
	}

	var offer webrtc.SessionDescription
	if err = json.Unmarshal(msg, &offer); err != nil {
		log.Println("Unmarshal offer error:", err)
		return
	}

	// 设置远端描述
	if err = peerConn.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription error:", err)
		return
	}

	// 创建 Answer
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer error:", err)
		return
	}

	// 设置本地描述
	if err = peerConn.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription error:", err)
		return
	}

	// 发送 Answer 给客户端
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		log.Println("Marshal answer error:", err)
		return
	}
	if err = wsConn.WriteMessage(websocket.TextMessage, answerJSON); err != nil {
		log.Println("WriteMessage answer error:", err)
		return
	}

	// 循环读取客户端 ICE candidate
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("WebSocket closed:", err)
			break
		}

		var ice webrtc.ICECandidateInit
		if err := json.Unmarshal(msg, &ice); err == nil && ice.Candidate != "" {
			if err := peerConn.AddICECandidate(ice); err != nil {
				log.Println("AddICECandidate error:", err)
			}
		}
	}
}
