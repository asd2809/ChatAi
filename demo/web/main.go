package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 消息记录（内存）
var messages []string
var msgMutex sync.Mutex

// 升级为 WebSocket 的 upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许跨域请求
		return true
	},
}

// WebSocket 连接处理
func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	defer conn.Close()

	log.Println("客户端已连接")

	for {
		// 读取消息
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}
		log.Printf("收到消息: %s\n", msg)

		// 保存消息（线程安全）
		msgMutex.Lock()
		messages = append(messages, string(msg))
		msgMutex.Unlock()

		// 回复客户端（模拟AI回复）
		reply := "你说的是: " + string(msg)
		conn.WriteMessage(msgType, []byte(reply))

		// 同样记录回复
		msgMutex.Lock()
		messages = append(messages, reply)
		msgMutex.Unlock()
	}
}

// 获取历史消息
func getMessagesHandler(c *gin.Context) {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func main() {
	router := gin.Default()

	// ✅ 添加全局中间件，手动设置 CORS 头部
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 或指定: http://localhost:5500
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// WebSocket 路由
	router.GET("/ws", websocketHandler)

	// 获取聊天记录
	router.GET("/get-messages", getMessagesHandler)

	// 提供静态页面（如你部署前端）
	// router.StaticFile("/", "./index.html")

	log.Println("服务器启动于 http://localhost:8083")
	if err := router.Run(":8083"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
