package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCourseList(c *gin.Context) {
	user := util.GetUserFromContext(c)
	category := c.Query("category")
	pageStr := c.Query("page")
	pageSize := 12

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	var total int64
	var courses []model.Course

	query := util.Db.Model(&model.Course{}).Preload("Author")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Order("priority desc, created_at desc").Offset(offset).Limit(pageSize).Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教程列表失败"})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	var categories []string
	util.Db.Model(&model.Course{}).Distinct().Pluck("category", &categories)

	c.HTML(http.StatusOK, "course_list.html", gin.H{
		"title":      "教程中心 - Go博客",
		"courses":    courses,
		"categories": categories,
		"user":       user,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"category":   category,
	})
}

func GetCourseDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)
	id := c.Param("id")

	var course model.Course
	result := util.Db.Preload("Author").Preload("Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc")
	}).First(&course, id)

	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"error": "教程不存在",
			"user":  user,
		})
		return
	}

	util.Db.Model(&course).Update("view_count", course.ViewCount+1)

	var learningRecord model.LearningRecord
	if user != nil {
		userMap, ok := user.(map[string]interface{})
		if ok {
			var userID uint
			switch v := userMap["ID"].(type) {
			case uint:
				userID = v
			case float64:
				userID = uint(v)
			case int:
				userID = uint(v)
			}
			util.Db.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&learningRecord)
		}
	}

	c.HTML(http.StatusOK, "course_detail.html", gin.H{
		"course":         course,
		"user":           user,
		"learningRecord": learningRecord,
	})
}

func GetChapterDetail(c *gin.Context) {
	user := util.GetUserFromContext(c)
	courseID := c.Param("course_id")
	chapterID := c.Param("chapter_id")

	var chapter model.CourseChapter
	result := util.Db.Preload("Course").First(&chapter, chapterID)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"error": "章节不存在",
			"user":  user,
		})
		return
	}

	util.Db.Model(&chapter).Update("view_count", chapter.ViewCount+1)

	if user != nil {
		userMap, ok := user.(map[string]interface{})
		if ok {
			var userID uint
			switch v := userMap["ID"].(type) {
			case uint:
				userID = v
			case float64:
				userID = uint(v)
			case int:
				userID = uint(v)
			}

			var courseIDUint uint
			for _, ch := range courseID {
				if ch >= '0' && ch <= '9' {
					courseIDUint = courseIDUint*10 + uint(ch-'0')
				}
			}

			now := time.Now()
			util.Db.Where("user_id = ? AND course_id = ?", userID, courseIDUint).
				Assign(model.LearningRecord{
					ChapterID:   &chapter.ID,
					Progress:    0,
					LastLearnAt: &now,
				}).
				FirstOrCreate(&model.LearningRecord{
					UserID:      userID,
					CourseID:    courseIDUint,
					ChapterID:   &chapter.ID,
					LastLearnAt: &now,
				})
		}
	}

	c.HTML(http.StatusOK, "chapter_detail.html", gin.H{
		"chapter": chapter,
		"user":    user,
	})
}
