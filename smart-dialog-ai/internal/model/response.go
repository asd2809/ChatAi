package model


type Result struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}
type ResponseWeb struct {
	Result   Result `json:"result"`      // 业务响应内容
}