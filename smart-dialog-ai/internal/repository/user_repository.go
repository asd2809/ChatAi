package repository

import (
	"errors"
	"fmt"
	"smart-dialog-ai/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// User 用户表结构体
type User struct {
    gorm.Model
    Username string `gorm:"type:varchar(255);unique;default:'default_username'"`
    Password string `gorm:"type:varchar(255);not null"`
    Phone    string `gorm:"type:varchar(20);unique"`
    Email    string `gorm:"type:varchar(255);unique;not null"`
}
// ---------------对聊天记录的操作-----------------------
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
// 删除指定用户的聊天记录硬删除
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
// --------------------------------------

// ----------------对用户表进行的数据库操作----------------------

// 好像只有在用户注册的时候才会使用默认值
// BeforeCreate 在创建记录之前设置默认值
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    if u.Username == "" {
        u.Username = "default_username"
    }
    return nil
}
func  GetUser(db *gorm.DB, userID string)(*User,error){
	logrus.Info("开始获取用户信息")
	// 初始化user对象
	var user User
	// 根据id查找用户
	result := db.First(&user,"id = ?",userID)
	// 检查是否找到用户
	if errors.Is(result.Error,gorm.ErrRecordNotFound){
		return nil,errors.New("用户不存在")
	}
	if result.Error != nil{
		// return nil, errors.Wrap(result.Error, "查询用户信息时发生错误")
	}
	logrus.Info("用户获取信息成功")
	return &user,nil
}
func  Register(db *gorm.DB, userID string){

}
func  UpdarerUser(db *gorm.DB, userID string){

}
func  DeleteUser(db *gorm.DB, userID string){

}

// --------------------------------------