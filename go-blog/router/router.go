package router

import (
	"go-blog/service"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由（对应Java的@RequestMapping）
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置session中间件
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(loggerMiddleware())
	r.Use(recoveryMiddleware())
	// 注册自定义模板函数
	r.SetFuncMap(template.FuncMap{
		// 替换字符串（解决文章内容换行）
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		// 安全HTML（允许输出<br/>等标签）
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		// 获取当前时间
		"now": func() time.Time {
			return time.Now()
		},
	})

	// 加载模板文件
	r.LoadHTMLGlob("templates/*")

	// 首页路由
	r.GET("/", service.GetIndex)

	// 文章相关路由
	articleGroup := r.Group("/article")
	{
		articleGroup.GET("/list", service.GetArticleList)         // 文章列表
		articleGroup.GET("/detail/:id", service.GetArticleDetail) // 文章详情
		articleGroup.GET("/create", service.GetArticleCreate)     // 发布文章页面
		articleGroup.POST("/create", service.PostArticleCreate)   // 处理发布文章请求
		articleGroup.GET("/edit/:id", service.GetArticleEdit)     // 编辑文章页面
		articleGroup.POST("/edit/:id", service.PostArticleEdit)   // 处理编辑文章请求
	}

	// 用户相关路由
	userGroup := r.Group("/user")
	{
		userGroup.GET("/list", service.GetUserList)         // 用户列表（管理员）
		userGroup.POST("/disable", service.PostDisableUser) // 禁用用户（管理员）
		userGroup.POST("/restore", service.PostRestoreUser) // 恢复用户（管理员）
	}

	// 类别相关路由
	categoryGroup := r.Group("/category")
	{
		categoryGroup.GET("/list", service.GetCategoryList)       // 类别列表（管理员）
		categoryGroup.GET("/create", service.GetCategoryCreate)   // 添加类别页面（管理员）
		categoryGroup.POST("/create", service.PostCategoryCreate) // 处理添加类别请求（管理员）
		categoryGroup.GET("/edit/:id", service.GetCategoryEdit)   // 编辑类别页面（管理员）
		categoryGroup.POST("/edit/:id", service.PostCategoryEdit) // 处理编辑类别请求（管理员）
	}

	// 登录注册路由
	r.GET("/login", service.GetLogin)         // 登录页面
	r.POST("/login", service.PostLogin)       // 处理登录请求
	r.GET("/register", service.GetRegister)   // 注册页面
	r.POST("/register", service.PostRegister) // 处理注册请求
	r.GET("/logout", service.GetLogout)       // 登出

	return r
}

// 日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		since := time.Since(start)
		clientIP := c.ClientIP()
		status := c.Writer.Status()
		log.Printf("[%s] %s %s %s %d %v\n", clientIP, path, raw, method, status, since)
	}
}

// 全局异常中间件
func recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录异常信息
		log.Printf("Panic recovered: %v", recovered)
		// 可以添加堆栈信息
		log.Printf("Stack trace: %s", debug.Stack())

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
	})
}
