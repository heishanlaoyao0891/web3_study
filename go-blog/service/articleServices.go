package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex 首页
func GetIndex(c *gin.Context) {
	// 跳转到首页模板（第三步完善）
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Go博客首页",
	})
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	var articles []model.Article
	// 查询所有文章，关联用户和分类
	result := util.Db.Preload("User").Preload("Category").Find(&articles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}
	// 渲染模板并传递数据
	c.HTML(http.StatusOK, "article_list.html", gin.H{
		"title":    "文章列表 - Go博客",
		"articles": articles,
	})
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	var article model.Article
	result := util.Db.Preload("User").Preload("Category").First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "文章不存在",
		})
		return
	}

	// 渲染模板
	c.HTML(http.StatusOK, "article_detail.html", gin.H{
		"article": article,
	})
}
