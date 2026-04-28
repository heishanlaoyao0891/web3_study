package main

import (
	"go-blog/model"
	"go-blog/router"
	"go-blog/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func initData() {
	// 检查是否已经存在admin用户
	var user model.User
	result := util.Db.Where("username = ?", "admin").First(&user)
	if result.Error != nil {
		// 创建测试用户
		password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user := model.User{
			Username: "admin",
			Password: string(password),
			Nickname: "管理员",
		}
		util.Db.Create(&user)
	}

	// 检查是否已经存在技术分类
	var category model.Category
	result = util.Db.Where("name = ?", "技术").First(&category)
	if result.Error != nil {
		// 创建测试分类
		category := model.Category{
			Name: "技术",
			Desc: "技术相关文章",
		}
		util.Db.Create(&category)
	}

	// 检查是否已经存在测试文章
	var article model.Article
	result = util.Db.Where("title = ?", "Go语言入门").First(&article)
	if result.Error != nil {
		// 创建测试文章
		article := model.Article{
			Title:      "Go语言入门",
			Content:    "Go语言是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。\n\nGo语言的设计目标是：\n1. 简洁、快速、安全\n2. 并行、有趣、开源\n3. 内存管理、数组安全、编译迅速",
			Status:     1,
			UserID:     user.ID,
			CategoryID: category.ID,
		}
		util.Db.Create(&article)
	}

	println("初始化数据成功！")
}

// 定时任务：自动恢复被禁用的用户
func startAutoRestoreTask() {
	go func() {
		// 每隔1分钟检查一次
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			// 查找所有被禁用且禁用时间已到的用户
			var users []model.User
			util.Db.Where("status = ? AND disable_until IS NOT NULL AND disable_until < ?", 0, time.Now()).Find(&users)

			// 恢复这些用户
			for _, user := range users {
				user.Status = 1
				user.DisableUntil = nil
				util.Db.Save(&user)
				println("自动恢复用户：" + user.Username)
			}
		}
	}()
}

func main() {
	// 1. 加载环境变量
	// 注意：CreateTestDB函数已经加载了环境变量，所以Redis初始化时可以直接使用
	// 2. 初始化数据库
	util.CreateTestDB(".env_prod")
	println("数据库初始化成功！")
	// 自动迁移（创建表，对应Java的DDL脚本/Gorm自动建表）
	// 注意：生产环境慎用，开发环境方便快捷
	err := util.Db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Article{},
		&model.Comment{},
		&model.Tag{},
		&model.ArticleTag{},
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
		panic("表创建失败：" + err.Error())
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
