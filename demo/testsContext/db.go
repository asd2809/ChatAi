package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "root:root@cbj@tcp(127.0.0.1:3306)/chat_ai?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ 数据库连接失败: " + err.Error())
	}

	// 自动建表
	if err := db.AutoMigrate(&ChatRecord{}); err != nil {
		panic("❌ 数据库迁移失败: " + err.Error())
	}

	return db
}
