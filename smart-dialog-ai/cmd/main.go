package main

import (
	"log"
	"smart-dialog-ai/internal/repository"
	"smart-dialog-ai/internal/service"
	"smart-dialog-ai/internal/websocket"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	apiUrl := "https://api.siliconflow.cn/v1/chat/completions"        // 替换为你的 API 地址
// 	apiKey := "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为有效的 API Token
// 	websocketServer := websocket.NewWebSocketServer()

// 	llm := service.NewSiliconFlowHandler(apiUrl, apiKey)

// 	historyManager := websocket.NewChatHistoryManager()

// 	messageHandler := websocket.NewMessageHandle(websocketServer, llm, historyManager)

// 	router := gin.Default()
// 	router.GET("/ws", func(c *gin.Context) {
// 		websocketServer.HandleConnection(c)
// 		// 启动消息处理协程，注意 websocketServer.conn 需要被正确赋值
// 		go messageHandler.HandleMessage()
// 	})

// 	log.Println("服务器启动于: http://localhost:8080")
// 	router.Run(":8080")
// }

func main() {
	// 加载配置

	// 初始化 API 服务
	db := repository.InitDB()
	// 获取db的实例
	dbNew := repository.NewDB(db,"user1")
	// 初始化llm
	url := "https://api.siliconflow.cn/v1/chat/completions"        // 替换为你的 API 地址
	token := "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为有效的 API Token
	llm := service.NewSiliconFlowHandler(url, token,db)

	// 建立web与go之间的连接,初始化
	websocketServer := websocket.NewWebSocketServer()
	// 初始化server消息的结构体
	message := websocket.NewMessageHandle(websocketServer, llm,dbNew)
	// 获取message实例
	websocketServer.GetMessageHandle(message)
	// 启动路由
	router := gin.Default()
	router.GET("/ws", websocketServer.HandleConnection)
	log.Println("服务器启动于: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalln("启动失败:", err)
	}

}
