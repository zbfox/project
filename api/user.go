package api

import (
	db "TestGin/config"
	"TestGin/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Produce json
// @Success 200 {object} string "成功"
// @Router /api/user/:id [get]
func GetUser(c *gin.Context) {
	var u model.User

	id, _ := strconv.Atoi(c.Param("id"))
	result := db.DB.Where("id = ?", id).First(&u)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "用户不存在",
			"id":      id,
		})
		return
	}
	us := model.UserToResponse(u)
	fmt.Printf("user: %+v\n", us)
	c.JSON(200, gin.H{
		"message": "user",
		"data":    us,
	})
}

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("接收到用户信息: %+v\n", user)
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
	//var user model.User
	id, _ := strconv.Atoi(c.Param("id"))
	//反射输出数据类型
	log.Println("id: ", reflect.TypeOf(id))

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "update user",
	})
}
