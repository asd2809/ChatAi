package api

import (
	"net/http"
	"smart-dialog-ai/internal/repository"

	"github.com/gin-gonic/gin"
)

// 存放web的请求


func (g *GinWrapper) HandleLoadHistory(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	records := repository.GetUserChatHistory(g.DB, userID)
	c.JSON(http.StatusOK, gin.H{"messages": records})
}

// HandleClearHistory 清空聊天记录的处理函数
func (g *GinWrapper) HandleClearHistory(c *gin.Context) {
	userID := c.Query("user_id") // 假设用户ID是通过查询参数传递的

	err := repository.ClearUserChatHistory(g.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat history cleared successfully"})
}

// 用户管理的方法
func (g *GinWrapper) HandleGetUser(c *gin.Context){

}
func (g *GinWrapper) HandleRegister(c *gin.Context){

}
func (g *GinWrapper) HandleUpdarerUser(c *gin.Context){

}
func (g *GinWrapper) HandleDeleteUser(c *gin.Context){

}