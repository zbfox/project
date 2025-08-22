package api

import (
	db "TestGin/config"
	res "TestGin/middleware"
	"TestGin/model"
	"TestGin/util"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type UploadItems struct {
	FileName string
	File     *multipart.FileHeader
}

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
// @Param files formData file true "上传文件"
// @Param   Authorization  header  string  true  "Bearer Token"
// @Success 200 {object} res.Response "{"code":200,"data":{},"msg":"操作成功"}"
// @Router /api/comment/add [post]
func AddComment(c *gin.Context) {

	content := c.PostForm("content")
	typeFile, _ := strconv.Atoi(c.PostForm("type"))

	postIDInt, _ := strconv.Atoi(c.PostForm("postId"))
	userIDInt, _ := strconv.Atoi(c.PostForm("userId"))
	parentIDInt, _ := strconv.Atoi(c.PostForm("parentId"))
	parentIDTo := uint(parentIDInt)
	parentIDPtr := &parentIDTo // parentIDPtr 类型为 *uint

	formFile, err := c.MultipartForm()
	files := formFile.File["files"]
	uploadPath := "static"

	if err != nil {
		res.Error(c, http.StatusInternalServerError, err)
		return
	}
	var allowedTypes = []string{"image/", "video/"}
	var allowedExt = []string{".png", ".jpg", ".jpeg", ".gif", ".mp4"}
	var fileList []string
	var imageCount, videoCount int

	var filePaths []UploadItems
	if formFile.File == nil {
		res.Error(c, http.StatusBadRequest, errors.New("请选择文件"))
		return
	}
	// 检查 MIME 类型
	for _, file := range files {
		fileType, err := util.ValidateFileType(file, util.FileTypeRule{AllowedMimePrefixes: allowedTypes, AllowedExtensions: allowedExt})
		if err != nil {
			res.Error(c, http.StatusInternalServerError, err)
			return
		}
		unix, _ := uuid.NewRandom()
		fileName := filepath.Join("image", time.Now().Format("20060102"), unix.String()+fileType.Extension)
		saveFilePath := filepath.Join(uploadPath, fileName)

		if strings.HasPrefix(fileType.MimeType, "image/") {
			imageCount++
		} else if strings.HasPrefix(fileType.MimeType, "video/") {
			videoCount++
		}
		log.Printf("file：%v\n", saveFilePath)
		filePaths = append(filePaths, UploadItems{
			File:     file,
			FileName: saveFilePath,
		})

	}

	// 不允许混合上传
	if imageCount > 0 && videoCount > 0 {
		res.Error(c, http.StatusBadRequest, errors.New("不能同时上传图片和视频"))
		return
	}
	if imageCount > 9 {
		res.Error(c, http.StatusBadRequest, errors.New("最多上传 9 张图片"))
		return
	}
	if videoCount > 1 {
		res.Error(c, http.StatusBadRequest, errors.New("最多上传 1 个视频"))
		return
	}

	for _, file := range filePaths {
		if err := c.SaveUploadedFile(file.File, file.FileName); err != nil {
			res.Error(c, http.StatusInternalServerError, err)
			return
		}
		// 统一为URL风格路径，存入数据库
		urlStyle := strings.ReplaceAll(file.FileName, "\\", "/")
		log.Printf("保存成功: %s", urlStyle)
		fileList = append(fileList, urlStyle)
	}

	form := model.Comment{
		Content:  content,
		PostID:   uint(postIDInt),
		UserID:   uint(userIDInt),
		ParentID: parentIDPtr,
	}

	// 将 fileList 序列化为 JSON 字符串
	urlJSON, err := json.Marshal(fileList)
	if err != nil {
		log.Printf("JSON 序列化失败: %v", err)
		res.Error(c, 500, err)
		return
	}

	tx := db.DB.Begin()
	if err1 := tx.Create(&form).Error; err1 != nil {
		tx.Rollback()
		res.Error(c, 500, err1)
		return
	}

	resource := model.Resource{
		CommentID: form.ID,
		Type:      model.ResourceType(typeFile),
		URLs:      string(urlJSON),
	}
	if err2 := tx.Debug().Create(&resource).Error; err2 != nil {
		tx.Rollback()
		res.Error(c, 500, err2)
		return
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		res.Error(c, 500, commitErr)
		return
	}
	res.Success(c, "")
}

// ListComments 获取评论列表
// @Summary 获取评论列表
// @Description 获取评论列表
// @Tags 评论
// @Param postId query int true "帖子ID"
// @Param   Authorization  header  string  true  "Bearer Token"
// @Success 200 {object} res.Response "{"code":200,"data":{},"msg":"操作成功"}"
// @Router /api/comment/list [get]
func ListComments(c *gin.Context) {
	postID := c.Query("postId")
	postIDInt, _ := strconv.Atoi(postID)

	var results []model.CommentResponse

	err := db.DB.
		Table("comments AS c").
		Select(`
    c.id, 
    c.content, 
    c.post_id, 
    c.user_id, 
    c.parent_id, 
    r.type, 
    r.urls, 
    c.created_at,
    u.username
  `).
		Joins("LEFT JOIN resources AS r ON r.comment_id = c.id ").
		Joins("LEFT JOIN users AS u ON u.id = c.user_id ").
		Where("c.post_id = ?", postIDInt).
		Find(&results).
		Error
	if err != nil {
		res.Error(c, 500, err)
		return
	}
	//log.Printf("results:%+v\n", results)
	resultsRes := model.CommentToResponse(results)
	//评论重组
	resultsRes, _ = model.CommentsReorganize(resultsRes)
	res.Success(c, resultsRes)
}
