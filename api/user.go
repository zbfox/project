package api

import (
	db "TestGin/config"
	res "TestGin/middleware"
	"TestGin/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Tags 用户
// @Produce json
// @Success 200 {object} string "成功"
// @Param id path int true "用户ID"
// @Router /api/user/get/:id [get]
func GetUser(c *gin.Context) {
	var u model.User

	// id, _ := strconv.Atoi(c.Param("id"))
	uuid := c.Param("id")
	redisClient := db.GetRedisClient()
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", uuid)
	userJson, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(userJson), &u); err == nil {
			us := model.UserToResponse(u)
			c.JSON(200, gin.H{
				"message": "user (from cache)",
				"data":    us,
			})
			return
		}
	}
	// 缓存未命中，查数据库
	result := db.DB.Where("uuid = ?", uuid).First(&u)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "用户不存在",
			"id":      uuid,
		})
		return
	}
	// 查到后写入缓存
	userBytes, _ := json.Marshal(u)
	redisClient.Set(ctx, cacheKey, userBytes, 0)

	us := model.UserToResponse(u)
	log.Printf("user: %+v\n", us)
	c.JSON(200, gin.H{
		"message": "user",
		"data":    us,
	})
}

// AddUser 添加用户
// @Summary 添加用户
// @Tags 用户
// @Produce json
// @Success 200 {object} string "成功"
// @Router /api/user/add [POST]
func AddUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(400, gin.H{
			"message": "参数绑定失败",
			"error":   err.Error(),
		})
		return
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "用户添加失败",
			"error":   err.Error(),
		})
		return
	}
	redisClient := db.GetRedisClient()
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	userBytes, _ := json.Marshal(user)
	redisClient.Set(ctx, cacheKey, userBytes, 0)

	c.JSON(200, gin.H{
		"message": "用户添加成功",
		"data":    user,
	})
}

// ListUsers 用户列表
// @Summary 获取用户列表
// @Tags 用户
// @Description 返回所有用户信息
// @Produce json
// @Success 200 {array} model.UserResponse "用户列表"
// @Router /api/user/list [GET]
func ListUsers(c *gin.Context) {
	var users []model.User

	// 查询用户数据
	if err := db.DB.Unscoped().Find(&users).Error; err != nil {
		res.Error(c, http.StatusInternalServerError, err)
		return
	}

	// 转换为响应结构
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserToResponse(user)
	}
	res.Success(c, userResponses)
}

// UpdateUser 更新用户数据
// @Summary 更新用户数据
// @Tags 用户
// @Description 更新用户数据
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} middleware.Response "成功"
// @Router /api/user/update/user [POST]
func UpdateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res.Error(c, http.StatusBadRequest, err)
		return
	}
	log.Printf("%+v\n", user)
	res.Success(c, "更新用户成功")
}

// UpdatePassword 更新用户密码
// @Summary 更新用户密码
// @Description 更新用户密码
// @Produce json
// @Tags 用户
// @Param user body model.User true "用户信息"
// @Success 200 {object} middleware.Response "成功"
// @Router /api/user/update/password [POST]
func UpdatePassword(c *gin.Context) {
	res.Success(c, "更新密码成功")
}
