package middleware

import (
	"go-blog/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Authorization字段格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		// 解析JWT令牌
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 从Redis中获取用户会话
		token, err := util.GetUserSession(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户会话已过期"})
			c.Abort()
			return
		}

		// 验证令牌是否匹配
		if token != parts[1] {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 设置用户信息到上下文中
		c.Set("user", map[string]interface{}{
			"ID":       claims.UserID,
			"Username": claims.Username,
		})

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Authorization字段格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Next()
			return
		}

		// 解析JWT令牌
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		// 从Redis中获取用户会话
		token, err := util.GetUserSession(claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		// 验证令牌是否匹配
		if token != parts[1] {
			c.Next()
			return
		}

		// 设置用户信息到上下文中
		c.Set("user", map[string]interface{}{
			"ID":       claims.UserID,
			"Username": claims.Username,
		})

		c.Next()
	}
}
