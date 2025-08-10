package api

import (
	
	"log"
	"smart-dialog-ai/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 控制gin的调用
type GinWrapper struct {
	Gin *gin.Engine
	DB *gorm.DB
}

func NewGinWrapper(db *gorm.DB) *GinWrapper {
	return &GinWrapper{
		Gin: gin.Default(),
		DB: db,
	}
}
// 调用路由的主要方法
func (g *GinWrapper) SetupWebSocketAndRoutes(websocketServer *websocket.WebSocketServer) {
	g.Gin.GET("/ws", websocketServer.HandleConnection)
	g.Gin.GET("/chatAll", g.HandleLoadHistory)
	g.Gin.GET("/chatHistory",g.HandleClearHistory)

	

	// 用户认证路由组
	// authGroup := g.Gin.Group("/api/v1/auth")
	// {
	// 	authGroup.POST("/sendVerificationCode",g.HandleSendVerificationCode)
	// 	authGroup.POST("/verifyCode",g.HandleVerfyCode)
	// 	authGroup.POST("/loginWithPhone",g.HandleLoginWithPhone)
	// 	authGroup.POST("/logout",g.HandleLogout)
	// }
	// 用户管理路由组
	userGroup := g.Gin.Group("/api/v1/users")
	{
	
		// 发起请求的方式GET /api/v1/users/:userID
		// GET /api/v1/users/123 , 123就是userID
		// 获取c.Param("userID")获取路径参数userID的值
		// c的数据类型是 *gin.Context

		userGroup.GET("/:userID",g.HandleGetUser)
		userGroup.POST("/register",g.HandleRegister)
		userGroup.PUT("/:userID",g.HandleUpdarerUser)
		// 通过userID删除数据
		userGroup.DELETE("/:userID",g.HandleDeleteUser)
		// userGroup.GET("/",g.HandleListUsers) //通常只对管理员开放
	
	}
	
	log.Println("服务器启动于: http://localhost:8080")
	if err := g.Gin.Run(":8080"); err != nil {
		log.Fatalln("启动失败:", err)
	}

}	