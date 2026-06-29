package main

import (
	"go-blog/crawler"
	"go-blog/model"
	"go-blog/router"
	"go-blog/util"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// initData 初始化默认数据
func initData() {
	// 检查是否已经存在admin用户
	var user model.User
	result := util.Db.Where("username = ?", "admin").First(&user)
	if result.Error != nil {
		password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		admin := model.User{
			Username: "admin",
			Password: string(password),
			Nickname: "管理员",
			Role:     "admin",
		}
		util.Db.Create(&admin)
	} else if user.Role != "admin" {
		user.Role = "admin"
		util.Db.Save(&user)
	}

	// 默认站点配置（M1.3）
	initSiteConfig()

	// 默认技术领域（顶层分类）
	initDefaultDomains()

	// 默认抓取源（M2 打样）
	initDefaultSources()

	println("初始化数据成功！")
}

// initSiteConfig 写入默认站点配置
func initSiteConfig() {
	defaults := map[string]string{
		model.CfgSiteTitle:      "技术学习平台",
		model.CfgSiteSubtitle:   "Java | Go | Python | AI | Web3 | 全栈技术社区",
		model.CfgSiteIntroTitle: "技术开发者社区",
		model.CfgFooterText:     "© 2026 技术学习平台 | 基于 Go + Gin + GORM 开发",
		model.CfgFooterSubtext:  "专注全栈技术学习 | 面试题库 | 资讯推送",
	}
	for key, val := range defaults {
		var cfg model.SiteConfig
		if err := util.Db.Where("config_key = ?", key).First(&cfg).Error; err != nil {
			util.Db.Create(&model.SiteConfig{Key: key, Value: val, Category: "brand", Desc: key})
		}
	}
}

// initDefaultDomains 初始化默认技术领域（顶层分类）
func initDefaultDomains() {
	domains := []model.Category{
		{Name: "Web3", Desc: "区块链与智能合约开发", Icon: "🔗", SortOrder: 50},
		{Name: "Java", Desc: "Java后端开发", Icon: "☕", SortOrder: 40},
		{Name: "Go", Desc: "Go语言开发", Icon: "🐹", SortOrder: 30},
		{Name: "Python", Desc: "Python开发", Icon: "🐍", SortOrder: 20},
		{Name: "AI", Desc: "人工智能与机器学习", Icon: "🤖", SortOrder: 10},
	}
	for _, d := range domains {
		var existing model.Category
		if err := util.Db.Where("name = ?", d.Name).First(&existing).Error; err != nil {
			util.Db.Create(&d)
		}
	}
}

// initDefaultSources 初始化默认抓取源
func initDefaultSources() {
	// 获取各领域 ID
	domains := make(map[string]uint)
	var cats []model.Category
	util.Db.Where("parent_id IS NULL").Find(&cats)
	for _, c := range cats {
		domains[c.Name] = c.ID
	}
	if len(domains) == 0 {
		return
	}

	// 默认 Hacker News 源
	var existing model.ContentSource
	if err := util.Db.Where("type = ? AND name = ?", "hackernews", "Hacker News").First(&existing).Error; err != nil {
		if aiID, ok := domains["AI"]; ok {
			util.Db.Create(&model.ContentSource{
				Name:     "Hacker News",
				Type:     "hackernews",
				URL:      "https://hacker-news.firebaseio.com/v0/topstories.json",
				DomainID: aiID,
				Cron:     "0 * * * *", // 每小时整点
				Enabled:  true,
			})
		}
	}

	// 默认 RSS 订阅源
	rssFeeds := []struct {
		Name     string
		URL      string
		Domain   string
		Cron     string
	}{
		{"Go 官方博客", "https://go.dev/blog/feed.atom", "Go", "0 */2 * * *"},
		{"InfoQ 中文", "https://www.infoq.cn/feed", "Java", "0 */2 * * *"},
		{"开源中国", "https://www.oschina.net/news/rss", "Python", "0 */3 * * *"},
		{"Ethereum 博客", "https://blog.ethereum.org/en/feed.xml", "Web3", "0 */2 * * *"},
		{"Solidity 博客", "https://soliditylang.org/blog/feed.xml", "Web3", "0 */3 * * *"},
	}
	for _, f := range rssFeeds {
		if err := util.Db.Where("type = ? AND name = ?", "rss", f.Name).First(&existing).Error; err != nil {
			if domainID, ok := domains[f.Domain]; ok {
				util.Db.Create(&model.ContentSource{
					Name:     f.Name,
					Type:     "rss",
					URL:      f.URL,
					DomainID: domainID,
					Cron:     f.Cron,
					Enabled:  true,
				})
			}
		}
	}

	// 默认掘金源
	if err := util.Db.Where("type = ? AND name = ?", "juejin", "掘金热门").First(&existing).Error; err != nil {
		if di, ok := domains["Go"]; ok {
			util.Db.Create(&model.ContentSource{
				Name:     "掘金热门",
				Type:     "juejin",
				URL:      "https://juejin.cn",
				DomainID: di,
				Cron:     "0 */2 * * *",
				Enabled:  true,
			})
		}
	}

	// 默认知乎源
	if err := util.Db.Where("type = ? AND name = ?", "zhihu", "知乎热榜").First(&existing).Error; err != nil {
		if di, ok := domains["AI"]; ok {
			util.Db.Create(&model.ContentSource{
				Name:     "知乎热榜",
				Type:     "zhihu",
				URL:      "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total",
				DomainID: di,
				Cron:     "0 */2 * * *",
				Enabled:  true,
			})
		}
	}
}

// startAutoRestoreTask 定时任务：自动恢复被禁用的用户
func startAutoRestoreTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			var users []model.User
			util.Db.Where("status = ? AND disable_until IS NOT NULL AND disable_until < ?", 0, time.Now()).Find(&users)

			for _, user := range users {
				user.Status = 1
				user.DisableUntil = nil
				util.Db.Save(&user)
			}
		}
	}()
}

func main() {
	// 0. 初始化结构化日志（M4.3）
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// 1. 初始化数据库
	if err := util.CreateTestDB(".env_prod"); err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	println("数据库初始化成功！")

	// 2. 自动迁移（M4.2：仅开发环境执行，生产环境请用 scripts/migrate.go）
	// 环境变量 APP_ENV=production 或 GO_ENV=production 时跳过 AutoMigrate
	var err error
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "production" || env == "prod" {
		slog.Info("生产环境：跳过 AutoMigrate，请使用 'go run scripts/migrate.go .env_prod' 执行迁移")
	} else {
		slog.Info("开发环境：执行 AutoMigrate")
		err = util.Db.AutoMigrate(
			&model.SiteConfig{},    // M1.3 前置
			&model.ContentSource{}, // M2.2 前置
			&model.CrawlLog{},      // M2.2 前置
			&model.User{},
			&model.Category{},
			&model.Article{},
			&model.Comment{},
			&model.Tag{},
			&model.ArticleTag{},
			&model.ArticleCategory{},
			&model.InterviewQuestionCategory{},
			&model.Like{},
			&model.Favorite{},
			&model.Question{},
			&model.Answer{},
			&model.Course{},
			&model.CourseChapter{},
			&model.LearningRecord{},
			&model.Checkin{},
			&model.LearningPath{},
			&model.LearningChapter{},
			&model.CodeSnippet{},
			&model.ContractTemplate{},
			&model.Resource{},
			&model.InterviewQuestion{},
			&model.TrendingTopic{}, // M3.1 风口话题
		)
		if err != nil {
			slog.Warn("表迁移警告：" + err.Error())
		}
	}

	// 3. 初始化Redis
	err = util.InitRedis(".env_prod")
	if err != nil {
		println("Redis初始化失败：" + err.Error())
	}

	// 4. 初始化测试数据
	initData()

	// 5. 启动自动恢复任务
	startAutoRestoreTask()
	println("自动恢复任务启动成功！")

	// 6. 启动抓取调度器（M2.1）
	crawler.GetScheduler().Start()

	// 7. 初始化路由
	r := router.InitRouter()

	// 8. 启动 HTTP 服务（带优雅关停）
	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	// 后台启动服务
	go func() {
		println("服务启动成功：http://localhost:8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP 服务启动失败", "error", err)
		}
	}()

	// 等待中断信号，优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("收到关停信号", "signal", sig.String())

	// 停止抓取调度器
	crawler.GetScheduler().Stop()

	// 给正在处理的请求最多 10 秒完成
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("服务关停异常", "error", err)
	}
	slog.Info("服务已安全退出")
}