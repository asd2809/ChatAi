package repository

import (
	"fmt"
	"smart-dialog-ai/internal/model"
	"smart-dialog-ai/internal/pkg"

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
func  GetUser(db *gorm.DB,userID string){
	
}
// 需要传入db,需要传入web的结构体,
func  Register(db *gorm.DB,reg model.RegisterRequest) error{
	logrus.Info("开始新建用户信息")
	// 检查前端传入数据的逻辑
	var count int64
	db.Model(&model.User{}).Where("username = ?",reg.Username).Count(&count)
	if count > 0 {
		return pkg.NewBizError(pkg.CodeUserExists)
	}
	// 检查邮箱是否已存在
	db.Model(&User{}).Where("email = ?",reg.Email).Count(&count)
	if count > 0{
		return pkg.NewBizError(pkg.CodeEmailExists)
	}
	// 创建用户
	user := model.User{
		Username: reg.Username,
		Password: reg.Password,
		Email: reg.Email,
	}
	// 向表中添加数据
	if err := db.Create(&user).Error; err !=nil{
		return pkg.NewBizError(pkg.CodeDBError)
	}
	logrus.Info("添加用户信息成功")
	return nil
}
func  UpdarerUser(db *gorm.DB, userID string){

}
func  DeleteUser(db *gorm.DB, userID string){

}

// --------------------------------------