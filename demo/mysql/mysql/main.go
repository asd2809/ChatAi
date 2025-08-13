package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// User 模型，ID为UUID字符串
type User struct {
    ID       string `gorm:"type:char(36);primaryKey" json:"id"`
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

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

func main() {
    dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("数据库连接失败: " + err.Error())
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        panic("自动迁移失败: " + err.Error())
    }

    r := gin.Default()

    // 使用 CORS 中间件
    r.Use(CORSMiddleware())

    r.POST("/api/v1/users/register", handleRegister)
    r.GET("/api/v1/users/:userID", handleGetUser)

    r.Run(":8080")
}

func handleRegister(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 唯一性校验
    var count int64
    db.Model(&User{}).Where("username = ?", req.Username).Count(&count)
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
        return
    }

    db.Model(&User{}).Where("email = ?", req.Email).Count(&count)
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已注册"})
        return
    }

    // 生成UUID作为主键
    userID := uuid.NewString()

    user := User{
        ID:       userID,
        Username: req.Username,
        Password: req.Password, // 生产环境请加密
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
        },
    })
}

func handleGetUser(c *gin.Context) {
    userID := c.Param("userID")

    var user User
    if err := db.First(&user, "id = ?", userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
}