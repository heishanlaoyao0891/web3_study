package router

import (
	"go-blog/service"
	"go-blog/util"
	"html/template"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func mul(a, b int) int {
	return a * b
}

func div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func iterate(count int) []int {
	var items []int
	for i := 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// InitRouter 初始化路由（对应Java的@RequestMapping）
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置session中间件
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(loggerMiddleware())
	r.Use(recoveryMiddleware())
	r.Use(authMiddleware()) // 全局认证中间件
	r.Use(corsMiddleware()) // 全局跨域中间件
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
		// 数学运算
		"add": add,
		"sub": sub,
		"mul": mul,
		"div": div,
		// 迭代器
		"iterate": iterate,
		// 字符串截取
		"slice": func(s string, start, end int) string {
			if len(s) < end {
				return s
			}
			if start < 0 {
				start = 0
			}
			if end > len(s) {
				end = len(s)
			}
			if start >= end {
				return ""
			}
			return s[start:end]
		},
	})

	// 加载模板文件
	r.LoadHTMLGlob("templates/*")

	// 首页路由
	r.GET("/", service.GetIndex)

	// 文章相关路由
	articleGroup := r.Group("/article")
	{
		articleGroup.GET("/list", service.GetArticleList)         // 文章列表（不需要认证）
		articleGroup.GET("/detail/:id", service.GetArticleDetail) // 文章详情（不需要认证）
		// 发布和编辑文章需要认证
		articleGroup.GET("/create", requireAuthMiddleware(), service.GetArticleCreate)     // 发布文章页面
		articleGroup.POST("/create", requireAuthMiddleware(), service.PostArticleCreate)   // 处理发布文章请求
		articleGroup.GET("/edit/:id", requireAuthMiddleware(), service.GetArticleEdit)     // 编辑文章页面
		articleGroup.POST("/edit/:id", requireAuthMiddleware(), service.PostArticleEdit)   // 处理编辑文章请求
		articleGroup.POST("/delete/:id", requireAuthMiddleware(), service.PostArticleDelete) // 删除文章
	}

	// 评论相关路由
	commentGroup := r.Group("/comment")
	{
		commentGroup.POST("/create", requireAuthMiddleware(), service.PostCommentCreate)   // 发布评论
		commentGroup.POST("/delete/:id", requireAuthMiddleware(), service.PostCommentDelete) // 删除评论
	}

	// 用户相关路由（需要认证）
	userGroup := r.Group("/user", requireAuthMiddleware())
	{
		userGroup.GET("/list", service.GetUserList)         // 用户列表（管理员）
		userGroup.POST("/disable", service.PostDisableUser) // 禁用用户（管理员）
		userGroup.POST("/restore", service.PostRestoreUser) // 恢复用户（管理员）
	}

	// 用户资料路由
	r.GET("/profile", requireAuthMiddleware(), service.GetProfile)
	r.POST("/profile", requireAuthMiddleware(), service.PostProfile)

	// 类别相关路由（需要认证）
	categoryGroup := r.Group("/category", requireAuthMiddleware())
	{
		categoryGroup.GET("/list", service.GetCategoryList)           // 类别列表（管理员）
		categoryGroup.GET("/create", service.GetCategoryCreate)       // 添加类别页面（管理员）
		categoryGroup.POST("/create", service.PostCategoryCreate)     // 处理添加类别请求（管理员）
		categoryGroup.GET("/edit/:id", service.GetCategoryEdit)       // 编辑类别页面（管理员）
		categoryGroup.POST("/edit/:id", service.PostCategoryEdit)     // 处理编辑类别请求（管理员）
		categoryGroup.POST("/delete/:id", service.PostCategoryDelete) // 删除分类（管理员）
	}

	// 登录注册路由
	r.GET("/login", service.GetLogin)                            // 登录页面（不需要认证）
	r.POST("/login", service.PostLogin)                          // 处理登录请求（不需要认证）
	r.GET("/register", service.GetRegister)                      // 注册页面（不需要认证）
	r.POST("/register", service.PostRegister)                    // 处理注册请求（不需要认证）
	r.GET("/logout", requireAuthMiddleware(), service.GetLogout) // 登出（需要认证）

	// API 路由
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/tags", service.GetTagList)
		apiGroup.GET("/tags/hot", service.GetHotTags)
		apiGroup.POST("/like", requireAuthMiddleware(), service.ToggleLike)
		apiGroup.POST("/favorite", requireAuthMiddleware(), service.ToggleFavorite)
		apiGroup.GET("/favorites", requireAuthMiddleware(), service.GetUserFavorites)
		apiGroup.POST("/checkin", requireAuthMiddleware(), service.PostCheckin)
		apiGroup.GET("/checkin/status", service.GetCheckinStatus)
		apiGroup.GET("/checkin/rank", service.GetCheckinRank)
		apiGroup.GET("/trending", service.GetTrendingTopicsAPI) // M3.1 风口话题
	}

	// 问答路由
	qaGroup := r.Group("/qa")
	{
		qaGroup.GET("/list", service.GetQuestionList)
		qaGroup.GET("/detail/:id", service.GetQuestionDetail)
		qaGroup.GET("/create", requireAuthMiddleware(), service.GetQuestionCreate)
		qaGroup.POST("/create", requireAuthMiddleware(), service.PostQuestionCreate)
		qaGroup.POST("/answer", requireAuthMiddleware(), service.PostAnswerCreate)
		qaGroup.POST("/accept/:id", requireAuthMiddleware(), service.PostAcceptAnswer)
	}

	// 教程路由
	courseGroup := r.Group("/course")
	{
		courseGroup.GET("/list", service.GetCourseList)
		courseGroup.GET("/detail/:id", service.GetCourseDetail)
		courseGroup.GET("/:course_id/chapter/:chapter_id", service.GetChapterDetail)
	}

	// Web3学习路径
	r.GET("/learning-paths", service.GetLearningPaths)
	r.GET("/learning-paths/:id", service.GetLearningPathDetail)

	// 代码片段
	r.GET("/snippets", service.GetCodeSnippets)

	// 合约模板
	r.GET("/templates", service.GetContractTemplates)

	// 资源导航
	r.GET("/resources", service.GetResources)

	// 面试题库
	r.GET("/interview", service.GetInterviewQuestions)
	r.GET("/interview/:id", service.GetInterviewQuestionDetail)

	// 风口话题管理路由（仅管理员）
	trendingGroup := r.Group("/trending", requireAuthMiddleware())
	{
		trendingGroup.GET("/list", service.GetTrendingAdminPage)
		trendingGroup.POST("/create", service.PostTrendingCreate)
		trendingGroup.POST("/delete/:id", service.PostTrendingDelete)
		trendingGroup.POST("/offline/:id", service.PostTrendingOffline)
	}

	// 抓取源管理路由（仅管理员 M2.2）
	sourceGroup := r.Group("/sources", requireAuthMiddleware())
	{
		sourceGroup.GET("", service.GetContentSources)
		sourceGroup.POST("/toggle/:id", service.PostContentSourceToggle)
		sourceGroup.POST("/delete/:id", service.PostContentSourceDelete)
		sourceGroup.POST("/run/:id", service.PostContentSourceRunNow)
	}

	// 健康检查（M4.4 无鉴权）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "go-blog",
			"time":    time.Now().Unix(),
		})
	})

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
		slog.Info("HTTP",
			"method", method,
			"path", path,
			"status", status,
			"ip", clientIP,
			"duration", since.String(),
			"raw", raw,
		)
	}
}

// 全局异常中间件
func recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		slog.Error("Panic recovered",
			"error", recovered,
			"path", c.Request.URL.Path,
			"stack", string(debug.Stack()),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
	})
}

// 认证中间件：将当前用户注入到上下文（未登录不阻断）
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := util.GetUserFromContext(c)
		c.Set("user", user)
		c.Next()
	}
}

// 要求认证中间件（用于需要登录的路由）
func requireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := util.GetUserFromContext(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}