package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)


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
	UserID    string    `gorm:"index"`
	Role      string    `gorm:"type:enum('user','assistant')"`
	Content   string
	CreatedAt time.Time
}


// 获取历史的聊天记录
func SaveMessage(db *gorm.DB,userID string ,msg Message){
	record := ChatRecord{
		UserID: userID,
		Role: msg.Role,
		Content: msg.Content,
	}

	if err := db.Create(&record).Error; err!=nil{
		fmt.Println("保存失败：",err)
	}
}
// 保存聊天记录
func LoadHistory(db *gorm.DB , userID string)[]Message{
	var records []ChatRecord
	var history []Message

	// 查询指定用户的所有聊天记录
	db.Where("user_id = ?",userID).Order("created_at asc").Find(&records)

	// 把所有的记录全部放到一个变量中
	for _,r := range records{
		history = append(history,Message{
			Role:r.Role,
			Content: r.Content,
		})
	}
	return history
}

