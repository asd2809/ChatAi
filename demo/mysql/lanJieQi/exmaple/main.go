package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("超级密钥123456")

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func loginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }
    if req.Username != "admin" || req.Password != "123456" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }
    token, err := GenerateToken(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
        return
    }

    // 设置 HttpOnly Cookie，自动发送给客户端浏览器
    c.SetCookie("jwt_token", token, 3600*24, "/", "localhost", false, true)

    c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("jwt_token")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带token"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("非法的签名算法")
            }
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token解析失败"})
            c.Abort()
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token缺少user_id"})
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Next()
    }
}

func profileHandler(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息丢失"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "欢迎访问个人信息页",
        "userID":  userID,
    })
}

func logoutHandler(c *gin.Context) {
    // 删除cookie，设置MaxAge为负数
    c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "已登出"})
}


func main() {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5500"}, // 允许你的前端地址
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        AllowCredentials: true, // 允许携带 cookie
        MaxAge:           12 * time.Hour,
    }))

    // 你的路由注册和逻辑...
    r.POST("/login", loginHandler)
    authGroup := r.Group("/user")
    authGroup.Use(AuthMiddleware())
    {
        authGroup.GET("/profile", profileHandler)
        authGroup.POST("/logout", logoutHandler)
    }

    r.Run(":8082")
}
