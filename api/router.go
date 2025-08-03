package api

import (
	"TestGin/config"
	_ "TestGin/docs"
	"TestGin/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine) {
	red := config.GetRedisClient()
	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/static", "./static")

	// API CURD 分组
	v1 := r.Group("/api")
	article := v1.Group("/article")
	{
		article.POST("/add", AddArticle)
		article.PUT("/update/:id", UpdateArticle)
		//更新文章状态
		article.PUT("/:uuid/status", UpdateArticleStatus)
		article.DELETE("/delete/:id", DeleteArticle)
		//article.GET("/list", ListArticle)
		article.GET("/get/:id", middleware.RedisCacheMiddleware(middleware.CacheOptions{RedisClient: red, TTL: 60 * time.Second}, GetArticle))
	}
	user := v1.Group("/user")
	{
		user.POST("/add", AddUser)
		user.GET("/list", middleware.RedisCacheMiddleware(middleware.CacheOptions{RedisClient: red, TTL: 60 * time.Second}, ListUsers))
		user.GET("/get/:id", middleware.RedisCacheMiddleware(middleware.CacheOptions{RedisClient: red, TTL: 60 * time.Second}, GetUser))
		update := user.Group("/update")
		{
			update.POST("/password", UpdatePassword)
			update.POST("/user", UpdateUser)
		}
	}
	//file := v1.Group("/upload")
	//{
	//	//file.POST("/resources")
	//}
	comment := v1.Group("/comment")
	{
		comment.POST("/add", AddComment)
		comment.GET("/list", middleware.RedisCacheMiddleware(middleware.CacheOptions{RedisClient: red, TTL: 60 * time.Second}, ListComments))
	}

	// 其他路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.GET("/", func(c *gin.Context) {
		//	重定向
		c.Redirect(302, "https://www.baidu.com")
	})
}
