package service

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

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
		// 检查用户是否是管理员  (user.(map[string]interface{})是类型断言，类似java的instanceof)
		userMap, ok := user.(map[string]interface{})
		if ok && userMap["Username"] == "admin" {
			// 管理员可以看到所有文章
			result := util.Db.Preload("User").Preload("Category").Find(&articles)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
				return
			}
		} else {
			// 普通用户可以看到公开的文章和自己的私有文章
			// 尝试获取用户ID，处理不同类型
			var userID uint
			var ok bool

			// 检查ID的类型
			idValue := userMap["ID"]
			switch v := idValue.(type) {
			case uint:
				userID = v
				ok = true
			case float64:
				userID = uint(v)
				ok = true
			case int:
				userID = uint(v)
				ok = true
			default:
				ok = false
			}

			if ok {
				// 查询公开文章或自己的文章
				result := util.Db.Preload("User").Preload("Category").Where("visibility = ? OR user_id = ?", 1, userID).Find(&articles)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
					return
				}
			} else {
				// 未获取到用户ID，只显示公开文章
				result := util.Db.Preload("User").Preload("Category").Where("visibility = ?", 1).Find(&articles)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
					return
				}
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

	// 检查文章访问权限
	if article.Visibility == 0 {
		// 私有文章，只有作者和管理员可以访问
		if user == nil {
			c.HTML(http.StatusForbidden, "index.html", gin.H{
				"title": "错误 - Go博客",
				"error": "无权限访问此文章",
				"user":  user,
			})
			return
		}

		userMap, ok := user.(map[string]interface{})
		if !ok {
			c.HTML(http.StatusForbidden, "index.html", gin.H{
				"title": "错误 - Go博客",
				"error": "无权限访问此文章",
				"user":  user,
			})
			return
		}

		// 检查是否是管理员
		if userMap["Username"] == "admin" {
			// 管理员可以访问所有文章
		} else {
			// 检查是否是文章作者
			// 尝试获取用户ID，处理不同类型
			var userID uint
			var ok bool

			// 检查ID的类型
			idValue := userMap["ID"]
			switch v := idValue.(type) {
			case uint:
				userID = v
				ok = true
			case float64:
				userID = uint(v)
				ok = true
			case int:
				userID = uint(v)
				ok = true
			default:
				ok = false
			}

			if !ok || userID != article.UserID {
				c.HTML(http.StatusForbidden, "index.html", gin.H{
					"title": "错误 - Go博客",
					"error": "无权限访问此文章",
					"user":  user,
				})
				return
			}
		}
	}

	// 渲染模板
	c.HTML(http.StatusOK, "article_detail.html", gin.H{
		"article": article,
		"user":    user,
	})
}

// GetArticleEdit 显示编辑文章页面
func GetArticleEdit(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

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

	// 检查编辑权限
	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	// 检查是否是管理员或文章作者
	isAdmin := userMap["Username"] == "admin"
	var userID uint
	var idOk bool

	// 检查ID的类型
	idValue := userMap["ID"]
	switch v := idValue.(type) {
	case uint:
		userID = v
		idOk = true
	case float64:
		userID = uint(v)
		idOk = true
	case int:
		userID = uint(v)
		idOk = true
	default:
		idOk = false
	}

	if !isAdmin && (!idOk || userID != article.UserID) {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	// 获取所有分类
	var categories []model.Category
	util.Db.Find(&categories)

	// 渲染编辑页面
	c.HTML(http.StatusOK, "article_edit.html", gin.H{
		"article":    article,
		"categories": categories,
		"user":       user,
	})
}

// PostArticleEdit 处理编辑文章请求
func PostArticleEdit(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id := c.Param("id")
	var article model.Article
	result := util.Db.First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "文章不存在",
			"user":  user,
		})
		return
	}

	// 检查编辑权限
	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	// 检查是否是管理员或文章作者
	isAdmin := userMap["Username"] == "admin"
	var userID uint
	var idOk bool

	// 检查ID的类型
	idValue := userMap["ID"]
	switch v := idValue.(type) {
	case uint:
		userID = v
		idOk = true
	case float64:
		userID = uint(v)
		idOk = true
	case int:
		userID = uint(v)
		idOk = true
	default:
		idOk = false
	}

	if !isAdmin && (!idOk || userID != article.UserID) {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	// 获取表单数据
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID := c.PostForm("category_id")
	visibility := c.PostForm("visibility")

	// 将categoryID转换为uint类型
	var categoryIDUint uint
	fmt.Sscanf(categoryID, "%d", &categoryIDUint)

	// 将visibility转换为int类型
	var visibilityInt int
	fmt.Sscanf(visibility, "%d", &visibilityInt)

	// 更新文章
	article.Title = title
	article.Content = content
	article.CategoryID = categoryIDUint
	article.Visibility = visibilityInt

	// 保存到数据库
	util.Db.Save(&article)

	// 重定向到文章详情页面
	c.Redirect(http.StatusFound, "/article/detail/"+id)
}
