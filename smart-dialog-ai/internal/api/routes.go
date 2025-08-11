package api

import (
	"errors"
	"net/http"
	"smart-dialog-ai/internal/model"
	"smart-dialog-ai/internal/pkg"
	"smart-dialog-ai/internal/repository"

	"github.com/gin-gonic/gin"
)

// 处理web请求中函数的逻辑
func (g *GinWrapper) HandleLoadHistory(c *gin.Context) {
	userID := "user1"
	// userID := c.Query("user_id")
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

//------------------------用户管理的方法-------------------------------
func (g *GinWrapper) HandleGetUser(c *gin.Context){
	// // 获取url路径中的id参数
	// id := c.Param("id")
	// var user User 
	// // 根据ID查询用户
	// if err := repository.GetUser(g.DB,id); err != nil{
	// 	c.JSON(htt)
	// }
}
// 增加一条新的用户
func (g *GinWrapper) HandleRegister(c *gin.Context){
	var reg model.RegisterRequest

	// 绑定前端传来的信息
	if err := c.ShouldBindJSON(&reg);err !=nil{
		Error(c,pkg.NewBizError(pkg.CodeParamInvalid))
		return 
	}
	// 写入数据库
	err := repository.Register(g.DB,reg)
	if err != nil{
		Error(c,err)
		return
	}
	Success(c,gin.H{
		"message":"注册成功",
	})
}
func (g *GinWrapper) HandleUpdarerUser(c *gin.Context){
	
}
func (g *GinWrapper) HandleDeleteUser(c *gin.Context){

}

// 统一响应结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 统一成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(pkg.CodeSuccess.HttpStatus, Result{
		Code:    pkg.CodeSuccess.Code,
		Message: pkg.CodeSuccess.Message,
		Data:    data,
	})
}

// 统一错误响应，支持 errors.As 识别业务错误接口
func Error(c *gin.Context, err error) {
	var bizErr pkg.BizErrorInterface
	if errors.As(err, &bizErr) {
		c.JSON(bizErr.HTTPStatus(), Result{
			Code:    bizErr.Code(),
			Message: bizErr.Message(),
			Data:    nil,
		})
	} else {
		c.JSON(pkg.CodeUnknownError.HttpStatus, Result{
			Code:    pkg.CodeUnknownError.Code,
			Message: err.Error(),
			Data:    nil,
		})
	}
}