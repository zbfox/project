package model

import (
	"time"

	"gorm.io/gorm"
)

// ArticleStatus 文章状态枚举
type ArticleStatus string

const (
	Draft     ArticleStatus = "draft"
	Published ArticleStatus = "published"
)

// Article 表结构映射 (articles)
type Article struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`                           // 主键
	UserID    int64          `gorm:"not null;index" json:"user_id"`                                // 外键关联用户ID
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`                      // 标题
	Content   string         `gorm:"type:text;not null" json:"content"`                            // 内容
	Status    ArticleStatus  `gorm:"type:enum('draft','published');default:'draft'" json:"status"` // 状态
	CreatedAt time.Time      `json:"created_at"`                                                   // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                                   // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                               // 软删除
}

// AutoMigrateArticle 创建或更新 Article 表结构
func AutoMigrateArticle(db *gorm.DB) {
	err := db.AutoMigrate(&Article{})
	if err != nil {
		panic("Article 表自动迁移失败: " + err.Error())
	}
}
