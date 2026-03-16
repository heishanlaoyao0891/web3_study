package main

import (
	"go-blog/model"
	"go-blog/router"
	"go-blog/util"
)

func main() {
	// 1. 初始化数据库
	util.CreateTestDB(".env_local")
	println("数据库初始化成功！")
	// 自动迁移（创建表，对应Java的DDL脚本/Gorm自动建表）
	// 注意：生产环境慎用，开发环境方便快捷
	err := util.Db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Article{},
	)
	if err != nil {
		panic("表创建失败：" + err.Error())
	}
	// 2. 初始化路由
	r := router.InitRouter()

	// 3. 启动服务
	println("服务启动成功：http://localhost:8080")
	r.Run(":8080")
}
