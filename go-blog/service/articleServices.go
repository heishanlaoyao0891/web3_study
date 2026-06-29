package service

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetIndex 首页
func GetIndex(c *gin.Context) {
	user := util.GetUserFromContext(c)

	var articleCount int64
	var categoryCount int64
	var userCount int64

	util.Db.Model(&model.Article{}).Count(&articleCount)
	util.Db.Model(&model.Category{}).Count(&categoryCount)
	util.Db.Model(&model.User{}).Count(&userCount)

	var latestArticles []model.Article
	util.Db.Preload("User").Preload("Category").Preload("Categories").Order("created_at desc").Limit(5).Find(&latestArticles)

	// 热门标签（M1.4：从 DB 取，替代硬编码）
	var hotTags []model.Tag
	util.Db.Order("use_count desc, id asc").Limit(10).Find(&hotTags)

	// 技术领域（顶层分类，供导航展示）
	var domains []model.Category
	util.Db.Where("parent_id IS NULL").Order("sort_order desc, id asc").Find(&domains)

	// 加载站点配置（M1.3 首页文案去硬编码）
	siteConfig := loadSiteConfig()

	// 风⼝话题（M3.1 有效话题）
	var trendingTopics []model.TrendingTopic
	util.Db.Where("status = ? AND (expire_at IS NULL OR expire_at > ?)", 1, time.Now()).
		Order("heat_score desc, created_at desc").
		Limit(5).
		Find(&trendingTopics)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": siteConfig["site_title"],
		"user":  user,
		"stats": map[string]interface{}{
			"articleCount":  articleCount,
			"categoryCount": categoryCount,
			"userCount":     userCount,
		},
		"latestArticles":  latestArticles,
		"hotTags":         hotTags,
		"domains":         domains,
		"site":            siteConfig,
		"trendingTopics":  trendingTopics,
	})
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	user := util.GetUserFromContext(c)
	keyword := c.Query("keyword")
	pageStr := c.Query("page")
	domainID := c.Query("domain") // 新增：按领域筛选
	pageSize := 10

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var total int64
	var articles []model.Article

	query := util.Db.Model(&model.Article{}).Preload("User").Preload("Category").Preload("Categories")

	isAdmin := util.IsAdmin(c)
	var userID uint
	if user != nil {
		userID = user.ID
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按领域筛选：取该领域下所有子分类的文章
	if domainID != "" {
		if dID, err := strconv.Atoi(domainID); err == nil && dID > 0 {
			var subCategoryIDs []uint
			util.Db.Model(&model.Category{}).Where("parent_id = ?", dID).Pluck("id", &subCategoryIDs)
			subCategoryIDs = append(subCategoryIDs, uint(dID))
			query = query.Where("category_id IN ?", subCategoryIDs)
		}
	}

	if !isAdmin {
		if userID > 0 {
			query = query.Where("visibility = ? OR user_id = ?", 1, userID)
		} else {
			query = query.Where("visibility = ?", 1)
		}
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&articles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 加载领域列表供筛选
	var domains []model.Category
	util.Db.Where("parent_id IS NULL").Order("sort_order desc, id asc").Find(&domains)

	c.HTML(http.StatusOK, "article_list.html", gin.H{
		"title":      "文章列表 - Go博客",
		"articles":   articles,
		"user":       user,
		"keyword":    keyword,
		"domain":     domainID,
		"domains":    domains,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
	})
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var article model.Article
	result := util.Db.Preload("User").Preload("Category").Preload("Categories").First(&article, id)
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
		isAdmin := util.IsAdmin(c)
		isAuthor := user != nil && user.ID == article.UserID
		if !isAdmin && !isAuthor {
			c.HTML(http.StatusForbidden, "index.html", gin.H{
				"title": "错误 - Go博客",
				"error": "无权限访问此文章",
				"user":  user,
			})
			return
		}
	}

	// 获取文章评论
	var comments []model.Comment
	util.Db.Preload("User").Where("article_id = ? AND parent_id IS NULL", article.ID).Order("created_at desc").Find(&comments)

	c.HTML(http.StatusOK, "article_detail.html", gin.H{
		"article":  article,
		"user":     user,
		"comments": comments,
	})
}

// GetArticleEdit 显示编辑文章页面
func GetArticleEdit(c *gin.Context) {
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var article model.Article
	result := util.Db.Preload("User").Preload("Category").Preload("Categories").First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "文章不存在",
			"user":  user,
		})
		return
	}

	if !canEditArticle(user, &article) {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	var categories []model.Category
	util.Db.Order("parent_id asc, sort_order desc, id asc").Find(&categories)

	c.HTML(http.StatusOK, "article_edit.html", gin.H{
		"article":    article,
		"categories": categories,
		"user":       user,
	})
}

// PostArticleEdit 处理编辑文章请求
func PostArticleEdit(c *gin.Context) {
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var article model.Article
	result := util.Db.Preload("User").Preload("Category").Preload("Categories").First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "文章不存在",
			"user":  user,
		})
		return
	}

	if !canEditArticle(user, &article) {
		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"title": "错误 - Go博客",
			"error": "无权限编辑此文章",
			"user":  user,
		})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryID := c.PostForm("category_id")
	visibility := c.PostForm("visibility")

	var categoryIDUint uint
	fmt.Sscanf(categoryID, "%d", &categoryIDUint)

	visibilityInt, err := strconv.Atoi(visibility)
	if err != nil {
		visibilityInt = 1
	}

	article.Title = title
	article.Content = content
	article.CategoryID = categoryIDUint
	article.Visibility = visibilityInt

	util.Db.Save(&article)

	c.Redirect(http.StatusFound, "/article/detail/"+id)
}

// PostArticleDelete 删除文章
func PostArticleDelete(c *gin.Context) {
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var article model.Article
	result := util.Db.First(&article, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if !canEditArticle(user, &article) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此文章"})
		return
	}

	util.Db.Delete(&article)

	c.Redirect(http.StatusFound, "/article/list")
}

// canEditArticle 判断用户是否有权编辑/删除文章（管理员或作者）
func canEditArticle(user *model.User, article *model.Article) bool {
	if user == nil {
		return false
	}
	if user.Role == "admin" {
		return true
	}
	return user.ID == article.UserID
}

// loadSiteConfig 从数据库加载站点配置，返回 map 供模板使用
func loadSiteConfig() map[string]string {
	config := make(map[string]string)
	var configs []model.SiteConfig
	util.Db.Find(&configs)
	for _, c := range configs {
		config[c.Key] = c.Value
	}
	// 默认值
	if config["site_title"] == "" {
		config["site_title"] = "技术学习平台"
	}
	if config["site_subtitle"] == "" {
		config["site_subtitle"] = "Java | Go | Python | AI | Web3 | 全栈技术社区"
	}
	return config
}