package model

import "time"

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}

// UserToResponse 用于将 User 转换为 UserResponse，时间字段格式化为字符串
func UserToResponse(u User) UserResponse {
	var deletedAt string
	if u.DeletedAt.Valid {
		deletedAt = FormatTime(u.DeletedAt.Time)
	} else {
		deletedAt = ""
	}
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: FormatTime(u.CreatedAt),
		UpdatedAt: FormatTime(u.UpdatedAt),
		DeletedAt: deletedAt,
	}
}
