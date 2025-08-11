package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义密码校验：必须包含大小写字母和数字，长度6-12
func PasswordValidator(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 6 || len(pwd) > 12 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pwd)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pwd)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(pwd)
	return hasUpper && hasLower && hasDigit
}

// 自定义中国手机号校验：11位数字，1开头
func PhoneCNValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值不校验（omitempty）
	}
	match, _ := regexp.MatchString(`^1\d{10}$`, phone)
	return match
}

// 注册自定义校验器函数，供外部调用
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", PasswordValidator)
	v.RegisterValidation("phonecn", PhoneCNValidator)
}
