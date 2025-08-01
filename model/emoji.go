package model

import (
	"gorm.io/gorm"
	"time"
)

// Emoji  表
type Emoji struct {
	ID        uint      `json:"id" gorm:"primaryKey"`       // 主键
	Name      string    `gorm:"not null;index" json:"name"` //
	Url       string    `gorm:"not null;index" json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

func AutoMigrateEmoji(db *gorm.DB) {
	err := db.AutoMigrate(&Emoji{})
	if err != nil {
		return
	}
}
