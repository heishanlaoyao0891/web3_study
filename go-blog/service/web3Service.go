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
	var paths []model.LearningPath
	util.Db.Preload("Chapters").Order("sort_order asc, id desc").Find(&paths)

	c.HTML(http.StatusOK, "learning_paths.html", gin.H{
		"title": "学习路径",
		"paths": paths,
	})
}

func GetLearningPathDetail(c *gin.Context) {
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
	})
}

func GetCodeSnippets(c *gin.Context) {
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
		"snippets":   snippets,
		"categories": categories,
		"current":    category,
	})
}

func GetContractTemplates(c *gin.Context) {
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
		"title":      "合约模板",
		"templates":  templates,
		"categories": categories,
		"current":    category,
	})
}

func GetResources(c *gin.Context) {
	category := c.Query("category")

	var resources []model.Resource
	query := util.Db

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Order("sort_order asc, id asc").Find(&resources)

	categories := []string{"文档", "工具", "教程", "社区"}

	c.HTML(http.StatusOK, "resources.html", gin.H{
		"title":      "资源导航",
		"resources":  resources,
		"categories": categories,
		"current":    category,
	})
}

func GetInterviewQuestions(c *gin.Context) {
	category := c.Query("category")
	difficulty := c.Query("difficulty")

	var questions []model.InterviewQuestion
	query := util.Db.Preload("User")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	query.Order("id desc").Find(&questions)

	var categories []string
	util.Db.Model(&model.InterviewQuestion{}).Distinct("category").Pluck("category", &categories)

	c.HTML(http.StatusOK, "interview_questions.html", gin.H{
		"title":       "面试题库",
		"questions":   questions,
		"categories":  categories,
		"currentCat":  category,
		"currentDiff": difficulty,
	})
}

func GetInterviewQuestionDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var question model.InterviewQuestion
	result := util.Db.Preload("User").First(&question, id)

	if result.Error != nil {
		c.Redirect(http.StatusFound, "/interview")
		return
	}

	util.Db.Model(&question).Update("view_count", question.ViewCount+1)

	c.HTML(http.StatusOK, "interview_question_detail.html", gin.H{
		"title":    question.Title,
		"question": question,
	})
}
