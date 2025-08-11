// package main

// import (
//     "fmt"
//     "net/http"

//     "github.com/gin-gonic/gin"
// )
// // 统一响应格式
// type Result struct {
//     Code    int         `json:"code"`
//     Message string      `json:"message"`
//     Data    interface{} `json:"data,omitempty"`
// }

// // 业务码和HTTP状态码定义
// type BizCode struct {
//     Code       int
//     HTTPStatus int
//     Message    string
// }
// // 自定义业务错误
// type BizError struct {
//     BizCode
// }
// // 这个是业务码与http状态码一一对应
// var (
//     CodeSuccess        = BizCode{0, http.StatusOK, "成功"}
//     CodeUserNotFound   = BizCode{10001, http.StatusNotFound, "用户不存在"}
//     CodeUserExists     = BizCode{10002, http.StatusBadRequest, "用户名已存在"}
//     CodeEmailExists    = BizCode{10003, http.StatusBadRequest, "邮箱已注册"}
//     CodePasswordError  = BizCode{10004, http.StatusUnauthorized, "密码错误"}
//     CodeParamInvalid   = BizCode{10005, http.StatusBadRequest, "参数校验失败"}
//     CodeDBError        = BizCode{10050, http.StatusInternalServerError, "数据库操作失败"}
//     CodeUnknownError   = BizCode{10051, http.StatusInternalServerError, "未知内部错误"}
// )

// func (e *BizError) Error() string {
//     return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
// }

// // 快捷创建业务错误
// func NewBizError(BizCode BizCode) *BizError {
//     return &BizError{
// 		BizCode:BizCode,
// 	}
// }

// // 预定义错误变量
// var (
//     ErrUserNotFound       = NewBizError(CodeUserNotFound)
//     ErrUserAlreadyExists  = NewBizError(CodeUserExists)
//     ErrEmailAlreadyExists = NewBizError(CodeEmailExists)
//     ErrPasswordIncorrect  = NewBizError(CodePasswordError)
//     ErrInvalidParameter   = NewBizError(CodeParamInvalid)
//     ErrDatabaseOperation  = NewBizError(CodeDBError)
//     ErrInternalServer     = NewBizError(CodeUnknownError)
// )

// // 统一成功响应
// func Success(c *gin.Context, data interface{}) {
//     c.JSON(CodeSuccess.HTTPStatus, Result{
//         Code:    CodeSuccess.Code,
//         Message: CodeSuccess.Message,
//         Data:    data,
//     })
// }

// // 统一错误响应（传入自定义错误）
// func Error(c *gin.Context, err error) {
// 	// 通过断言的方式，拿到业务码和对应的http状态码
// 	// 这样可以确保状态码与业务码一一对应
//     if be, ok := err.(*BizError); ok {
//         c.JSON(be.HTTPStatus, Result{
//             Code:    be.Code,
//             Message: be.Message,
//             Data:    nil,
//         })
//     } else {
//         // 非业务错误，返回服务器内部错误
//         c.JSON(CodeUnknownError.HTTPStatus, Result{
//             Code:    CodeUnknownError.Code,
//             Message: err.Error(),
//             Data:    nil,
//         })
//     }
// }

// // 业务逻辑示例
// func RegisterUser(username, email string) error {
//     // 模拟参数校验
//     if username == "" || email == "" {
//         return ErrInvalidParameter
//     }

//     // 模拟查重逻辑
//     if username == "existuser" {
//         return ErrUserAlreadyExists
//     }
//     if email == "exist@example.com" {
//         return ErrEmailAlreadyExists
//     }

//     // 模拟数据库插入操作失败的情况（这里假设成功）
//     // 如果失败，返回 ErrDatabaseOperation
//     // if err := dbInsertUser(...); err != nil {
//     //     return ErrDatabaseOperation
//     // }

//     // 注册成功
//     return nil
// }

// // Gin 控制器示例
// func RegisterHandler(c *gin.Context) {
//     var req struct {
//         Username string `json:"username" binding:"required"`
//         Email    string `json:"email" binding:"required,email"`
//     }

//     if err := c.ShouldBindJSON(&req); err != nil {
//         // ErrInvalidParameter是自定义的错误信息
// 		Error(c, ErrInvalidParameter)
//         return
//     }
// 	// 模拟数据库操作
//     err := RegisterUser(req.Username, req.Email)
//     if err != nil {
//         Error(c, err)
//         return
//     }

//     Success(c, gin.H{"message": "注册成功"})
// }

// // main.go 中如何启动
// func main() {
//     r := gin.Default()
//     r.POST("/register", RegisterHandler)
//     r.Run(":8083")
// }

// 进阶版本
package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

// 业务码与 HTTP 状态码定义
type BizCode struct {
	code       int
	httpStatus int
	message    string
}
// 业务码常量
var (
	CodeSuccess       = BizCode{0, http.StatusOK, "成功"}
	CodeUserNotFound  = BizCode{10001, http.StatusNotFound, "用户不存在"}
	CodeUserExists    = BizCode{10002, http.StatusBadRequest, "用户名已存在"}
	CodeEmailExists   = BizCode{10003, http.StatusBadRequest, "邮箱已注册"}
	CodePasswordError = BizCode{10004, http.StatusUnauthorized, "密码错误"}
	CodeParamInvalid  = BizCode{10005, http.StatusBadRequest, "参数校验失败"}
	CodeDBError       = BizCode{10050, http.StatusInternalServerError, "数据库操作失败"}
	CodeUnknownError  = BizCode{10051, http.StatusInternalServerError, "未知内部错误"}
)

// 业务错误快速创建
// 不再返回一个单一错误的结构体，而是可以返回只要能实现这个接口的错误结构体，
// 如果后续有其他类型错误的结构体可以进行扩展
// 类比就是，以前你只能认识一个人，现在你可以通过性别，姓名等认识很多人
func NewBizError(code BizCode) BizErrorInterface {
	return &BizError{
		BizCode: code,
	}
}
func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}
func (e *BizError) Code() int       { return e.code }
func (e *BizError) HTTPStatus() int { return e.httpStatus }
func (e *BizError) Message() string { return e.message }

// 添加这个接口的作用是，可以使其他类型的错误(包含其中三个方法)也可以调用这个接口
// 业务错误接口
type BizErrorInterface interface {
	// 还是不太明白为什么要加error
	// 但意思是：为了让业务错误和普通错误无缝互通，同时又能保留额外的业务信息
	error
	Code() int
	HTTPStatus() int
	Message() string
}

// 业务错误实现
type BizError struct {
	BizCode
}

// 统一响应结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 统一成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(CodeSuccess.httpStatus, Result{
		Code:    CodeSuccess.code,
		Message: CodeSuccess.message,
		Data:    data,
	})
}

// 统一错误响应
func Error(c *gin.Context, err error) {
	var bizErr BizErrorInterface
	//errors.As 会顺藤摸瓜地从一个可能被多层包装的错误里，把你想要的那种错误类型“剥”出来用
	// 不管你是 BizError 还是别的实现，只要符合接口，就能处理
	if errors.As(err, &bizErr) {
		c.JSON(bizErr.HTTPStatus(), Result{
			Code:    bizErr.Code(),
			Message: bizErr.Message(),
			Data:    nil,
		})
	} else {
		// 兜底未知错误
		c.JSON(CodeUnknownError.httpStatus, Result{
			Code:    CodeUnknownError.code,
			Message: err.Error(),
			Data:    nil,
		})
	}
}

// 业务层模拟注册函数
func RegisterUser(username, email string) error {
	if username == "" || email == "" {
		return NewBizError(CodeParamInvalid)
	}
	if username == "existuser" {
		return NewBizError(CodeUserExists)
	}
	if email == "exist@example.com" {
		return NewBizError(CodeEmailExists)
	}
	// 模拟数据库错误示范:
	// return NewBizError(CodeDBError)

	// 模拟成功
	return nil
}

// Gin 处理器
func RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 第二个参数是自定义的错误，里面包含状态码、业务码以及相对于的信息
		Error(c, NewBizError(CodeParamInvalid))
		return
	}

	err := RegisterUser(req.Username, req.Email)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, gin.H{"message": "注册成功"})
}

func main() {
	r := gin.Default()
	r.POST("/register", RegisterHandler)
	r.Run(":8080")
}
