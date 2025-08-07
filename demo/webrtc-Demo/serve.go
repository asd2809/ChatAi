package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)
// 升级为ws
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	r := gin.Default()

	// 先注册具体路由
	r.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})
	log.Println("🚀 Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gin server error:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//1. 升级ws连接
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer wsConn.Close()

	//2. 创建 WebRTC API 实例
	m := webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		log.Println("Register codecs error:", err)
		return
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	//3. 创建 PeerConnection
	peerConn, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("Create PeerConnection error:", err)
		return
	}
	defer peerConn.Close()

	//4. 添加音频接收器
	if _, err := peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		log.Println("AddTransceiver error:", err)
		return
	}

	//5. 音频轨道处理
	peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Println("✅ Got audio track:", track.Codec().MimeType)

		go func() {
			for {
				rtpPacket, _, err := track.ReadRTP()
				if err != nil {
					log.Println("Track read error:", err)
					return
				}
				log.Printf("📦 Received RTP packet: SSRC=%d Seq=%d TS=%d Size=%d\n",
					rtpPacket.SSRC, rtpPacket.SequenceNumber, rtpPacket.Timestamp, len(rtpPacket.Payload))
			}
		}()
	})

	//6. ICE Candidate 回传
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
		if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("WriteMessage ICE error:", err)
		}
	})

	//7. 接收 SDP Offer
	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		log.Println("Read offer error:", err)
		return
	}
	var offer webrtc.SessionDescription
	if err := json.Unmarshal(msg, &offer); err != nil {
		log.Println("Unmarshal offer error:", err)
		return
	}
	if err := peerConn.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription error:", err)
		return
	}

	//8. 生成并发送 SDP Answer
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer error:", err)
		return
	}
	if err := peerConn.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription error:", err)
		return
	}
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		log.Println("Marshal answer error:", err)
		return
	}
	if err := wsConn.WriteMessage(websocket.TextMessage, answerJSON); err != nil {
		log.Println("Write answer error:", err)
		return
	}

	//9. 处理客户端 ICE candidates
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
