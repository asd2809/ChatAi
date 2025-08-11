package main

import (
	"net/http"
	"smart-dialog-ai/internal/api"
	"smart-dialog-ai/internal/repository"
	"smart-dialog-ai/internal/service"
	"smart-dialog-ai/internal/utils"
	"smart-dialog-ai/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CORSMiddleware 是一个中间件，用于处理 CORS 相关的响应头
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")                             // 允许所有域的请求
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")           // 允许 GET 和 POST 请求
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept") // 允许自定义请求头
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next() // 继续处理请求
		}
	}
}

func main() {
	// 注册自定义校验器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		utils.RegisterCustomValidators(v)
	}
	// 加载配置

	// 初始化mysql服务
	db := repository.InitDB()
	// 获取db的实例
	dbNew := repository.NewDB(db, "user1")
	// 初始化llm
	url := "https://api.siliconflow.cn/v1/chat/completions"        // 替换为你的 API 地址
	token := "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为有效的 API Token
	llm := service.NewSiliconFlowHandler(url, token, db)

	// 建立web与go之间的连接,初始化
	websocketServer := websocket.NewWebSocketServer()
	// 初始化server消息的结构体
	message := websocket.NewMessageHandle(websocketServer, llm, dbNew)
	// 获取message实例
	websocketServer.GetMessageHandle(message)
	// 启动路由
	router := api.NewGinWrapper(db)
	router.SetupWebSocketAndRoutes(websocketServer)
	// // 前端获取所有聊天记录的gin
	// router.Gin.GET("/chatAll",router.Ss(debNew))

	// 使用 CORS 中间件
	router.Gin.Use(CORSMiddleware())

}
