你现在的后端日志表明：

成功收到了音频轨道

并且持续不断地收到了 RTP 包，且打印出了包的序号（Seq）、时间戳（TS）、大小（Size）

这就说明音频流正在实时传输过来了，WebRTC 完整建立且音频数据已经传输成功。

下一步你可以做的事情示例：
保存音频数据：将 rtpPacket.Payload 以 opus 格式写文件，或者做转码保存为 PCM

接入语音识别：把这些 RTP 或解码后的音频流送给腾讯/讯飞等 ASR 做语音转文字

实现语音交互：识别结果回传前端做展示，甚至做 TTS 回复音频

如果你需要，我可以帮你写：

RTP 数据保存为 .opus 文件示例

RTP 解码转 PCM 示例

连接第三方 ASR 接口示例

或者前端显示文字结果的完整方案

你想往哪个方向推进？



package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade:", err)
		return
	}
	defer wsConn.Close()

	m := webrtc.MediaEngine{}
	m.RegisterDefaultCodecs()
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))

	peerConn, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatal("Create PeerConnection:", err)
		return
	}
	defer func() {
		err := peerConn.Close()
		if err != nil {
			log.Println("PeerConnection close error:", err)
		}
	}()

	// 明确告知接收音频轨道
	_, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio)
	if err != nil {
		log.Fatal("AddTransceiver error:", err)
		return
	}

	peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Println("✅ Got audio track:", track.Codec().MimeType)

		// 创建文件保存音频负载
		file, err := os.Create("output.opus")
		if err != nil {
			log.Println("Failed to create output file:", err)
			return
		}
		defer file.Close()

		go func() {
			for {
				rtpPacket, _, err := track.ReadRTP()
				if err != nil {
					log.Println("Track read error:", err)
					return
				}

				// 写入 RTP 负载（Opus 编码）
				_, err = file.Write(rtpPacket.Payload)
				if err != nil {
					log.Println("File write error:", err)
					return
				}
			}
		}()
	})

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
		err = wsConn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("WriteMessage ICE candidate error:", err)
		}
	})

	// 读 Offer
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

	if err = peerConn.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription error:", err)
		return
	}

	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer error:", err)
		return
	}

	if err = peerConn.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription error:", err)
		return
	}

	answerJSON, err := json.Marshal(answer)
	if err != nil {
		log.Println("Marshal answer error:", err)
		return
	}

	if err = wsConn.WriteMessage(websocket.TextMessage, answerJSON); err != nil {
		log.Println("WriteMessage answer error:", err)
		return
	}

	// 循环处理客户端发来的 ICE candidates
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
保存音频
说明
运行这个程序后，客户端音频流的 Opus 编码负载会被写入当前目录的 output.opus 文件

录制结束后可用 VLC、Foobar2000 等播放器打开播放该文件

如果想要转 PCM 或 WAV，需要额外解码，后续可以帮你补充

如果你需要，我可以帮你写：

如何将 .opus 文件转 PCM 或 WAV

录制更长时间的处理（比如断开时关闭文件）

或者把录音直接传给语音识别接口