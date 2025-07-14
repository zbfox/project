package main

import (
	"TestGin/api"
	"TestGin/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin引擎
	r := gin.Default()
	config.InitConfig()
	config.InitDB()
	config.InitRedis()

	//fmt.Println("MySQL Host:", config.Conf.MySQL.Host)
	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	r.Run(":8080")
}
