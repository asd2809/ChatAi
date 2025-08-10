package repository

import (
	"fmt"
	"smart-dialog-ai/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 保存聊天记录
func SaveMessage(db *gorm.DB,userID string ,msg model.Message){
	logrus.Info("开始保存聊天记录")
	record := model.ChatRecord{
		UserID: userID,
		Role: msg.Role,
		Content: msg.Content,
	}

	if err := db.Create(&record).Error; err!=nil{
		fmt.Println("保存失败：",err)
	}
}
// 获取聊天记录
func LoadHistory(db *gorm.DB , userID string)[] model.Message{
	logrus.Info("开始获取聊天记录")

	var records []model.ChatRecord
	var history []model.Message

	// 查询指定用户的所有聊天记录
	db.Where("user_id = ?",userID).Order("created_at asc").Find(&records)

	// 把所有的记录全部放到一个变量中
	for _,r := range records{
		history = append(history,model.Message{
			Role:r.Role,
			Content: r.Content,
		})
	}
	return history
}

// GetUserChatHistory 获取指定用户的聊天记录，用于前端显示
func GetUserChatHistory(db *gorm.DB, userID string) []model.ChatRecord {
	var records []model.ChatRecord
	db.Where("user_id = ?", userID).Order("created_at asc").Find(&records)
	return records
}
func ClearUserChatHistory(db *gorm.DB, userID string) error {
	// 执行删除操作，删除指定用户的所有聊天记录
	result := db.Where("user_id = ?", userID).Delete(&model.ChatRecord{})
	if result.Error != nil {
		// 如果删除过程中出现错误，返回错误
		return result.Error
	}
	if result.RowsAffected == 0 {
		// 如果没有记录被删除，可能表示没有找到对应的聊天记录，或者用户ID错误
		return fmt.Errorf("no chat records found for user %s", userID)
	}
	return nil
}

