package main

import (
	"TestGin/api"
	"TestGin/config"
	"TestGin/middleware"
	"TestGin/util"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin引擎
	r := gin.Default()
	config.InitConfig()
	rdb := config.InitRedis()
	middleware.InitJWTMiddleware(rdb)
	config.InitDB()
	util.InitWebsocket(r)

	//fmt.Println("MySQL Host:", config.Conf.MySQL.Host)
	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
