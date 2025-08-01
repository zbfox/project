package util

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// FileTypeRule 定义文件类型白名单规则
type FileTypeRule struct {
	AllowedMimePrefixes []string // 允许的 MIME 前缀，例如 image/, video/
	AllowedExtensions   []string // 允许的文件后缀，例如 .jpg, .png
}

// FileType 返回的文件类型
type FileType struct {
	MimeType  string
	Extension string
}

// ValidateFileType 校验文件类型
func ValidateFileType(fileHeader *multipart.FileHeader, rule FileTypeRule) (FileType, error) {
	// 1. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return FileType{"", ""}, err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 2. 读取前 512 字节，检测 MIME
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return FileType{"", ""}, err
	}

	mimeType := http.DetectContentType(buffer[:n])
	//log.Printf("mimeType:%s\n", mimeType)

	// 3. 判断 MIME 类型
	allowedMime := false
	for _, prefix := range rule.AllowedMimePrefixes {
		if strings.HasPrefix(mimeType, prefix) {
			allowedMime = true
			break
		}
	}
	if !allowedMime {
		return FileType{mimeType, ""}, errors.New("不支持的文件 MIME 类型: " + mimeType)
	}

	// 4. 判断文件后缀
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExt := false
	for _, e := range rule.AllowedExtensions {
		if ext == e {
			allowedExt = true
			break
		}
	}
	if !allowedExt {
		return FileType{mimeType, ext}, errors.New("不支持的文件扩展名: " + ext)
	}
	ext, err = MimeToExtension(mimeType)
	return FileType{mimeType, ext}, err
}

// MimeToExtension 根据 MIME Type 返回标准文件后缀
func MimeToExtension(mimeType string) (string, error) {
	// 去除 MIME 参数，如 image/jpeg; charset=utf-8
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = mimeType[:idx]
	}

	mimeMap := map[string]string{
		"image/jpeg":                   ".jpg",
		"image/png":                    ".png",
		"image/gif":                    ".gif",
		"image/webp":                   ".webp",
		"image/bmp":                    ".bmp",
		"video/mp4":                    ".mp4",
		"video/quicktime":              ".mov",
		"application/pdf":              ".pdf",
		"application/zip":              ".zip",
		"application/x-rar-compressed": ".rar",
		"text/plain":                   ".txt",
		"application/msword":           ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		// 可根据需要扩展更多类型
	}

	if ext, ok := mimeMap[mimeType]; ok {
		return ext, nil
	}

	return "", errors.New("不支持的 MIME 类型: " + mimeType)
}
