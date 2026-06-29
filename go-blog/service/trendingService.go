package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// getActiveTrendingTopics 获取当前有效的风口话题
func getActiveTrendingTopics(limit int) []model.TrendingTopic {
	var topics []model.TrendingTopic
	now := time.Now()
	util.Db.Where("status = ? AND (expire_at IS NULL OR expire_at > ?)", 1, now).
		Order("heat_score desc, created_at desc").
		Limit(limit).
		Find(&topics)
	return topics
}

// GetTrendingTopicsAPI API：返回风口话题列表（供前端 AJAX 调用）
func GetTrendingTopicsAPI(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	topics := getActiveTrendingTopics(limit)
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

// GetTrendingAdminPage 风口话题管理页面（仅管理员）
func GetTrendingAdminPage(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	var topics []model.TrendingTopic
	util.Db.Order("status desc, heat_score desc, id desc").Find(&topics)

	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

// PostTrendingCreate 创建风口话题（仅管理员）
func PostTrendingCreate(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	title := c.PostForm("title")
	summary := c.PostForm("summary")
	url := c.PostForm("url")
	domainIDStr := c.PostForm("domain_id")
	heatScoreStr := c.PostForm("heat_score")
	expireStr := c.PostForm("expire_hours")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	domainID, _ := strconv.Atoi(domainIDStr)
	heatScore, _ := strconv.Atoi(heatScoreStr)

	topic := model.TrendingTopic{
		Title:     title,
		Summary:   summary,
		URL:       url,
		DomainID:  uint(domainID),
		HeatScore: heatScore,
		Status:    1,
	}

	// 处理过期时间
	if expireStr != "" {
		if hours, err := strconv.Atoi(expireStr); err == nil && hours > 0 {
			expire := time.Now().Add(time.Duration(hours) * time.Hour)
			topic.ExpireAt = &expire
		}
	}

	if err := util.Db.Create(&topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "topic": topic})
}

// PostTrendingDelete 下线/删除风口话题（仅管理员）
func PostTrendingDelete(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	id := c.Param("id")
	var topic model.TrendingTopic
	if err := util.Db.First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "话题不存在"})
		return
	}

	util.Db.Delete(&topic)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PostTrendingOffline 下线风口话题（仅管理员，改状态为0而非删除）
func PostTrendingOffline(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	id := c.Param("id")
	var topic model.TrendingTopic
	if err := util.Db.First(&topic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "话题不存在"})
		return
	}

	topic.Status = 0
	util.Db.Save(&topic)
	c.JSON(http.StatusOK, gin.H{"success": true})
}