package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostCommentCreate(c *gin.Context) {
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
	idValue := userMap["ID"]
	switch v := idValue.(type) {
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

	articleIDStr := c.PostForm("article_id")
	content := c.PostForm("content")
	parentIDStr := c.PostForm("parent_id")

	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID无效"})
		return
	}

	var article model.Article
	if result := util.Db.First(&article, articleID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	comment := model.Comment{
		Content:   content,
		ArticleID: uint(articleID),
		UserID:    userID,
	}

	if parentIDStr != "" {
		parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err == nil {
			comment.ParentID = util.UintPtr(uint(parentID))
		}
	}

	util.Db.Create(&comment)

	c.Redirect(http.StatusFound, "/article/detail/"+articleIDStr)
}

func PostCommentDelete(c *gin.Context) {
	user := util.GetUserFromContext(c)

	id := c.Param("id")
	var comment model.Comment
	result := util.Db.First(&comment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此评论"})
		return
	}

	isAdmin := userMap["Username"] == "admin"
	var userID uint
	idValue := userMap["ID"]
	switch v := idValue.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		userID = 0
	}

	if !isAdmin && userID != comment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此评论"})
		return
	}

	articleID := comment.ArticleID
	util.Db.Delete(&comment)

	c.Redirect(http.StatusFound, "/article/detail/"+strconv.FormatUint(uint64(articleID), 10))
}
