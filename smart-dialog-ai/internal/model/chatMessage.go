package model

import "time"


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
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Role      string `gorm:"type:enum('user','assistant')"`
	Content   string
	CreatedAt time.Time
}