package websocket

// 用于处理信息的发送与接收

import (
	"encoding/json"
	"smart-dialog-ai/internal/model"
	"smart-dialog-ai/internal/repository"
	"time"

	"smart-dialog-ai/internal/service"

	"github.com/sirupsen/logrus"
)

// MessageHandle 结构体用于处理WebSocket消息
type MessageHandle struct {
	WebSocketServer *WebSocketServer            // go与web之间的连接
	Logrus          logrus.Logger               // 日志记录器
	llm             *service.SiliconFlowHandler // 硅谷流处理器(大模型)
	db              *repository.DB                     // 数据库
}

// 初始化结构体
func NewMessageHandle(server *WebSocketServer, llm *service.SiliconFlowHandler, db *repository.DB) *MessageHandle {
	return &MessageHandle{
		WebSocketServer: server,
		Logrus:          *logrus.New(),
		llm:             llm, // 初始化硅谷流处理器
		db:              db,  //数据库的结构体
	}
}

// 处理web与go之间的发送与接收消息
func (s *MessageHandle) HandleMessage() {
	s.Logrus.Info("开始处理消息")

	for {
		messageType, msg, err := s.WebSocketServer.conn.ReadMessage()
		
		// 保存聊天记录
		userContent := model.Message{
			Role:    "user",
			Content: string(msg),
		}
		s.llm.History = append(s.llm.History,userContent)
		repository.SaveMessage(s.db.DB,s.db.UserID,userContent)
		
		if err != nil {
			s.Logrus.Error("消息读取失败:", err)
			break
		}

		// 1. 解析前端发送的 JSON 消息
		var userMsg model.ChatMessage
		if err := json.Unmarshal(msg, &userMsg); err != nil {
			s.Logrus.Error("JSON解析失败:", err)
			break
		}
		s.Logrus.Infof("收到用户消息: %+v", userMsg)

		// 2. 传给 LLM 生成回复
		replyText, err := s.llm.GenerateText(userMsg.Text)
		s.Logrus.Printf("打印出llm回复内容%s",replyText)
		
		// 保存聊天记录
		aiContent := model.Message{
			Role: "assistant",
			Content: replyText,
		}
		s.llm.History = append(s.llm.History, aiContent)
		repository.SaveMessage(s.db.DB,s.db.UserID,aiContent)


		if err != nil {
			s.Logrus.Error("调用LLM失败:", err)
			break
		}

		// 3. 构造返回消息
		reply := model.ChatMessage{
			Role:      "assistant",
			Text:      replyText,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		// 4. 编码 JSON 发送给前端
		jsonBytes, err := json.Marshal(reply)
		if err != nil {
			s.Logrus.Error("JSON编码失败:", err)
			break
		}

		err = s.WebSocketServer.conn.WriteMessage(messageType, jsonBytes)
		if err != nil {
			s.Logrus.Error("发送消息失败:", err)
			return
		}
	}
}
