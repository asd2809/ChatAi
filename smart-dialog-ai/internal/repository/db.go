package repository
// 主要是用于初始化数据库
import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
type DB struct{
	DB *gorm.DB
	UserID string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// 存储到 MySQL 的结构
type ChatRecord struct {
	ID        uint      `gorm:"primaryKey"`
	// UserID    string    `gorm:"index"`
	Role      string    `gorm:"type:enum('user','assistant')"`
	Content   string
	CreatedAt time.Time
}

// 获取实例
func NewDB(db *gorm.DB,userID string) *DB{
	return &DB{
		DB: db,
		UserID: userID,
	}
}

// 初始化
func InitDB() *gorm.DB{
	dsn := "root:root@cbj@tcp(127.0.0.1:3306)/chat_ai?charset=utf8mb4&parseTime=True&loc=Local"
	db,err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil{
		panic("mysql connect failed"  +  err.Error())
	}

	// 自动建表

	if err := db.AutoMigrate(&ChatRecord{}) ; err != nil{
		panic("mysql automigrate failed：" + err.Error())
	}

	return db
}