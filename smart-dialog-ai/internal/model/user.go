package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"type:varchar(255);unique;not null"`    // 必填且唯一
	Password string  `gorm:"type:varchar(255);not null"`           // 必填
	Phone    *string `gorm:"type:varchar(20);unique" json:"phone"` // 可选，唯一
	Email    string  `gorm:"type:varchar(255);unique;not null"`    // 必填且唯一
}
// 注册请求参数结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}