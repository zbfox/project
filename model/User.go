package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string         `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`                        // 用户唯一标识
	Username  string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"username" binding:"required"` // 用户名，唯一
	Password  string         `gorm:"type:varchar(100);not null" json:"password"  binding:"required"`           // 密码
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email" binding:"omitempty,email"`     // 邮箱，唯一
	Phone     string         `gorm:"type:varchar(20)" json:"phone" binding:"omitempty"`                        // 手机号
	Role      string         `gorm:"type:varchar(20);default:user" json:"role"`                                // 角色
	Status    string         `gorm:"type:varchar(20);default:active" json:"status"`                            // 状态 active/disabled
	CreatedAt time.Time      `json:"created_at"`                                                               // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                                               // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                           // 软删除                                                       // 软删除
}

// AutoMigrate 创建或更新表结构
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic("数据库自动迁移失败: " + err.Error())
	}
}

// AfterCreate 插入之后执行
func (u *User) AfterCreate(db *gorm.DB) (err error) {
	fmt.Println("AfterCreate")
	fmt.Printf("")
	return nil
}

// AfterFind 查询结束后执行
func (u *User) AfterFind(db *gorm.DB) (err error) {
	fmt.Println("AfterFind", u)
	return nil
}

// BeforeCreate 创建之前执行
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return
}
