package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex 首页
func GetIndex(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	// 跳转到首页模板
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Go博客首页",
		"user":  user,
	})
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	var articles []model.Article
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user != nil {
		// 检查用户是否是管理员
		userMap, ok := user.(map[string]interface{})
		if ok && userMap["Username"] == "admin" {
			// 管理员可以看到所有文章
			result := util.Db.Preload("User").Preload("Category").Find(&articles)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
				return
			}
		} else {
			// 普通用户只能看到公开的文章
			result := util.Db.Preload("User").Preload("Category").Where("visibility = ?", 1).Find(&articles)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
				return
			}
		}
	} else {
		// 未登录用户只能看到公开的文章
		result := util.Db.Preload("User").Preload("Category").Where("visibility = ?", 1).Find(&articles)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
			return
		}
	}
	// 渲染模板并传递数据
	c.HTML(http.StatusOK, "article_list.html", gin.H{
		"title":    "文章列表 - Go博客",
		"articles": articles,
		"user":     user,
	})
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var article model.Article
	result := util.Db.Preload("User").Preload("Category").First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "文章不存在",
			"user":  user,
		})
		return
	}

	// 渲染模板
	c.HTML(http.StatusOK, "article_detail.html", gin.H{
		"article": article,
		"user":    user,
	})
}
