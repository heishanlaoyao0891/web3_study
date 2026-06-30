package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLearningPaths(c *gin.Context) {
	user := util.GetUserFromContext(c)

	var paths []model.LearningPath
	util.Db.Preload("Chapters").Order("sort_order asc, id desc").Find(&paths)

	c.HTML(http.StatusOK, "learning_paths.html", gin.H{
		"title": "学习路径",
		"paths": paths,
		"user":  user,
	})
}

func GetLearningPathDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var path model.LearningPath
	result := util.Db.Preload("Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc")
	}).First(&path, id)

	if result.Error != nil {
		c.Redirect(http.StatusFound, "/learning-paths")
		return
	}

	c.HTML(http.StatusOK, "learning_path_detail.html", gin.H{
		"title": path.Title,
		"path":  path,
		"user":  user,
	})
}

func GetCodeSnippets(c *gin.Context) {
	user := util.GetUserFromContext(c)
	category := c.Query("category")

	var snippets []model.CodeSnippet
	query := util.Db.Preload("User")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("id desc").Find(&snippets)

	var categories []string
	util.Db.Model(&model.CodeSnippet{}).Distinct("category").Pluck("category", &categories)

	c.HTML(http.StatusOK, "code_snippets.html", gin.H{
		"title":      "代码片段",
		"snippets":    snippets,
		"categories": categories,
		"current":     category,
		"user":       user,
	})
}

func GetContractTemplates(c *gin.Context) {
	user := util.GetUserFromContext(c)
	category := c.Query("category")

	var templates []model.ContractTemplate
	query := util.Db.Preload("User")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("id desc").Find(&templates)

	var categories []string
	util.Db.Model(&model.ContractTemplate{}).Distinct("category").Pluck("category", &categories)

	c.HTML(http.StatusOK, "contract_templates.html", gin.H{
		"title":      "代码模板",
		"templates":   templates,
		"categories": categories,
		"current":     category,
		"user":       user,
	})
}

func GetResources(c *gin.Context) {
	user := util.GetUserFromContext(c)
	category := c.Query("category")

	var resources []model.Resource
	query := util.Db

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("sort_order asc, id asc").Find(&resources)

	// 分类从 DB 动态获取，替代硬编码
	var categories []string
	util.Db.Model(&model.Resource{}).Distinct("category").Pluck("category", &categories)
	if len(categories) == 0 {
		categories = []string{"文档", "工具", "教程", "社区"}
	}

	c.HTML(http.StatusOK, "resources.html", gin.H{
		"title":      "资源导航",
		"resources":   resources,
		"categories": categories,
		"current":     category,
		"user":       user,
	})
}

// GetInterviewQuestions 面试题库列表（M3.3：增加按领域筛选）
func GetInterviewQuestions(c *gin.Context) {
	user := util.GetUserFromContext(c)
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	domainID := c.Query("domain") // M3.3 新增：技术领域筛选

	var questions []model.InterviewQuestion
	query := util.Db.Model(&model.InterviewQuestion{}).Preload("Categories")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	// M3.3：按技术领域筛选——通过 interview_question_categories 关联到顶层 Category
	if domainID != "" {
		if dID, err := strconv.Atoi(domainID); err == nil && dID > 0 {
			// 找到该领域及其所有子分类的 ID
			var categoryIDs []uint
			categoryIDs = append(categoryIDs, uint(dID))
			var subCategories []model.Category
			util.Db.Where("parent_id = ?", dID).Find(&subCategories)
			for _, sub := range subCategories {
				categoryIDs = append(categoryIDs, sub.ID)
			}
			// 通过中间表筛选
			query = query.Where(
				"id IN (SELECT interview_question_id FROM interview_question_categories WHERE category_id IN ?)",
				categoryIDs,
			)
		}
	}

	result := query.Order("id desc").Find(&questions)
	_ = result

	// 获取分类列表（从 DB 动态获取）
	var categories []string
	util.Db.Model(&model.InterviewQuestion{}).Distinct("category").Pluck("category", &categories)

	// 获取技术领域列表供筛选
	var domains []model.Category
	util.Db.Where("parent_id IS NULL").Order("sort_order desc, id asc").Find(&domains)

	c.HTML(http.StatusOK, "interview_questions.html", gin.H{
		"title":       "面试题库",
		"questions":    questions,
		"categories":  categories,
		"currentCat":  category,
		"currentDiff": difficulty,
		"currentDomain": domainID,
		"domains":      domains,
		"user":         user,
	})
}

// GetInterviewQuestionDetail 面试题详情
func GetInterviewQuestionDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)
	id, _ := strconv.Atoi(c.Param("id"))
	domainID := c.Query("domain") // 保留领域上下文用于返回

	var question model.InterviewQuestion
	result := util.Db.Preload("Categories").First(&question, id)

	if result.Error != nil {
		c.Redirect(http.StatusFound, "/interview")
		return
	}

	util.Db.Model(&question).Update("view_count", question.ViewCount+1)

	c.HTML(http.StatusOK, "interview_question_detail.html", gin.H{
		"title":         question.Title,
		"question":      question,
		"user":          user,
		"currentDomain": domainID,
	})
}