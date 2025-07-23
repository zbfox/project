package api

import (
	res "TestGin/middleware"
	"TestGin/model"
	"TestGin/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// AddComment 评论
// @Summary 添加评论
// @Description 添加评论
// @Tags 评论
// @Accept  multipart/form-data
// @Param content formData string true "评论内容"
// @Param postId formData int true "所属帖子ID"
// @Param userId formData int true "评论用户ID"
// @Param parentId formData int false "父评论ID"
// @Param type formData int true "资源类型 image/video"
// @Param files formData files true "上传文件"
// @Success 200 {object} res.Response "{"code":200,"data":{},"msg":"操作成功"}"
// @Router /api/comment/add [post]
func AddComment(c *gin.Context) {

	content := c.PostForm("content")
	postID := c.PostForm("postId")
	userID := c.PostForm("userId")
	parentID := c.PostForm("parentId")

	postIDInt, _ := strconv.Atoi(postID)
	userIDInt, _ := strconv.Atoi(userID)
	parentIDInt, _ := strconv.Atoi(parentID)
	parentIDTo := uint(parentIDInt)
	parentIDPtr := &parentIDTo // parentIDPtr 类型为 *uint

	form := model.CommentRequest{
		Content:  content,
		PostID:   uint(postIDInt),
		UserID:   uint(userIDInt),
		ParentID: parentIDPtr,
		Type:     1,
	}
	log.Printf("content:%s,postID:%s,userID:%s,parentID:%s", content, postID, userID, parentID)
	formFile, err := c.MultipartForm()
	files := formFile.File["files"]
	uploadPath := "./static"

	if err != nil {
		res.Error(c, http.StatusInternalServerError, err)
		return
	}
	var allowedTypes = []string{"image", "video/"}
	var allowedExt = []string{".png", ".jpg", ".jpeg", ".gif", ".mp4"}
	if formFile.File != nil {
		for _, file := range files {
			if len(files) > 9 {
				err := errors.New("最多上传9张图片")
				res.Error(c, http.StatusInternalServerError, err)
			}
			fileType := util.FileType{}
			fileType, err := util.ValidateFileType(file, util.FileTypeRule{AllowedMimePrefixes: allowedTypes, AllowedExtensions: allowedExt})
			if err != nil {
				res.Error(c, http.StatusInternalServerError, err)
				return
			}
			// 生成目标文件名
			fileName := "/image/" + time.Now().Format("20060102150405") + fileType.Extension
			filePath := filepath.Join(uploadPath, fileName)

			if err := c.SaveUploadedFile(file, filePath); err != nil {
				res.Error(c, http.StatusInternalServerError, err)
				return
			}
		}
	}
	res.Success(c, "操作成功")
}
