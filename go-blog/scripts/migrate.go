package main

// M4.2 独立迁移脚本
// 用法：go run scripts/migrate.go
// 生产环境用此脚本替代 main.go 里的 AutoMigrate
// main.go 中的 AutoMigrate 应在 M4.2 完成后关闭或改为开发环境限定
//
// 说明：scripts/ 目录下原已有多个 main 包文件（main 冲突），
// 运行此脚本时需指定文件：go run scripts/migrate.go

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
	"os"
)

func main() {
	// 连接数据库
	envFile := ".env_prod"
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}

	fmt.Printf("正在连接数据库 (env: %s)...\n", envFile)
	if err := util.CreateTestDB(envFile); err != nil {
		fmt.Printf("❌ 数据库连接失败：%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("✅ 数据库连接成功")

	// 按依赖顺序迁移（无 FK 冲突的表前置）
	fmt.Println("开始迁移...")

	tables := []interface{}{
		// 独立表先建
		&model.SiteConfig{},
		&model.ContentSource{},
		&model.CrawlLog{},
		&model.User{},
		&model.Category{},
		&model.Tag{},
		// 依赖上述表的
		&model.Article{},
		&model.Comment{},
		&model.ArticleTag{},
		&model.ArticleCategory{},
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
		&model.InterviewQuestionCategory{},
	}

	for _, table := range tables {
		if err := util.Db.AutoMigrate(table); err != nil {
			fmt.Printf("⚠️  迁移警告: %T - %s\n", table, err.Error())
		} else {
			fmt.Printf("✅ %T\n", table)
		}
	}

	fmt.Println("迁移完成！")
}