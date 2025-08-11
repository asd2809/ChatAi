package pkg

import "fmt"

// 业务错误接口
type BizErrorInterface interface {
	// 还是不太明白为什么要加error
	// 但意思是：为了让业务错误和普通错误无缝互通，同时又能保留额外的业务信息
	error
	Code() int
	HTTPStatus() int
	Message() string
}
	// 业务错误快速创建
	// 不再返回一个单一错误的结构体，而是可以返回只要能实现这个接口的错误结构体，
	// 如果后续有其他类型错误的结构体可以进行扩展
	// 类比就是，以前你只能认识一个人，现在你可以通过性别，姓名等认识很多人

	// 工厂函数，快速创建业务错误
func NewBizError(code BizCode) BizErrorInterface {
	return &BizError{
		BizCode: code,
	}
}
// 业务错误实现，内嵌 BizCode
type BizError struct {
	BizCode	BizCode
}
// 实现BizeError的四个方法，这样可以调用定义的接口
func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code(), e.Message())
}
func (e *BizError) Code() int       { return e.BizCode.Code }
func (e *BizError) HTTPStatus() int { return e.BizCode.HttpStatus }
func (e *BizError) Message() string { return e.BizCode.Message }


