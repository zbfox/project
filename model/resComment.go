package model

import (
	ti "TestGin/util"
	"encoding/json"
	"time"
)

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint     `json:"id"`
	PostID    uint     `json:"post_id"`
	UserID    uint     `json:"user_id"`
	Content   string   `json:"content"`
	ParentID  uint     `json:"parent_id"`
	Type      uint8    `json:"type"`
	URLs      string   `json:"urls"`
	Resources []string `json:"resources"`
	CreatedAt string   `json:"created_at"`
}

// CommentToResponse  评论转为响应
func CommentToResponse(comments []CommentResponse) []CommentResponse {
	var commentRes []CommentResponse
	for _, c := range comments {
		ctime, _ := time.Parse(time.RFC3339, c.CreatedAt)
		resp := CommentResponse{
			ID:       c.ID,
			PostID:   c.PostID,
			UserID:   c.UserID,
			Content:  c.Content,
			ParentID: c.ParentID,
			Type:     uint8(c.Type),
			//转为时间 time
			CreatedAt: ti.FormatTime(ctime),
		}

		if c.ID != 0 && c.URLs != "" {
			var urls []string
			err := json.Unmarshal([]byte(c.URLs), &urls)
			if err != nil {
				break
			}
			resp.Resources = urls

		}
		commentRes = append(commentRes, resp)
	}
	return commentRes
}
