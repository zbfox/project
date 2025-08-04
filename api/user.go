package api

import (
	db "TestGin/config"
	res "TestGin/middleware"
	"TestGin/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var redisClient = db.GetRedisClient()

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

	// 缓存未命中，查数据库
	result := db.DB.Where("uuid = ?", uuid).First(&u)
	if result.RowsAffected == 0 {
		res.Error(c, http.StatusNotFound, errors.New("用户不存在"))
		return
	}
	us := model.UserToResponse(u)
	res.Success(c, us)
	return
}

// Register  添加用户
// @Summary 添加用户
// @Tags 用户
// @Produce json
// @Param user body model.User true "用户信息"
// @Success 200 {object} string "成功"
// @Router /api/user/add [POST]
func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "参数绑定失败",
			"error":   err.Error(),
		})
		return
	}
	if err := db.DB.Where("email = ?", user.Email).Table("users").Error; err != nil {
		res.Error(c, http.StatusBadRequest, errors.New("邮箱已被占用"))
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
		"data":    "",
	})
	return
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
	//开启事务
	tx := db.DB.Begin()
	if err := c.ShouldBindJSON(&user); err != nil {
		res.Error(c, http.StatusBadRequest, err)
		return
	}
	err := tx.Where("id = ?", user.UUID).Updates(&user).Error
	if err != nil {
		//回滚事务
		tx.Rollback()
		res.Error(c, 500, err)
		return
	}
	tx.Commit()
	//同步更新到Redis
	ctx := context.Background()
	err = redisClient.Set(ctx, fmt.Sprintf("user:%s", user.UUID), user, 0).Err()
	if err != nil {
		res.Error(c, 500, err)
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
	//加强修改密码的流程
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res.Error(c, http.StatusBadRequest, err)
		return
	}
	//if user.Password != user.NewPassword {
	//	res.Error(c, http.StatusBadRequest, errors.New("新密码和旧密码不一致"))
	//	return
	//}
	if err := db.DB.Where("id = ?", user.UUID).Updates(&user).Error; err != nil {
		res.Error(c, 500, err)
		return
	}

	res.Success(c, "更新密码成功")
}

// Login 登录
// @Summary 登录
// @Description 登录
// @Produce json
// @Tags 用户
// @Param user body model.User true "用户信息"
// @Success 200 {object} middleware.Response "成功"
// @Router /api/user/login [POST]
func Login(c *gin.Context) {
}
