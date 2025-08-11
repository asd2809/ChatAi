package model

import "gorm.io/gorm"

// type User struct {
// 	gorm.Model
// 	Username string  `gorm:"type:varchar(255);unique;not null"`    // 必填且唯一
// 	Password string  `gorm:"type:varchar(255);not null"`           // 必填
// 	Phone    *string `gorm:"type:varchar(20);unique" json:"phone"` // 可选，唯一
// 	Email    string  `gorm:"type:varchar(255);unique;not null"`    // 必填且唯一
// }
// 添加了格式校验
type User struct {
	gorm.Model
	Username string  `gorm:"type:varchar(255);unique;not null" json:"username" binding:"required"`
	Password string  `gorm:"type:varchar(255);not null" json:"password" binding:"required,password"`
	Phone    *string `gorm:"type:varchar(20);unique" json:"phone" binding:"omitempty,phonecn"`
	Email    string  `gorm:"type:varchar(255);unique;not null" json:"email" binding:"required,email"`
}

// 注册请求参数结构体
type RegisterRequest struct {
	// 代表这个手机号不是必须的，如果填写了必须是合法的中国手机号
	Phone    string `json:"phone" binding:"omitempty,phonecn"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}