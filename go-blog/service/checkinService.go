package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostCheckin(c *gin.Context) {
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

	today := time.Now().Format("2006-01-02")
	todayTime, _ := time.Parse("2006-01-02", today)

	var existingCheckin model.Checkin
	result := util.Db.Where("user_id = ? AND checkin_date = ?", userID, todayTime).First(&existingCheckin)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "今日已打卡"})
		return
	}

	var dbUser model.User
	util.Db.First(&dbUser, userID)

	checkinDays := 1
	expGained := 10
	coinsGained := 5

	if dbUser.LastCheckinAt != nil {
		lastCheckin := dbUser.LastCheckinAt.Format("2006-01-02")
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if lastCheckin == yesterday {
			checkinDays = dbUser.CheckinDays + 1
			expGained = 10 + checkinDays
			coinsGained = 5 + checkinDays
		}
	}

	checkin := model.Checkin{
		UserID:      userID,
		CheckinDate: todayTime,
		ExpGained:   expGained,
		CoinsGained: coinsGained,
	}
	util.Db.Create(&checkin)

	now := time.Now()
	util.Db.Model(&dbUser).Updates(map[string]interface{}{
		"checkin_days":    checkinDays,
		"last_checkin_at": &now,
		"exp":             dbUser.Exp + expGained,
		"coins":           dbUser.Coins + coinsGained,
	})

	newLevel := dbUser.Level
	for expNeeded := dbUser.Level * 100; dbUser.Exp+expGained >= expNeeded; expNeeded = newLevel * 100 {
		newLevel++
	}
	if newLevel > dbUser.Level {
		util.Db.Model(&dbUser).Update("level", newLevel)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "打卡成功",
		"checkin_days": checkinDays,
		"exp_gained":   expGained,
		"coins_gained": coinsGained,
		"level":        newLevel,
	})
}

func GetCheckinStatus(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"checked_in": false})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusOK, gin.H{"checked_in": false})
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

	today := time.Now().Format("2006-01-02")
	todayTime, _ := time.Parse("2006-01-02", today)

	var existingCheckin model.Checkin
	result := util.Db.Where("user_id = ? AND checkin_date = ?", userID, todayTime).First(&existingCheckin)

	c.JSON(http.StatusOK, gin.H{
		"checked_in": result.Error == nil,
	})
}

func GetCheckinRank(c *gin.Context) {
	var users []model.User
	util.Db.Where("checkin_days > 0").Order("checkin_days desc").Limit(10).Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}
