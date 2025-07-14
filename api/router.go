package api

import (
	_ "TestGin/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine) {
	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API CURD 分组
	v1 := r.Group("/api/")
	article := v1.Group("/article")
	{
		//article.POST("/add", AddArticle)
		//article.PUT("/update/:id", UpdateArticle)
		article.DELETE("/delete/:id", DeleteArticle)
		//article.GET("/list", ListArticle)
		article.GET("/:id", GetArticle)
	}
	user := v1.Group("/user")
	{
		user.POST("/add", AddUser)
		//user.PUT("/update/:id", UpdateUser)
		//user.DELETE("/delete/:id", DeleteUser)
		user.GET("/list", ListUsers)
		user.GET("/:id", GetUser)
	}

	// 其他路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
