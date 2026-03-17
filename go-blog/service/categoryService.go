package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategoryList 显示类别列表页面（仅管理员可见）
func GetCategoryList(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return
	}

	// 获取所有类别
	var categories []model.Category
	util.Db.Find(&categories)

	c.HTML(http.StatusOK, "category_list.html", gin.H{
		"categories": categories,
		"user":       user,
	})
}

// GetCategoryCreate 显示添加类别页面（仅管理员可见）
func GetCategoryCreate(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return
	}

	// 渲染添加类别页面
	c.HTML(http.StatusOK, "category_edit.html", gin.H{
		"category": model.Category{},
		"user":     user,
	})
}

// PostCategoryCreate 处理添加类别请求（仅管理员可见）
func PostCategoryCreate(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return
	}

	// 获取表单数据
	name := c.PostForm("name")
	desc := c.PostForm("desc")

	// 检查类别名称是否已存在
	var existingCategory model.Category
	result := util.Db.Where("name = ?", name).First(&existingCategory)
	if result.Error == nil {
		c.HTML(http.StatusOK, "category_edit.html", gin.H{
			"category": model.Category{Name: name, Desc: desc},
			"error":    "类别名称已存在",
			"user":     user,
		})
		return
	}

	// 创建新类别
	newCategory := model.Category{
		Name: name,
		Desc: desc,
	}
	util.Db.Create(&newCategory)

	// 重定向到类别列表页面
	c.Redirect(http.StatusFound, "/category/list")
}

// GetCategoryEdit 显示编辑类别页面（仅管理员可见）
func GetCategoryEdit(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return
	}

	// 获取类别ID
	id := c.Param("id")

	// 查找类别
	var category model.Category
	result := util.Db.First(&category, id)
	if result.Error != nil {
		c.HTML(http.StatusOK, "category_list.html", gin.H{
			"error":      "类别不存在",
			"categories": []model.Category{},
			"user":       user,
		})
		return
	}

	// 渲染编辑类别页面
	c.HTML(http.StatusOK, "category_edit.html", gin.H{
		"category": category,
		"user":     user,
	})
}

// PostCategoryEdit 处理编辑类别请求（仅管理员可见）
func PostCategoryEdit(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return
	}

	// 获取类别ID
	id := c.Param("id")

	// 查找类别
	var category model.Category
	result := util.Db.First(&category, id)
	if result.Error != nil {
		c.HTML(http.StatusOK, "category_list.html", gin.H{
			"error":      "类别不存在",
			"categories": []model.Category{},
			"user":       user,
		})
		return
	}

	// 获取表单数据
	name := c.PostForm("name")
	desc := c.PostForm("desc")

	// 检查类别名称是否已存在（排除当前类别）
	var existingCategory model.Category
	result = util.Db.Where("name = ? AND id != ?", name, id).First(&existingCategory)
	if result.Error == nil {
		c.HTML(http.StatusOK, "category_edit.html", gin.H{
			"category": model.Category{ID: category.ID, Name: name, Desc: desc},
			"error":    "类别名称已存在",
			"user":     user,
		})
		return
	}

	// 更新类别
	category.Name = name
	category.Desc = desc
	util.Db.Save(&category)

	// 重定向到类别列表页面
	c.Redirect(http.StatusFound, "/category/list")
}
