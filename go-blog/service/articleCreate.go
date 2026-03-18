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
	// 从上下文中获取用户信息（中间件已确保用户已登录）
	user := util.GetUserFromContext(c)

	var categories []model.Category
	util.Db.Find(&categories)
	c.HTML(http.StatusOK, "article_create.html", gin.H{
		"categories": categories,
		"user":       user,
	})
}

// PostArticleCreate 处理发布文章请求
func PostArticleCreate(c *gin.Context) {
	// 从上下文中获取用户信息（中间件已确保用户已登录）
	user := util.GetUserFromContext(c)

	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID := c.PostForm("category_id")
	visibility := c.PostForm("visibility")

	// 将categoryID转换为uint类型
	var categoryIDUint uint
	fmt.Sscanf(categoryID, "%d", &categoryIDUint)

	// 将visibility转换为int类型
	println("表单提交的visibility值:", visibility)
	visibilityInt, err := strconv.Atoi(visibility)
	if err != nil {
		println("转换错误:", err.Error())
		visibilityInt = 1 // 默认公开
	}
	println("转换后的visibility值:", visibilityInt)

	// 获取当前用户ID
	userMap, ok := user.(map[string]interface{})
	var userID uint
	if ok {
		// 尝试获取用户ID，处理不同类型
		idValue := userMap["ID"]
		switch v := idValue.(type) {
		case uint:
			userID = v
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		default:
			userID = 1 // 默认使用admin用户
		}
	} else {
		userID = 1 // 默认使用admin用户
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

	// 重定向到文章列表页面
	c.Redirect(http.StatusFound, "/article/list")
}
