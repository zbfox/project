package config

import (
	"TestGin/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func InitDB() {
	c := Conf.MySQL

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 配置GORM以仅打印API执行的SQL语句
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	log.Println("数据库连接成功")
	model.AutoMigrate(db) // 创建用户表结构
	model.AutoMigrateArticle(db)
	err = model.AutoMigrateComment(db)
	if err != nil {
		return
	}
	model.AutoMigrateEmoji(db) // 创建表情包表结构
	DB = db
}
