package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ChatRecord struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Role      string    `gorm:"type:enum('user','assistant')"`
	Content   string
	CreatedAt time.Time
}

func main() {
	// 数据库连接配置
	dsn := "root:root@cbj@tcp(127.0.0.1:3306)/chat_ai?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("无法连接数据库: " + err.Error())
	}        

	// 自动迁移表结构（如果表不存在则创建）
	err = db.AutoMigrate(&ChatRecord{})
	if err != nil {
		panic("迁移表结构失败: " + err.Error())
	}

	// 示例1: 获取单条记录（根据ID）
	var record ChatRecord
	result := db.First(&record) // 获取ID=1的记录
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("未找到记录")
		} else {
			fmt.Println("查询失败:", result.Error)
		}
	} else {
		fmt.Printf("找到记录: %+v\n", record)
	}

	// 示例2: 获取多条记录（根据用户ID）
	var records []ChatRecord
	result = db.Where("user_id = ?", "user1").Order("created_at desc").Find(&records)
	if result.Error != nil {
		fmt.Println("查询失败:", result.Error)
	} else {
		fmt.Printf("找到%d条记录:\n", len(records))
		for _, r := range records {
			fmt.Printf("%+v\n", r)
		}
	}

	// 示例3: 获取表的所有记录
	var allRecords []ChatRecord
	result = db.Find(&allRecords)
	if result.Error != nil {
		fmt.Println("查询所有记录失败:", result.Error)
	} else {
		fmt.Printf("表中共有%d条记录\n", len(allRecords))
	}
}
