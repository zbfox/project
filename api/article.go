package api

import (
	sql "TestGin/config"
	"TestGin/model"
	"github.com/gin-gonic/gin"
)

// DeleteArticle 删除文章
// @Summary 删除文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	// 获取路径参数 id
	id := c.Param("id")

	// TODO: 实现删除逻辑，例如通过数据库操作删除对应文章

	c.JSON(200, gin.H{
		"message": "文章已删除",
		"id":      id,
	})
}

// GetArticle 查询文章
func GetArticle(c *gin.Context) {
	//创建
	println("GetArticle")
	id := c.Param("id")
	s := sql.DB.Table("articles").Where("id = ? AND status = 1", id).Select("id, title, content").First(&model.Article{})
	c.JSON(200, gin.H{
		"message": "查询成功",
		"data":    s,
	})
}
