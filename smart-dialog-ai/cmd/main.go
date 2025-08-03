package main

import (
	"log"
	"smart-dialog-ai/internal/websocket"
	"smart-dialog-ai/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
    // 加载配置
    
    // 初始化 API 服务
    

	// 建立web与go之间的连接
	websocketServer := websocket.NewWebSocketServer()
	router := gin.Default()
	router.GET("/ws", websocketServer.HandleConnection)
	log.Println("服务器启动于: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalln("启动失败:", err)
	}
	// 初始化llm
	url := "https://api.siliconflow.cn/v1/chat/completions"        // 替换为你的 API 地址
	token := "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为有效的 API Token
	llm := service.NewSiliconFlowHandler(url, token)
	// 初始化server消息的结构体
	 websocket.NewMessageHandle(websocketServer,llm)

}
