type WebSocketServer struct {
	upgrader  *websocket.Upgrader // 用于升级HTTP连接到WebSocket
	conn      *websocket.Conn       //用于web与go之间联系websocket的客户端
	clients   map[*websocket.Conn]bool // 连接的客户端
	broadcast chan []byte              // 广播消息的通道
	messageHandle *MessageHandle       //用于调用读写与web之间的消息的客户端
}
方法
    用于初始化WebSocketServer
    NewWebSocketServer
    用于向NeWebSockerServer初始化messageHandle这样可以调用
    GetMessageHandle