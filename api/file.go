package api

import (
	res "TestGin/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

// 文件上传
func uploadFile(c *gin.Context) {
	log.Println("上传文件")

	res.Success(c, "上传成功")
}
