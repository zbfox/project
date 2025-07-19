package api

import (
	db "TestGin/config"
	"TestGin/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Produce json
// @Success 200 {object} string "成功"
// @Router /api/user/:id [get]
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
func ListUsers(c *gin.Context) {
	var users []model.User

	// 查询用户数据
	if err := db.DB.Unscoped().Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应结构
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserToResponse(user)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    userResponses,
	})
}

// UpdateUser 更新用户数据
func UpdateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "参数绑定失败",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("%+v\n", user)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "update user",
	})
}

// UpdatePassword 更新用户密码
func UpdatePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "update user",
	})
}
