package service

import (
	"go-blog/model"
	"go-blog/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetLogin 显示登录页面
func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// PostLogin 处理登录请求
func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user model.User
	result := util.Db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		if user.DisableUntil == nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error":         "账号已被永久禁用",
				"is_disabled":   true,
				"is_permanent":  true,
			})
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error":          "账号已被禁用",
				"is_disabled":    true,
				"is_permanent":   false,
				"disable_until":  user.DisableUntil,
			})
		}
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "生成令牌失败",
		})
		return
	}

	// 将令牌存储到Redis中
	err = util.SetUserSession(user.ID, token)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "存储会话失败",
		})
		return
	}

	c.SetCookie("token", token, 30*60, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

// GetRegister 显示注册页面
func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// PostRegister 处理注册请求
func PostRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")

	var user model.User
	result := util.Db.Where("username = ?", username).First(&user)
	if result.Error == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "用户名已存在",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "注册失败",
		})
		return
	}

	newUser := model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
	}
	util.Db.Create(&newUser)

	c.Redirect(http.StatusFound, "/login")
}

// GetUserList 显示用户列表页面（仅管理员可见）
func GetUserList(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	var users []model.User
	util.Db.Find(&users)

	c.HTML(http.StatusOK, "user_list.html", gin.H{
		"users": users,
		"user":  user,
	})
}

// PostDisableUser 禁用用户（仅管理员可见）
func PostDisableUser(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	id := c.PostForm("id")
	disableTimeStr := c.PostForm("disable_time")

	var targetUser model.User
	result := util.Db.First(&targetUser, id)
	if result.Error != nil {
		var users []model.User
		util.Db.Find(&users)
		c.HTML(http.StatusOK, "user_list.html", gin.H{
			"users": users,
			"user":  user,
			"error": "用户不存在",
		})
		return
	}

	// 不能禁用管理员
	if targetUser.Role == "admin" {
		var users []model.User
		util.Db.Find(&users)
		c.HTML(http.StatusOK, "user_list.html", gin.H{
			"users": users,
			"user":  user,
			"error": "不能禁用管理员",
		})
		return
	}

	var disableUntil *time.Time
	disableTime, err := strconv.Atoi(disableTimeStr)
	if err != nil || disableTime <= 0 {
		disableUntil = nil
	} else {
		endTime := time.Now().Add(time.Duration(disableTime) * time.Minute)
		disableUntil = &endTime
	}

	targetUser.Status = 0
	targetUser.DisableUntil = disableUntil
	util.Db.Save(&targetUser)

	util.DeleteUserSession(targetUser.ID)

	c.Redirect(http.StatusFound, "/user/list")
}

// GetLogout 处理登出请求
func GetLogout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil && token != "" {
		claims, err := util.ParseToken(token)
		if err == nil {
			util.DeleteUserSession(claims.UserID)
		}
	}

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

// PostRestoreUser 恢复用户（仅管理员可见）
func PostRestoreUser(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if !requireAdmin(c) {
		return
	}

	id := c.PostForm("id")

	var targetUser model.User
	result := util.Db.First(&targetUser, id)
	if result.Error != nil {
		var users []model.User
		util.Db.Find(&users)
		c.HTML(http.StatusOK, "user_list.html", gin.H{
			"users": users,
			"user":  user,
			"error": "用户不存在",
		})
		return
	}

	targetUser.Status = 1
	targetUser.DisableUntil = nil
	util.Db.Save(&targetUser)

	c.Redirect(http.StatusFound, "/user/list")
}

// GetProfile 显示用户资料页面
func GetProfile(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var dbUser model.User
	util.Db.First(&dbUser, user.ID)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": dbUser,
	})
}

// PostProfile 更新用户资料
func PostProfile(c *gin.Context) {
	user := util.GetUserFromContext(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	nickname := c.PostForm("nickname")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	var dbUser model.User
	result := util.Db.First(&dbUser, user.ID)
	if result.Error != nil {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"error": "用户不存在",
			"user":  user,
		})
		return
	}

	if nickname != "" {
		dbUser.Nickname = nickname
	}

	if oldPassword != "" && newPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(oldPassword))
		if err != nil {
			c.HTML(http.StatusOK, "profile.html", gin.H{
				"error": "原密码错误",
				"user":  dbUser,
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusOK, "profile.html", gin.H{
				"error": "密码加密失败",
				"user":  dbUser,
			})
			return
		}
		dbUser.Password = string(hashedPassword)
	}

	util.Db.Save(&dbUser)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"success": "资料更新成功",
		"user":    dbUser,
	})
}