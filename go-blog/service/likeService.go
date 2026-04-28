package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToggleLike(c *gin.Context) {
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
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	targetType := c.PostForm("target_type")
	targetID := c.PostForm("target_id")

	var tType int
	switch targetType {
	case "article":
		tType = 1
	case "comment":
		tType = 2
	case "question":
		tType = 3
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标类型"})
		return
	}

	var targetIDUint uint
	for _, ch := range targetID {
		if ch < '0' || ch > '9' {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
			return
		}
	}
	for i := 0; i < len(targetID); i++ {
		targetIDUint = targetIDUint*10 + uint(targetID[i]-'0')
	}

	var existingLike model.Like
	result := util.Db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, tType, targetIDUint).First(&existingLike)

	if result.Error == nil {
		util.Db.Delete(&existingLike)
		switch tType {
		case 1:
			util.Db.Model(&model.Article{}).Where("id = ?", targetIDUint).Update("like_count", gorm.Expr("like_count - 1"))
		case 3:
			util.Db.Model(&model.Question{}).Where("id = ?", targetIDUint).Update("like_count", gorm.Expr("like_count - 1"))
		}
		c.JSON(http.StatusOK, gin.H{"liked": false, "message": "取消点赞"})
	} else {
		like := model.Like{
			UserID:     userID,
			TargetType: tType,
			TargetID:   targetIDUint,
		}
		util.Db.Create(&like)
		switch tType {
		case 1:
			util.Db.Model(&model.Article{}).Where("id = ?", targetIDUint).Update("like_count", gorm.Expr("like_count + 1"))
		case 3:
			util.Db.Model(&model.Question{}).Where("id = ?", targetIDUint).Update("like_count", gorm.Expr("like_count + 1"))
		}
		c.JSON(http.StatusOK, gin.H{"liked": true, "message": "点赞成功"})
	}
}

func ToggleFavorite(c *gin.Context) {
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
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	targetType := c.PostForm("target_type")
	targetID := c.PostForm("target_id")

	var tType int
	switch targetType {
	case "article":
		tType = 1
	case "course":
		tType = 2
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标类型"})
		return
	}

	var targetIDUint uint
	for _, ch := range targetID {
		if ch < '0' || ch > '9' {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
			return
		}
	}
	for i := 0; i < len(targetID); i++ {
		targetIDUint = targetIDUint*10 + uint(targetID[i]-'0')
	}

	var existingFav model.Favorite
	result := util.Db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, tType, targetIDUint).First(&existingFav)

	if result.Error == nil {
		util.Db.Delete(&existingFav)
		if tType == 1 {
			util.Db.Model(&model.Article{}).Where("id = ?", targetIDUint).Update("favorite_count", gorm.Expr("favorite_count - 1"))
		}
		c.JSON(http.StatusOK, gin.H{"favorited": false, "message": "取消收藏"})
	} else {
		fav := model.Favorite{
			UserID:     userID,
			TargetType: tType,
			TargetID:   targetIDUint,
		}
		util.Db.Create(&fav)
		if tType == 1 {
			util.Db.Model(&model.Article{}).Where("id = ?", targetIDUint).Update("favorite_count", gorm.Expr("favorite_count + 1"))
		}
		c.JSON(http.StatusOK, gin.H{"favorited": true, "message": "收藏成功"})
	}
}

func GetUserFavorites(c *gin.Context) {
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

	var favorites []model.Favorite
	util.Db.Where("user_id = ?", userID).Find(&favorites)

	var articleIDs []uint
	for _, fav := range favorites {
		if fav.TargetType == 1 {
			articleIDs = append(articleIDs, fav.TargetID)
		}
	}

	var articles []model.Article
	if len(articleIDs) > 0 {
		util.Db.Preload("User").Preload("Category").Where("id IN ?", articleIDs).Find(&articles)
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}
