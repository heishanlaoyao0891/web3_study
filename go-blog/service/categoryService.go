package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// requireAdmin 要求管理员权限，非管理员渲染错误页。返回 false 表示已响应。
func requireAdmin(c *gin.Context) bool {
	if !util.IsAdmin(c) {
		user := util.GetUserFromContext(c)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
			"user":  user,
		})
		return false
	}
	return true
}

// GetCategoryList 显示类别列表页面（仅管理员可见）
func GetCategoryList(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	// 获取所有类别，按层级排序
	var categories []model.Category
	util.Db.Order("parent_id asc, sort_order desc, id asc").Find(&categories)

	c.HTML(http.StatusOK, "category_list.html", gin.H{
		"categories": categories,
		"user":       user,
	})
}

// GetCategoryCreate 显示添加类别页面（仅管理员可见）
func GetCategoryCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	// 加载顶层领域供选择父分类
	var domains []model.Category
	util.Db.Where("parent_id IS NULL").Order("sort_order desc, id asc").Find(&domains)

	c.HTML(http.StatusOK, "category_edit.html", gin.H{
		"category": model.Category{},
		"domains":  domains,
		"user":     user,
	})
}

// PostCategoryCreate 处理添加类别请求（仅管理员可见）
func PostCategoryCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	name := c.PostForm("name")
	desc := c.PostForm("desc")
	parentIDStr := c.PostForm("parent_id")
	sortOrderStr := c.PostForm("sort_order")
	icon := c.PostForm("icon")

	// 检查类别名称是否已存在
	var existingCategory model.Category
	result := util.Db.Where("name = ?", name).First(&existingCategory)
	if result.Error == nil {
		var domains []model.Category
		util.Db.Where("parent_id IS NULL").Order("sort_order desc, id asc").Find(&domains)
		c.HTML(http.StatusOK, "category_edit.html", gin.H{
			"category": model.Category{Name: name, Desc: desc, Icon: icon},
			"domains":  domains,
			"error":    "类别名称已存在",
			"user":     user,
		})
		return
	}

	newCategory := model.Category{
		Name:      name,
		Desc:      desc,
		Icon:      icon,
		SortOrder: 0,
	}

	// 处理父分类
	if parentIDStr != "" {
		if pid, err := strconv.Atoi(parentIDStr); err == nil && pid > 0 {
			newCategory.ParentID = util.UintPtr(uint(pid))
		}
	}

	// 处理排序
	if so, err := strconv.Atoi(sortOrderStr); err == nil {
		newCategory.SortOrder = so
	}

	util.Db.Create(&newCategory)

	c.Redirect(http.StatusFound, "/category/list")
}

// GetCategoryEdit 显示编辑类别页面（仅管理员可见）
func GetCategoryEdit(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	id := c.Param("id")

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

	var domains []model.Category
	util.Db.Where("parent_id IS NULL AND id != ?", id).Order("sort_order desc, id asc").Find(&domains)

	c.HTML(http.StatusOK, "category_edit.html", gin.H{
		"category": category,
		"domains":  domains,
		"user":     user,
	})
}

// PostCategoryEdit 处理编辑类别请求（仅管理员可见）
func PostCategoryEdit(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	id := c.Param("id")

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

	name := c.PostForm("name")
	desc := c.PostForm("desc")
	parentIDStr := c.PostForm("parent_id")
	sortOrderStr := c.PostForm("sort_order")
	icon := c.PostForm("icon")

	// 检查类别名称是否已存在（排除当前类别）
	var existingCategory model.Category
	result = util.Db.Where("name = ? AND id != ?", name, id).First(&existingCategory)
	if result.Error == nil {
		var domains []model.Category
		util.Db.Where("parent_id IS NULL AND id != ?", id).Order("sort_order desc, id asc").Find(&domains)
		c.HTML(http.StatusOK, "category_edit.html", gin.H{
			"category": model.Category{ID: category.ID, Name: name, Desc: desc, Icon: icon},
			"domains":  domains,
			"error":    "类别名称已存在",
			"user":     user,
		})
		return
	}

	category.Name = name
	category.Desc = desc
	category.Icon = icon

	// 处理父分类
	if parentIDStr != "" {
		if pid, err := strconv.Atoi(parentIDStr); err == nil && pid > 0 {
			category.ParentID = util.UintPtr(uint(pid))
		}
	} else {
		category.ParentID = nil
	}

	if so, err := strconv.Atoi(sortOrderStr); err == nil {
		category.SortOrder = so
	}

	util.Db.Save(&category)

	c.Redirect(http.StatusFound, "/category/list")
}

// PostCategoryDelete 删除分类
func PostCategoryDelete(c *gin.Context) {
	if !util.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	id := c.Param("id")

	var category model.Category
	result := util.Db.First(&category, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	// 检查是否有子分类
	var childCount int64
	util.Db.Model(&model.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下有子分类，无法删除"})
		return
	}

	var articleCount int64
	util.Db.Model(&model.Article{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下还有文章，无法删除"})
		return
	}

	util.Db.Delete(&category)

	c.Redirect(http.StatusFound, "/category/list")
}