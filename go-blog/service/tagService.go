package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTagList(c *gin.Context) {
	var tags []model.Tag
	util.Db.Order("sort_order asc, use_count desc").Find(&tags)
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func GetHotTags(c *gin.Context) {
	var tags []model.Tag
	util.Db.Order("use_count desc").Limit(10).Find(&tags)
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
