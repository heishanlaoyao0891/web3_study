package util

import (
	"go-blog/model"

	"github.com/gin-gonic/gin"
)

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *gin.Context) interface{} {
	// 从cookie中获取令牌
	token, err := c.Cookie("token")
	var user interface{}

	if err == nil && token != "" {
		// 解析JWT令牌
		claims, err := ParseToken(token)
		if err == nil {
			// 从Redis中获取用户会话
			tokenFromRedis, err := GetUserSession(claims.UserID)
			if err == nil && tokenFromRedis == token {
				// 检查用户状态
				var dbUser model.User
				result := Db.First(&dbUser, claims.UserID)
				if result.Error == nil && dbUser.Status == 1 {
					// 设置用户信息
					user = map[string]interface{}{
						"ID":       claims.UserID,
						"Username": claims.Username,
					}
				} else {
					// 用户被禁用，删除会话
					DeleteUserSession(claims.UserID)
					// 清除cookie
					c.SetCookie("token", "", -1, "/", "", false, true)
				}
			}
		}
	}

	return user
}

// RequireAuth 要求用户登录
func RequireAuth(c *gin.Context) (interface{}, bool) {
	user := GetUserFromContext(c)
	if user == nil {
		c.Redirect(302, "/login")
		return nil, false
	}
	return user, true
}
