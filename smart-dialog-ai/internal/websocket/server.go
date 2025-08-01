// web与go之间建立websocket连接,主要是用于web与go之间测试llm的反应
package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)
type WebSocketServer struct {
	upgrader  *websocket.Upgrader // 用于升级HTTP连接到WebSocket
	clients map[*websocket.Conn]bool // 连接的客户端
	broadcast chan []byte            // 广播消息的通道
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有跨域请求
    },
}

// 初刷化结构体webSocetServer
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		upgrader:  &upgrader,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}
// 建立web与go之间的连接
func (s *WebSocketServer) HandleConnection(c *gin.Context) {
	// 升级HTTP连接到WebSocket
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 使用 Gin 的错误处理方法返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
        return
	}
	defer conn.Close()

	s.clients[conn] = true

	for {
		messageType, msg, err := conn.ReadMessage()
		log.Printf("Received message: %s", msg)
		if err != nil {
			delete(s.clients, conn)
			break
		}


		for client := range s.clients {
			if err := client.WriteMessage(messageType, msg); err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
	}

}