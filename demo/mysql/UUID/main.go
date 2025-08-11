package main

import (
	"errors"
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

// 统一响应结构体
type Result struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据
}

var db *gorm.DB

func main() {
	// 连接数据库，替换为你自己的 DSN
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("自动迁移失败: " + err.Error())
	}

	r := gin.Default()

	r.GET("/api/v1/users/:id", getUserByID)

	r.Run(":8080")
}

// 统一响应函数，直接传 Result 结构体
func Response(c *gin.Context, httpCode int, res Result) {
	c.JSON(httpCode, res)
}

// 辅助封装 - 成功响应，简化调用
func Success(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, Result{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// 辅助封装 - 错误响应，简化调用
func Error(c *gin.Context, httpCode, code int, msg string) {
	Response(c, httpCode, Result{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// 示例接口：获取用户 - 直接使用 Response 调用（标准调用方式）
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	var user User
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 直接调用 Response 返回完整结构体
			Response(c, http.StatusNotFound, Result{
				Code:    40401,
				Message: "用户不存在",
				Data:    nil,
			})
		} else {
			Response(c, http.StatusInternalServerError, Result{
				Code:    50001,
				Message: "数据库查询失败",
				Data:    nil,
			})
		}
		return
	}

	// 直接调用 Response 返回数据
	Response(c, http.StatusOK, Result{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// 你也可以用下面的辅助函数写法，简化代码：
//
// func getUserByID(c *gin.Context) {
//     id := c.Param("id")
//
//     var user User
//     err := db.First(&user, id).Error
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             Error(c, http.StatusNotFound, 40401, "用户不存在")
//         } else {
//             Error(c, http.StatusInternalServerError, 50001, "数据库查询失败")
//         }
//         return
//     }
//
//     Success(c, gin.H{
//         "id":       user.ID,
//         "username": user.Username,
//         "email":    user.Email,
//     })
// }
