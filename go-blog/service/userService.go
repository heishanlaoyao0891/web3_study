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
			// 永久禁用
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error":        "账号已被永久禁用",
				"is_disabled":  true,
				"is_permanent": true,
			})
		} else {
			// 限时禁用
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error":         "账号已被禁用",
				"is_disabled":   true,
				"is_permanent":  false,
				"disable_until": user.DisableUntil,
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

	// 设置令牌到cookie中，过期时间为30分钟
	c.SetCookie("token", token, 30*60, "/", "localhost", false, true)

	// 重定向到首页
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

	// 检查用户名是否已存在
	var user model.User
	result := util.Db.Where("username = ?", username).First(&user)
	if result.Error == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "用户名已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "注册失败",
		})
		return
	}

	// 创建新用户
	newUser := model.User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
	}
	util.Db.Create(&newUser)

	// 注册成功，跳转到登录页面
	c.Redirect(http.StatusFound, "/login")
}

// GetUserList 显示用户列表页面（仅管理员可见）
func GetUserList(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
		})
		return
	}

	// 获取所有用户
	var users []model.User
	util.Db.Find(&users)

	c.HTML(http.StatusOK, "user_list.html", gin.H{
		"users": users,
		"user":  user,
	})
}

// PostDisableUser 禁用用户（仅管理员可见）
func PostDisableUser(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
		})
		return
	}

	// 获取用户ID和禁用时间
	id := c.PostForm("id")
	disableTimeStr := c.PostForm("disable_time")

	// 找到对应的用户
	var targetUser model.User
	result := util.Db.First(&targetUser, id)
	if result.Error != nil {
		// 获取所有用户
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
	if targetUser.Username == "admin" {
		// 获取所有用户
		var users []model.User
		util.Db.Find(&users)
		c.HTML(http.StatusOK, "user_list.html", gin.H{
			"users": users,
			"user":  user,
			"error": "不能禁用管理员",
		})
		return
	}

	// 计算禁用结束时间
	var disableUntil *time.Time
	disableTime, err := strconv.Atoi(disableTimeStr)
	if err != nil || disableTime <= 0 {
		// 永久禁用
		disableUntil = nil
	} else {
		// 按分钟禁用
		endTime := time.Now().Add(time.Duration(disableTime) * time.Minute)
		disableUntil = &endTime
	}

	// 更新用户状态和禁用结束时间
	targetUser.Status = 0
	targetUser.DisableUntil = disableUntil
	util.Db.Save(&targetUser)

	// 从Redis中删除用户会话，踢出用户下线
	util.DeleteUserSession(targetUser.ID)

	// 重定向到用户列表页面
	c.Redirect(http.StatusFound, "/user/list")
}

// GetLogout 处理登出请求
func GetLogout(c *gin.Context) {
	// 从cookie中获取令牌
	token, err := c.Cookie("token")

	if err == nil && token != "" {
		// 解析JWT令牌
		claims, err := util.ParseToken(token)
		if err == nil {
			// 从Redis中删除用户会话
			util.DeleteUserSession(claims.UserID)
		}
	}

	// 清除cookie中的令牌
	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	// 重定向到首页
	c.Redirect(http.StatusFound, "/")
}

// PostRestoreUser 恢复用户（仅管理员可见）
func PostRestoreUser(c *gin.Context) {
	// 从上下文中获取用户信息
	user := util.GetUserFromContext(c)

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 检查是否是管理员
	userMap, ok := user.(map[string]interface{})
	if !ok || userMap["Username"] != "admin" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "权限不足",
		})
		return
	}

	// 获取用户ID
	id := c.PostForm("id")

	// 找到对应的用户
	var targetUser model.User
	result := util.Db.First(&targetUser, id)
	if result.Error != nil {
		// 获取所有用户
		var users []model.User
		util.Db.Find(&users)
		c.HTML(http.StatusOK, "user_list.html", gin.H{
			"users": users,
			"user":  user,
			"error": "用户不存在",
		})
		return
	}

	// 更新用户状态为正常
	targetUser.Status = 1
	targetUser.DisableUntil = nil
	util.Db.Save(&targetUser)

	// 重定向到用户列表页面
	c.Redirect(http.StatusFound, "/user/list")
}
