package router

import (
	"go-blog/service"
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由（对应Java的@RequestMapping）
func InitRouter() *gin.Engine {
	r := gin.Default()

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
	})

	// 加载模板文件（第三步用）
	r.LoadHTMLGlob("templates/*")

	// 首页路由
	r.GET("/", service.GetIndex)

	// 文章相关路由
	articleGroup := r.Group("/article")
	{
		articleGroup.GET("/list", service.GetArticleList)         // 文章列表
		articleGroup.GET("/detail/:id", service.GetArticleDetail) // 文章详情
	}

	return r
}
