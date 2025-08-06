package model

import (
	ti "TestGin/util"
)

type UserResponse struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

// UserToResponse 用于将 User 转换为 UserResponse，时间字段格式化为字符串
func UserToResponse(u User) UserResponse {
	var deletedAt string
	if u.DeletedAt.Valid {
		deletedAt = ti.FormatTime(u.DeletedAt.Time)
	} else {
		deletedAt = ""
	}
	return UserResponse{
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: ti.FormatTime(u.CreatedAt),
		UpdatedAt: ti.FormatTime(u.UpdatedAt),
		DeletedAt: deletedAt,
	}
}
