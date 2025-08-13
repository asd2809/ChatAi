package middleware

import (
	
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
    userGroup := r.Group("/user")
    userGroup.Use(AuthMiddleware())
    {
        userGroup.GET("/profile", profileHandler)
        // 更多用户相关路由
    }
}
func profileHandler(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息丢失"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "欢迎访问个人信息页",
        "userID":  userID,
    })
}