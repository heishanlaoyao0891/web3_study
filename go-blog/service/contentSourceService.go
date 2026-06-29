package service

import (
	"go-blog/crawler"
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetContentSources 抓取源管理页面
func GetContentSources(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	var sources []model.ContentSource
	util.Db.Order("enabled desc, created_at desc").Find(&sources)

	// 获取领域名称映射
	var domains []model.Category
	util.Db.Where("parent_id IS NULL").Find(&domains)
	domainNames := make(map[uint]string)
	for _, d := range domains {
		domainNames[d.ID] = d.Name
	}

	// 统计各源最近一次运行状态
	type sourceStatus struct {
		Status string
		Count  int
	}
	lastStatus := make(map[uint]sourceStatus)
	for _, s := range sources {
		var log model.CrawlLog
		if err := util.Db.Where("source_id = ?", s.ID).Order("created_at desc").First(&log).Error; err == nil {
			lastStatus[s.ID] = sourceStatus{Status: log.Status, Count: log.SavedCount}
		}
	}

	c.HTML(http.StatusOK, "content_sources.html", gin.H{
		"sources":     sources,
		"domainNames": domainNames,
		"lastStatus":  lastStatus,
		"user":        user,
	})
}

// PostContentSourceToggle 启用/禁用抓取源
func PostContentSourceToggle(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效 ID"})
		return
	}

	var source model.ContentSource
	if err := util.Db.First(&source, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "源不存在"})
		return
	}

	source.Enabled = !source.Enabled
	util.Db.Save(&source)

	// 动态更新调度器
	if source.Enabled {
		crawler.GetScheduler().AddSource(source)
	} else {
		crawler.GetScheduler().RemoveSource(source.ID)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "enabled": source.Enabled})
}

// PostContentSourceDelete 删除抓取源
func PostContentSourceDelete(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效 ID"})
		return
	}

	crawler.GetScheduler().RemoveSource(uint(id))
	util.Db.Delete(&model.ContentSource{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PostContentSourceRunNow 立即执行一次抓取
func PostContentSourceRunNow(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效 ID"})
		return
	}

	var source model.ContentSource
	if err := util.Db.First(&source, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "源不存在"})
		return
	}

	log, err := crawler.RunOnce(source)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "log": log})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "log": log})
}
