package main

import (
	"go-blog/model"
	"go-blog/router"
	"go-blog/util"
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
		// 兼容旧数据：给已存在的 admin 用户补上 Role
		user.Role = "admin"
		util.Db.Save(&user)
	}

	// 默认站点配置（M1.3）
	initSiteConfig()

	// 默认技术领域（顶层分类）
	initDefaultDomains()

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
	// 1. 初始化数据库
	if err := util.CreateTestDB(".env_prod"); err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	println("数据库初始化成功！")

	// 2. 自动迁移（开发环境使用，生产环境应改为独立迁移脚本）
	err := util.Db.AutoMigrate(
		&model.SiteConfig{}, // M1.3 新增（前置，避免其他表迁移错误中断）
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
	)
	if err != nil {
		println("表迁移警告：" + err.Error())
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

	// 6. 初始化路由
	r := router.InitRouter()

	// 7. 启动服务
	println("服务启动成功：http://localhost:8081")
	r.Run(":8081")
}