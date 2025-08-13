package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 是签发和验证 JWT 的密钥，建议用环境变量或安全存储
var jwtSecret = []byte("超级密钥123456")

// LoginRequest 绑定登录请求的 JSON 参数结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GenerateToken 用于根据 userID 生成带有过期时间的 JWT Token
func GenerateToken(userID string) (string, error) {
	// jwt包含三部分,header(由jwt库自动生成) , payload,以及Signature(签名)
	// jwt.MapClaims 是 JWT 载荷（Payload），这里放自定义字段和标准字段 exp（过期时间）
	claims := jwt.MapClaims{
		"user_id": userID, // 自定义字段，通常表示你想要在token里带的用户相关信息,比如用户id，用户名等
		// 可以定义很多
		"exp": time.Now().Add(24 * time.Hour).Unix(), // 过期时间：当前时间加24小时，Unix时间戳格式
	}

	// 使用 HS256 签名算法创建一个新 token，并赋予上述 claims
	// 创建一个新的 JWT Token 对象，并告诉它用 HS256 算法来对载荷（claims）进行签名。
	// jwt本身就是一个token
	// 第一个参数是代表后续签名使用的算法(HS256),第二个参数是之前准备好的Payload
	// 这一步并不生成最终的 JWT 字符串，只是把 Header（签名算法等元信息）和 Payload（用户数据）
	// 封装成一个 Token 对象，准备后续签名。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对 token 进行签名，生成最终的 token 字符串
	return token.SignedString(jwtSecret)
}

// 登录处理函数，接受用户名密码，返回 JWT token
func loginHandler(c *gin.Context) {
	var req LoginRequest
	// 绑定 JSON 请求参数到结构体，并校验必填
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 简单示例，固定用户名密码，实际应查询数据库验证
	if req.Username != "admin" || req.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 调用 GenerateToken 生成 JWT
	token, err := GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 返回 token 给客户端
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// JWT 验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 字段
		// 这个是一种约定
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带token"})
			c.Abort() // 终止请求，阻止访问后续 Handler
			return
		}

		// 按空格拆分，期望格式是 "Bearer <token>"
		// 一种约定与规范
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求token格式错误"})
			c.Abort()
			return
		}

		tokenString := parts[1] // 取出真正的 token 字符串部分

		// 解析 token，第二个参数是一个回调函数，返回用于校验的密钥
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 校验签名算法是否符合预期，防止被篡改
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("非法的签名算法")
			}
			return jwtSecret, nil
		})

		// 解析失败或者 token 无效（过期、篡改等）
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 断言 token.Claims 类型为 jwt.MapClaims，即 map[string]interface{}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token解析失败"})
			c.Abort()
			return
		}

		// 从 claims 中取出自定义字段 user_id
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token缺少user_id"})
			c.Abort()
			return
		}

		// 将 userID 写入 gin 上下文，供后续 handler 使用
		c.Set("userID", userID)

		// 继续执行后续 handler
		c.Next()
	}
}

// 受保护接口示例，返回当前用户信息
func profileHandler(c *gin.Context) {
	// 从上下文读取 userID
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
func main() {
	r := gin.Default()

	// 进行跨域处理
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"}, // 允许前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 你的路由和中间件
	r.POST("/login", loginHandler)

	authGroup := r.Group("/user")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/profile", profileHandler)
	}

	r.Run(":8082")
}
