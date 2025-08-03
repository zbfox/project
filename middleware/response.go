package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Code    int         `json:"code"`    // 自定义状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据
	Time    string      `json:"time"`
}

// Success 成功返回 生成swagger提示
// @Summary 成功返回
// @Description 成功返回
// @Tags 公共
// @Accept json
// @Produce json
// @Param data body interface{} true "返回数据"
// @Success 200 {object} middleware.Response "{"code":0,"message":"success"}"
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
		Time:    strconv.FormatInt(time.Now().Unix(), 11),
	})
}

// Error httpCode int不是必要参数
func Error(c *gin.Context, code int, err error) {
	//不写http_code就默认使用200 AbortWithStatusJSON中断请求并且返回json数据，不会执行后续处理
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": err.Error(),
		"data":    nil,
		"time":    strconv.FormatInt(time.Now().Unix(), 11),
	})
}
