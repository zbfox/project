package api

import (
	res "TestGin/middleware"
	"TestGin/model"
	"TestGin/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
// @Success 200 {object} res.Response "{"code":200,"data":{},"msg":"操作成功"}"
// @Router /api/comment/add [post]
func AddComment(c *gin.Context) {

	content := c.PostForm("content")
	postID := c.PostForm("postId")
	userID := c.PostForm("userId")
	parentID := c.PostForm("parentId")
	//typeFile := c.PostForm("type")/

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
	log.Printf("form:%+v\n", form)
	formFile, err := c.MultipartForm()
	files := formFile.File["files"]
	uploadPath := "./static"

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
		// 生成目标文件名 128位的uuid
		newV6, err := uuid.NewV6()
		fileName := "/image/" + time.Now().Format("2006010215040") + newV6.String() + fileType.Extension
		saveFilePath := filepath.Join(uploadPath, fileName)

		if strings.HasPrefix(fileType.MimeType, "image/") {
			imageCount++
		} else if strings.HasPrefix(fileType.MimeType, "video/") {
			videoCount++
		} else {
			res.Error(c, http.StatusBadRequest, fmt.Errorf("不支持的文件类型: %s", fileType.MimeType))
			return
		}
		log.Printf("imageCount:%v,videoCount:%v\n", imageCount, videoCount)
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
		log.Printf("上传成功: %s", file.FileName)
		fileList = append(fileList, file.FileName)
	}
	res.Success(c, fileList)
}
