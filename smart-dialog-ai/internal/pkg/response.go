package pkg

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// 统一响应结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 统一成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(CodeSuccess.HttpStatus, Result{
		Code:    CodeSuccess.Code,
		Message: CodeSuccess.Message,
		Data:    data,
	})
}

// 统一错误响应，支持 errors.As 识别业务错误接口
func Error(c *gin.Context, err error) {
	var bizErr BizErrorInterface
	if errors.As(err, &bizErr) {
		c.JSON(bizErr.HTTPStatus(), Result{
			Code:    bizErr.Code(),
			Message: bizErr.Message(),
			Data:    nil,
		})
	} else {
		c.JSON(CodeUnknownError.HttpStatus, Result{
			Code:    CodeUnknownError.Code,
			Message: err.Error(),
			Data:    nil,
		})
	}
}
