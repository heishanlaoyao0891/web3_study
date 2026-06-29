package service

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArticleCreate 显示发布文章页面
func GetArticleCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)

	var categories []model.Category
	util.Db.Order("parent_id asc, sort_order desc, id asc").Find(&categories)

	c.HTML(http.StatusOK, "article_create.html", gin.H{
		"categories": categories,
		"user":       user,
	})
}

// PostArticleCreate 处理发布文章请求
func PostArticleCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)

	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID := c.PostForm("category_id")
	visibility := c.PostForm("visibility")

	var categoryIDUint uint
	fmt.Sscanf(categoryID, "%d", &categoryIDUint)

	visibilityInt, err := strconv.Atoi(visibility)
	if err != nil {
		visibilityInt = 1 // 默认公开
	}

	var userID uint
	if user != nil {
		userID = user.ID
	}

	article := model.Article{
		Title:      title,
		Content:    content,
		Status:     1,
		Visibility: visibilityInt,
		UserID:     userID,
		CategoryID: categoryIDUint,
	}

	util.Db.Create(&article)

	c.Redirect(http.StatusFound, "/article/list")
}