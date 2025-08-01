package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论主表
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index;comment:所属帖子ID" json:"post_id"`
	UserID    uint      `gorm:"not null;index;comment:评论用户" json:"user_id"`
	Content   string    `gorm:"type:text;comment:评论内容" json:"content"`
	ParentID  *uint     `gorm:"index;comment:父评论ID" json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	// 关联
	Resources Resource `gorm:"foreignKey:CommentID" json:"resources,omitempty"`
}

// ResourceType 资源类型
type ResourceType uint8

const (
	ResourceTypeImage ResourceType = iota // 图片
	ResourceTypeVideo                     // 视频
)

// Resource 资源表
type Resource struct {
	ID        uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	CommentID uint         `gorm:"not null;index;comment:关联的评论ID" json:"comment_id"`
	Type      ResourceType `gorm:"type:tinyint;not null;comment:资源类型" json:"type"`
	URLs      string       `gorm:"type:json;not null;comment:资源URL列表" json:"urls"`
	CreatedAt time.Time    `json:"created_at"`
}

// CommentRequest 创建评论请求
type CommentRequest struct {
	PostID   uint         `json:"post_id" binding:"required"`
	UserID   uint         `json:"user_id" binding:"required"`
	Content  string       `json:"content" binding:"required"`
	ParentID *uint        `json:"parent_id"`
	Type     ResourceType `json:"type" binding:"required_with=Files"`
	Files    []string     `json:"files"`
}

// AutoMigrateComment 数据库迁移
func AutoMigrateComment(db *gorm.DB) error {
	return db.AutoMigrate(&Comment{}, &Resource{})
}
