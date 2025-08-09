// web与go之间建立websocket连接,主要是用于web与go之间测试llm的反应
package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	upgrader  *websocket.Upgrader // 用于升级HTTP连接到WebSocket
	conn      *websocket.Conn    //用于web与go之间联系websocket的客户端
	clients   map[*websocket.Conn]bool // 连接的客户端
	broadcast chan []byte              // 广播消息的通道
	messageHandle *MessageHandle
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
		upgrader: &upgrader,
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
	// 把web与go之间客户端存放如websocketServer
	s.conn =conn
	s.clients[conn] = true
	// 读取web发来的消息
	// 实例化 MessageHandle，并调用 HandleMessage
	// 获取message实例
	s.messageHandle.HandleMessage()
	

}
// 存放messageHandle
func(s *WebSocketServer) GetMessageHandle(messageHandler *MessageHandle){
	s.messageHandle = messageHandler
}