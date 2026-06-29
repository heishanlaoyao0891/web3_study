package util

import (
	"go-blog/model"

	"github.com/gin-gonic/gin"
)

func UintPtr(v uint) *uint {
	return &v
}

// contextUserKey 存入 gin.Context 的用户键
const contextUserKey = "currentUser"

// GetUserFromContext 从上下文中获取当前登录用户，未登录返回 nil
// 一次请求内最多查库一次，后续从 context 缓存读取
func GetUserFromContext(c *gin.Context) *model.User {
	// 优先从本次请求缓存读
	if v, exists := c.Get(contextUserKey); exists {
		if u, ok := v.(*model.User); ok {
			return u
		}
	}

	// 从 cookie 取 JWT
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		return nil
	}

	claims, err := ParseToken(token)
	if err != nil {
		return nil
	}

	// 校验 Redis 会话（顺带续期）
	tokenFromRedis, err := GetUserSession(claims.UserID)
	if err != nil || tokenFromRedis != token {
		return nil
	}

	// 查库确认用户状态
	var user model.User
	if result := Db.First(&user, claims.UserID); result.Error != nil {
		return nil
	}
	if user.Status != 1 {
		DeleteUserSession(claims.UserID)
		c.SetCookie("token", "", -1, "/", "", false, true)
		return nil
	}

	// 缓存 + 续期 cookie
	c.Set(contextUserKey, &user)
	c.SetCookie("token", token, 30*60, "/", "", false, true)
	return &user
}

// IsAdmin 判断当前用户是否管理员（基于 Role 字段）
func IsAdmin(c *gin.Context) bool {
	u := GetUserFromContext(c)
	return u != nil && u.Role == "admin"
}

// CurrentUserID 获取当前用户 ID，未登录返回 0
func CurrentUserID(c *gin.Context) uint {
	u := GetUserFromContext(c)
	if u == nil {
		return 0
	}
	return u.ID
}

// RequireAuth 要求用户登录，返回用户指针；未登录重定向到登录页
func RequireAuth(c *gin.Context) (*model.User, bool) {
	u := GetUserFromContext(c)
	if u == nil {
		c.Redirect(302, "/login")
		return nil, false
	}
	return u, true
}