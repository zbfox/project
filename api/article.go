package api

import (
	db "TestGin/config"
	res "TestGin/middleware"
	"TestGin/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// DeleteArticle 删除文章
// @Summary 删除文章
// @Tags 文章
// @Description 删除文章
// @Param id path int true "文章ID"
// @Success 200 {object} string  "文章信息"
// @Router /api/articles/delete/:id [delete]
func DeleteArticle(c *gin.Context) {
	// 获取路径参数 id
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := db.DB.Delete(model.Article{}, id)
	log.Printf("结果：%v\n", ctx)

	// TODO: 实现删除逻辑，例如通过数据库操作删除对应文章
	res.Success(c, "文章已删除")
}

// GetArticle 查询文章
// @Summary 查询文章
// @Description 查询文章
// @Tags 文章
// @Param id path int true "文章ID"
// @Success 200 {object} model.ArticleResponse  "文章信息"
// @Router /api/article/get/:id [get]
func GetArticle(c *gin.Context) {
	var article model.Article
	id, _ := strconv.Atoi(c.Param("id"))
	// 查询文章
	if err := db.DB.Table("articles").
		Where("id = ? AND status = 0", id).
		First(&article).Error; err != nil {
		// 处理不同类型的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Error(c, 404, err)
		} else {
			res.Error(c, 400, err)
		}
		return
	}
	// 转换为响应格式并返回
	ar := model.ArticleToResponse(article)
	res.Success(c, ar)
}

// UpdateArticleStatus 更新文章状态
// @Summary 更新文章状态
// @Description 更新文章状态
// @Tags 文章
// @Param id path int true "文章ID"
// @Param request body model.ArticleStatus true "请求体 (status: 0=Draft, 1=Pending, 2=Published)"
// @Success 200 {object} middleware.Response "更新成功返回"
// @Router /api/article/{id}/status [put]
func UpdateArticleStatus(c *gin.Context) {
	var article model.Article
	article.Status = model.Draft
	aid, _ := strconv.Atoi(c.Param("id"))
	//打印请求体
	sta := c.ShouldBindJSON(&article.Status)
	log.Printf("id:%v----status:%s\n", aid, sta)
	res.Success(c, "操作成功")
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 更新文章
// @Tags 文章
// @Param id path int true "文章ID"
// @Param request body model.Article true "请求体"
// @Success 200 {object} middleware.Response "更新成功返回"
// @Router /api/article/{id} [put]
func UpdateArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		res.Error(c, 400, err)
		return
	}
	log.Printf("%+v\n", article)
	res.Success(c, "更新文章成功")
}

// AddArticle 添加文章
// @Summary 添加文章
// @Description 添加文章
// @Tags 文章
// @Param request body model.Article true "请求体"
// @Success 200 {object} middleware.Response "添加成功返回"
// @Router /api/article/add [post]
func AddArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		res.Error(c, 400, err)
		return
	}

	log.Printf("%+v\n", article)
	if err := db.DB.Create(&article).Error; err != nil {
		res.Error(c, 500, err)
		return
	}
	res.Success(c, "添加文章成功")
}
