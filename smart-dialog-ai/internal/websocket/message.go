package websocket
// 用于处理信息的发送与接收

type MessageHandle struct {
	WebSocketServer *WebSocketServer  //go与web之间的连接

}

// 接收Web发来的消息