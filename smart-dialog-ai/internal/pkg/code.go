package pkg

import "net/http"

// 业务码与 HTTP 状态码定义
type BizCode struct {
	Code       int
	HttpStatus int
	Message    string
}

var (
	CodeSuccess              = BizCode{0, http.StatusOK, "成功"}
	CodeUserNotFound         = BizCode{10001, http.StatusNotFound, "用户不存在"}
	CodeUserExists           = BizCode{10002, http.StatusBadRequest, "用户名已存在"}
	CodeEmailExists          = BizCode{10003, http.StatusBadRequest, "邮箱已注册"}
	CodePasswordError        = BizCode{10004, http.StatusUnauthorized, "密码错误"}
	CodeParamInvalid         = BizCode{10005, http.StatusBadRequest, "参数校验失败"}
	CodePasswordEncryptError = BizCode{10006, http.StatusInternalServerError, "密码加密失败"}
	CodeInvalidUsername      = BizCode{10007, http.StatusBadRequest, "用户名不能为空"}
	CodeInvalidPassword      = BizCode{10008, http.StatusBadRequest, "密码不能为空"}
	CodeInvalidEmail         = BizCode{10009, http.StatusBadRequest, "邮箱不能为空"}
	CodePhoneExists          = BizCode{100010, http.StatusBadRequest, "手机号已存在"}
	CodeDBError              = BizCode{10050, http.StatusInternalServerError, "数据库操作失败"}
	CodeUnknownError         = BizCode{10051, http.StatusInternalServerError, "未知内部错误"}
)
