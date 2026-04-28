package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetQuestionList(c *gin.Context) {
	user := util.GetUserFromContext(c)
	pageStr := c.Query("page")
	pageSize := 10

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	status := c.Query("status")
	var total int64
	var questions []model.Question

	query := util.Db.Model(&model.Question{}).Preload("User")
	if status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", s)
		}
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&questions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取问答列表失败"})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "question_list.html", gin.H{
		"title":      "问答社区 - Go博客",
		"questions":  questions,
		"user":       user,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"status":     status,
	})
}

func GetQuestionDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)
	id := c.Param("id")

	var question model.Question
	result := util.Db.Preload("User").Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("is_best desc, created_at asc")
	}).Preload("Answers.User").First(&question, id)

	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"error": "问题不存在",
			"user":  user,
		})
		return
	}

	util.Db.Model(&question).Update("view_count", question.ViewCount+1)

	c.HTML(http.StatusOK, "question_detail.html", gin.H{
		"question": question,
		"user":     user,
	})
}

func GetQuestionCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "question_create.html", gin.H{
		"user": user,
	})
}

func PostQuestionCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息错误"})
		return
	}

	var userID uint
	switch v := userMap["ID"].(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	}

	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能为空"})
		return
	}

	question := model.Question{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	util.Db.Create(&question)
	c.Redirect(http.StatusFound, "/qa/detail/"+strconv.FormatUint(uint64(question.ID), 10))
}

func PostAnswerCreate(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息错误"})
		return
	}

	var userID uint
	switch v := userMap["ID"].(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	}

	questionID := c.PostForm("question_id")
	content := c.PostForm("content")

	var qID uint
	for _, ch := range questionID {
		if ch >= '0' && ch <= '9' {
			qID = qID*10 + uint(ch-'0')
		}
	}

	answer := model.Answer{
		QuestionID: qID,
		UserID:     userID,
		Content:    content,
	}

	util.Db.Create(&answer)
	util.Db.Model(&model.Question{}).Where("id = ?", qID).Update("answer_count", gorm.Expr("answer_count + 1"))

	c.Redirect(http.StatusFound, "/qa/detail/"+questionID)
}

func PostAcceptAnswer(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息错误"})
		return
	}

	var userID uint
	switch v := userMap["ID"].(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	}

	answerID := c.Param("id")
	var answer model.Answer
	result := util.Db.First(&answer, answerID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "回答不存在"})
		return
	}

	var question model.Question
	util.Db.First(&question, answer.QuestionID)

	if question.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有提问者可以采纳答案"})
		return
	}

	util.Db.Model(&answer).Update("is_best", 1)
	util.Db.Model(&question).Updates(map[string]interface{}{
		"status":        1,
		"best_answer_id": answer.ID,
	})

	c.Redirect(http.StatusFound, "/qa/detail/"+strconv.FormatUint(uint64(question.ID), 10))
}
