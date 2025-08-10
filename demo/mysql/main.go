package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"`
}

// 注册请求参数结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

var db *gorm.DB

func main() {
	// 连接数据库，记得替换成你的账号密码和数据库名
	dsn := "root:root@cbj@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移，创建 users 表
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("自动迁移失败: " + err.Error())
	}

	r := gin.Default()

	r.POST("/api/v1/users/register", handleRegister)

	r.Run(":8080")
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest

	// 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var count int64
	db.Model(&User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	db.Model(&User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已注册"})
		return
	}

	// 创建用户
	user := User{
		Username: req.Username,
		Password: req.Password, // 生产环境密码请加密！
		Email:    req.Email,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功", 
		"user": gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}})
}
