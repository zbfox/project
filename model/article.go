package model

import (
	ti "TestGin/util"
	"time"

	"gorm.io/gorm"
)

// ArticleStatus 文章状态枚举
type ArticleStatus int

const (
	Draft     ArticleStatus = iota // 草稿
	Pending                        // 待审核
	Published                      // 已发布
)

// Article 表结构映射 (articles)
type Article struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`                  // 主键
	UserID     int64          `gorm:"not null;index" json:"user_id"`                       // 外键关联用户ID
	Title      string         `gorm:"type:varchar(200);not null" json:"title"`             // 标题
	Content    string         `gorm:"type:text;not null" json:"content"`                   // 内容
	Status     ArticleStatus  `gorm:"default: 0" json:"status" enums:"0,1,2"`              // 状态
	StatusName string         `gorm:"type:varchar(8);default:'draft'" json:"status_name" ` // 状态名称
	CreatedAt  time.Time      `json:"created_at"`                                          // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"`                                          // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`                                      // 软删除
}

// AfterFind 插入前
func (a *Article) AfterFind(db *gorm.DB) (err error) {
	switch a.Status {
	case 1:
		a.StatusName = "draft"
	case 2:
		a.StatusName = "pending"
	default:
		a.StatusName = "published"
	}
	return
}

// ArticleResponse 响应结构体
type ArticleResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ArticleToResponse 将 Article 转换为 ArticleResponse
func ArticleToResponse(a Article) ArticleResponse {
	return ArticleResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		Title:      a.Title,
		Content:    a.Content,
		Status:     int(a.Status),
		StatusName: a.StatusName,
		CreatedAt:  ti.FormatTime(a.CreatedAt),
		UpdatedAt:  ti.FormatTime(a.UpdatedAt),
	}
}

// AutoMigrateArticle 创建或更新 Article 表结构
func AutoMigrateArticle(db *gorm.DB) {
	err := db.AutoMigrate(&Article{})
	if err != nil {
		panic("Article 表自动迁移失败: " + err.Error())
	}
}
