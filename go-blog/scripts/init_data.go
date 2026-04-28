package main

import (
	"go-blog/model"
	"go-blog/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	util.CreateTestDB(".env_prod")
	println("数据库初始化成功！")

	initUsers()
	initCategories()
	initTags()
	initArticles()
	initQuestions()
	initCourses()
	initCheckins()

	println("基础数据初始化完成！")
}

func initUsers() {
	users := []struct {
		Username string
		Password string
		Nickname string
		Level    int
		Exp      int
		Coins    int
		Bio      string
	}{
		{"admin", "123456", "管理员", 5, 1500, 500, "系统管理员，负责网站维护"},
		{"zhangsan", "123456", "张三", 3, 450, 120, "Java开发工程师，热爱技术分享"},
		{"lisi", "123456", "李四", 2, 200, 60, "前端开发，Vue技术栈"},
		{"wangwu", "123456", "王五", 4, 800, 200, "全栈工程师，擅长Go和React"},
		{"zhaoliu", "123456", "赵六", 1, 50, 20, "Python爱好者，数据分析方向"},
		{"sunqi", "123456", "孙七", 3, 500, 150, "后端架构师，微服务专家"},
		{"zhouba", "123456", "周八", 2, 300, 80, "数据库专家，MySQL/Redis"},
		{"wujiu", "123456", "吴九", 1, 80, 30, "学生，正在学习编程"},
	}

	for _, u := range users {
		var user model.User
		result := util.Db.Where("username = ?", u.Username).First(&user)
		if result.Error == nil {
			continue
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		now := time.Now()
		user = model.User{
			Username: u.Username,
			Password: string(password),
			Nickname: u.Nickname,
			Level:    u.Level,
			Exp:      u.Exp,
			Coins:    u.Coins,
			Bio:      u.Bio,
			Status:   1,
		}
		if u.Username == "zhangsan" || u.Username == "wangwu" {
			user.LastCheckinAt = &now
			user.CheckinDays = 7
		}
		util.Db.Create(&user)
	}
	println("用户数据初始化完成！")
}

func initCategories() {
	categories := []struct {
		Name string
		Desc string
	}{
		{"技术", "技术相关文章，包含编程语言、框架、工具等"},
		{"生活", "生活随笔，记录日常点滴"},
		{"职场", "职场经验，求职面试心得"},
		{"学习", "学习笔记，课程总结"},
		{"项目", "项目实战，开源作品分享"},
	}

	for _, c := range categories {
		var category model.Category
		result := util.Db.Where("name = ?", c.Name).First(&category)
		if result.Error == nil {
			continue
		}
		category = model.Category{
			Name: c.Name,
			Desc: c.Desc,
		}
		util.Db.Create(&category)
	}
	println("分类数据初始化完成！")
}

func initTags() {
	tags := []struct {
		Name  string
		Color string
	}{
		{"Java", "#f89820"},
		{"Go", "#00add8"},
		{"Python", "#3776ab"},
		{"前端", "#42b883"},
		{"后端", "#6c757d"},
		{"数据库", "#336791"},
		{"Redis", "#dc382d"},
		{"MySQL", "#4479a1"},
		{"微服务", "#ff6b6b"},
		{"Docker", "#2496ed"},
		{"Kubernetes", "#326ce5"},
		{"算法", "#ffc107"},
		{"面试", "#667eea"},
		{"求职", "#764ba2"},
		{"Vue", "#42b883"},
		{"React", "#61dafb"},
		{"Spring", "#6db33f"},
		{"Gin", "#00add8"},
		{"Linux", "#fcc624"},
		{"Git", "#f05032"},
	}

	for i, t := range tags {
		var tag model.Tag
		result := util.Db.Where("name = ?", t.Name).First(&tag)
		if result.Error == nil {
			continue
		}
		tag = model.Tag{
			Name:      t.Name,
			Color:     t.Color,
			SortOrder: i + 1,
			UseCount:  0,
		}
		util.Db.Create(&tag)
	}
	println("标签数据初始化完成！")
}

func initArticles() {
	var adminUser, zhangsanUser, lisiUser, wangwuUser model.User
	util.Db.Where("username = ?", "admin").First(&adminUser)
	util.Db.Where("username = ?", "zhangsan").First(&zhangsanUser)
	util.Db.Where("username = ?", "lisi").First(&lisiUser)
	util.Db.Where("username = ?", "wangwu").First(&wangwuUser)

	var techCat, lifeCat, workCat, studyCat model.Category
	util.Db.Where("name = ?", "技术").First(&techCat)
	util.Db.Where("name = ?", "生活").First(&lifeCat)
	util.Db.Where("name = ?", "职场").First(&workCat)
	util.Db.Where("name = ?", "学习").First(&studyCat)

	articles := []struct {
		Title      string
		Content    string
		UserID     uint
		CategoryID uint
		ViewCount  int
		LikeCount  int
	}{
		{
			Title:      "Go语言入门教程",
			Content:    "# Go语言入门\n\nGo语言是Google开发的一种静态强类型、编译型、并发型编程语言。\n\n## 特点\n\n1. 语法简洁\n2. 编译速度快\n3. 天然支持并发\n4. 内存安全\n\n## 安装\n\n访问 https://golang.org 下载对应平台的安装包。\n\n## Hello World\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n## 学习建议\n\n1. 先学习基础语法\n2. 理解goroutine和channel\n3. 学习标准库的使用\n4. 尝试写一些小项目",
			UserID:     adminUser.ID,
			CategoryID: techCat.ID,
			ViewCount:  1256,
			LikeCount:  89,
		},
		{
			Title: "Gin框架实战指南",
			Content: "# Gin框架实战\n\nGin是一个用Go语言编写的高性能Web框架。\n\n## 安装\n\n```bash\ngo get -u github.com/gin-gonic/gin\n```\n\n## 快速开始\n\n```go\npackage main\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc main() {\n    r := gin.Default()\n    r.GET(\"/ping\", func(c *gin.Context) {\n        c.JSON(200, gin.H{\n            \"message\": \"pong\",\n        })\n    })\n    r.Run(\":8080\")\n}\n```\n\n## 路由分组\n\n```go\nv1 := r.Group(\"/v1\")\n{\n    v1.GET(\"/users\", getUsers)\n    v1.POST(\"/users\", createUser)\n}\n```\n\n## 中间件\n\nGin支持中间件，可以用于日志记录、认证等功能。",
			UserID:     wangwuUser.ID,
			CategoryID: techCat.ID,
			ViewCount:  892,
			LikeCount:  67,
		},
		{
			Title: "MySQL性能优化实践",
			Content: "# MySQL性能优化\n\n## 索引优化\n\n1. 为常用查询字段创建索引\n2. 使用EXPLAIN分析查询计划\n3. 避免在索引列上使用函数\n\n## 查询优化\n\n```sql\n-- 避免SELECT *\nSELECT id, name FROM users WHERE status = 1;\n\n-- 使用LIMIT\nSELECT * FROM orders ORDER BY created_at DESC LIMIT 100;\n```\n\n## 表结构优化\n\n1. 选择合适的数据类型\n2. 避免NULL值\n3. 适当的反范式设计\n\n## 缓存优化\n\n1. 使用Redis缓存热点数据\n2. 合理设置缓存过期时间",
			UserID:     zhangsanUser.ID,
			CategoryID: techCat.ID,
			ViewCount:  567,
			LikeCount:  45,
		},
		{
			Title: "Vue3组合式API学习笔记",
			Content: "# Vue3组合式API\n\n## setup函数\n\n```javascript\nimport { ref, reactive } from 'vue'\n\nexport default {\n  setup() {\n    const count = ref(0)\n    const state = reactive({ name: 'Vue' })\n    \n    return { count, state }\n  }\n}\n```\n\n## 响应式API\n\n- ref: 用于基本类型\n- reactive: 用于对象\n- computed: 计算属性\n- watch: 侦听器\n\n## 生命周期\n\n- onMounted\n- onUpdated\n- onUnmounted",
			UserID:     lisiUser.ID,
			CategoryID: techCat.ID,
			ViewCount:  423,
			LikeCount:  34,
		},
		{
			Title: "程序员的一天",
			Content: "# 程序员的一天\n\n## 早上\n\n7:30 起床\n8:30 到公司\n9:00 晨会\n\n## 上午\n\n写代码、修bug、看文档\n\n## 中午\n\n12:00 午饭\n13:00 午休\n\n## 下午\n\n继续写代码，偶尔开会\n\n## 晚上\n\n18:00 下班\n19:00 健身\n21:00 学习新技术\n23:00 睡觉",
			UserID:     zhangsanUser.ID,
			CategoryID: lifeCat.ID,
			ViewCount:  234,
			LikeCount:  28,
		},
		{
			Title: "面试常见问题汇总",
			Content: "# 面试常见问题\n\n## Java基础\n\n1. HashMap原理\n2. JVM内存模型\n3. 垃圾回收机制\n\n## 数据库\n\n1. 索引原理\n2. 事务隔离级别\n3. MVCC实现\n\n## 网络\n\n1. TCP三次握手\n2. HTTP与HTTPS\n3. 浏览器输入URL后发生什么\n\n## 算法\n\n1. 反转链表\n2. 二叉树遍历\n3. 动态规划",
			UserID:     wangwuUser.ID,
			CategoryID: workCat.ID,
			ViewCount:  1567,
			LikeCount:  123,
		},
		{
			Title: "Docker入门教程",
			Content: "# Docker入门\n\n## 什么是Docker\n\nDocker是一个开源的容器化平台。\n\n## 安装\n\n```bash\n# Ubuntu\nsudo apt-get install docker.io\n\n# Mac\nbrew install docker\n```\n\n## 常用命令\n\n```bash\n# 构建镜像\ndocker build -t myapp .\n\n# 运行容器\ndocker run -d -p 8080:80 myapp\n\n# 查看容器\ndocker ps\n\n# 进入容器\ndocker exec -it container_id /bin/bash\n```\n\n## Dockerfile\n\n```dockerfile\nFROM golang:1.20\nWORKDIR /app\nCOPY . .\nRUN go build -o main .\nCMD [\"./main\"]\n```",
			UserID:     adminUser.ID,
			CategoryID: techCat.ID,
			ViewCount:  678,
			LikeCount:  56,
		},
		{
			Title: "如何高效学习编程",
			Content: "# 如何高效学习编程\n\n## 明确目标\n\n- 确定学习方向\n- 制定学习计划\n- 设定阶段性目标\n\n## 学习方法\n\n1. 看官方文档\n2. 跟着教程做项目\n3. 阅读优秀源码\n4. 参与开源项目\n\n## 实践\n\n- 多写代码\n- 解决实际问题\n- 记录学习笔记\n\n## 坚持\n\n每天至少写1小时代码",
			UserID:     zhangsanUser.ID,
			CategoryID: studyCat.ID,
			ViewCount:  345,
			LikeCount:  42,
		},
	}

	for _, a := range articles {
		var article model.Article
		result := util.Db.Where("title = ?", a.Title).First(&article)
		if result.Error == nil {
			continue
		}
		article = model.Article{
			Title:         a.Title,
			Content:       a.Content,
			UserID:        a.UserID,
			CategoryID:    a.CategoryID,
			ViewCount:     a.ViewCount,
			LikeCount:     a.LikeCount,
			FavoriteCount: a.LikeCount / 3,
			CommentCount:  a.LikeCount / 5,
			Status:        1,
			Visibility:    1,
		}
		util.Db.Create(&article)
	}
	println("文章数据初始化完成！")
}

func initQuestions() {
	var zhangsanUser, lisiUser, wangwuUser model.User
	util.Db.Where("username = ?", "zhangsan").First(&zhangsanUser)
	util.Db.Where("username = ?", "lisi").First(&lisiUser)
	util.Db.Where("username = ?", "wangwu").First(&wangwuUser)

	questions := []struct {
		Title   string
		Content string
		UserID  uint
	}{
		{
			Title:   "Go语言如何实现优雅关闭？",
			Content: "在Go的Web服务中，如何实现优雅关闭？希望能够在收到关闭信号时，等待正在处理的请求完成后再退出进程。",
			UserID:  zhangsanUser.ID,
		},
		{
			Title:   "MySQL索引什么时候会失效？",
			Content: "我在面试中被问到MySQL索引失效的场景，想了解具体哪些情况会导致索引无法使用？",
			UserID:  lisiUser.ID,
		},
		{
			Title:   "Vue3的ref和reactive有什么区别？",
			Content: "Vue3中ref和reactive都是响应式API，它们之间有什么区别？什么场景下使用哪个更好？",
			UserID:  wangwuUser.ID,
		},
		{
			Title:   "如何设计一个高并发系统？",
			Content: "假设要设计一个秒杀系统，QPS可能达到10万+，应该如何设计架构？需要考虑哪些问题？",
			UserID:  zhangsanUser.ID,
		},
		{
			Title:   "Redis如何保证数据一致性？",
			Content: "在使用Redis作为缓存时，如何保证缓存和数据库的数据一致性？有哪些常见的策略？",
			UserID:  wangwuUser.ID,
		},
	}

	for _, q := range questions {
		var question model.Question
		result := util.Db.Where("title = ?", q.Title).First(&question)
		if result.Error == nil {
			continue
		}
		question = model.Question{
			Title:      q.Title,
			Content:    q.Content,
			UserID:     q.UserID,
			ViewCount:  100 + int(time.Now().Unix()%200),
			AnswerCount: 0,
			LikeCount:  10 + int(time.Now().Unix()%30),
			Status:     0,
		}
		util.Db.Create(&question)
	}

	var q1 model.Question
	util.Db.Where("title = ?", "Go语言如何实现优雅关闭？").First(&q1)
	answer1 := model.Answer{
		QuestionID: q1.ID,
		UserID:     wangwuUser.ID,
		Content:    "可以使用os.Signal配合context来实现优雅关闭：\n\n```go\nimport (\n    \"context\"\n    \"os\"\n    \"os/signal\"\n    \"syscall\"\n)\n\nfunc main() {\n    r := gin.Default()\n    // ...路由配置\n    \n    srv := &http.Server{Addr: \":8080\", Handler: r}\n    \n    go srv.ListenAndServe()\n    \n    quit := make(chan os.Signal, 1)\n    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n    <-quit\n    \n    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n    defer cancel()\n    srv.Shutdown(ctx)\n}\n```",
		LikeCount:  15,
	}
	util.Db.Create(&answer1)
	util.Db.Model(&q1).Update("answer_count", 1)

	println("问答数据初始化完成！")
}

func initCourses() {
	var adminUser model.User
	util.Db.Where("username = ?", "admin").First(&adminUser)

	courses := []struct {
		Title       string
		Description string
		Category    string
		AuthorID    uint
		Chapters    []struct {
			Title   string
			Content string
		}
	}{
		{
			Title:       "Go语言从入门到精通",
			Description: "全面学习Go语言，从基础语法到并发编程，适合零基础学员",
			Category:    "Go",
			AuthorID:    adminUser.ID,
			Chapters: []struct {
				Title   string
				Content string
			}{
				{"Go语言简介与环境搭建", "# Go语言简介\n\nGo语言是Google开发的开源编程语言...\n\n## 安装Go\n\n1. 下载安装包\n2. 配置环境变量"},
				{"基础语法", "# 基础语法\n\n## 变量声明\n\n```go\nvar name string = \"hello\"\nname := \"world\"\n```"},
				{"函数与包", "# 函数\n\n```go\nfunc add(a, b int) int {\n    return a + b\n}\n```"},
				{"结构体与接口", "# 结构体\n\n```go\ntype Person struct {\n    Name string\n    Age  int\n}\n```"},
				{"并发编程", "# Goroutine\n\n```go\ngo func() {\n    fmt.Println(\"Hello\")\n}()\n```"},
			},
		},
		{
			Title:       "Vue3实战开发",
			Description: "学习Vue3的组合式API，构建现代化的Web应用",
			Category:    "前端",
			AuthorID:    adminUser.ID,
			Chapters: []struct {
				Title   string
				Content string
			}{
				{"Vue3简介", "# Vue3新特性\n\n1. 组合式API\n2. 更好的TypeScript支持\n3. 更快的渲染速度"},
				{"创建项目", "# 创建Vue3项目\n\n```bash\nnpm create vue@latest\n```"},
				{"组合式API", "# setup函数\n\n组合式API是Vue3的核心特性..."},
				{"路由管理", "# Vue Router\n\n使用Vue Router管理页面路由..."},
			},
		},
		{
			Title:       "MySQL数据库优化",
			Description: "深入学习MySQL性能优化，掌握索引、查询优化等核心技能",
			Category:    "数据库",
			AuthorID:    adminUser.ID,
			Chapters: []struct {
				Title   string
				Content string
			}{
				{"MySQL架构", "# MySQL架构\n\n了解MySQL的整体架构..."},
				{"索引原理", "# 索引\n\nB+树索引原理..."},
				{"查询优化", "# 慢查询\n\n如何分析慢查询..."},
			},
		},
		{
			Title:       "Docker容器技术",
			Description: "学习Docker容器化技术，掌握容器编排与部署",
			Category:    "运维",
			AuthorID:    adminUser.ID,
			Chapters: []struct {
				Title   string
				Content string
			}{
				{"Docker基础", "# 什么是Docker\n\n容器化技术介绍..."},
				{"Dockerfile", "# 编写Dockerfile\n\n如何编写高效的Dockerfile..."},
				{"Docker Compose", "# 多容器编排\n\n使用Compose管理多个容器..."},
			},
		},
	}

	for _, c := range courses {
		var course model.Course
		result := util.Db.Where("title = ?", c.Title).First(&course)
		if result.Error == nil {
			continue
		}
		course = model.Course{
			Title:         c.Title,
			Description:   c.Description,
			Category:      c.Category,
			AuthorID:      c.AuthorID,
			ViewCount:     500 + int(time.Now().Unix()%500),
			LikeCount:     50 + int(time.Now().Unix()%50),
			FavoriteCount: 30 + int(time.Now().Unix()%30),
			ChapterCount:  len(c.Chapters),
			IsFree:        1,
			Priority:      10,
			Status:        1,
		}
		util.Db.Create(&course)

		for i, ch := range c.Chapters {
			chapter := model.CourseChapter{
				CourseID:  course.ID,
				Title:     ch.Title,
				Content:   ch.Content,
				SortOrder: i + 1,
			}
			util.Db.Create(&chapter)
		}
	}
	println("教程数据初始化完成！")
}

func initCheckins() {
	var users []model.User
	util.Db.Where("username IN ?", []string{"zhangsan", "lisi", "wangwu"}).Find(&users)

	for _, u := range users {
		for i := 0; i < 3; i++ {
			date := time.Now().AddDate(0, 0, -i)
			var existing model.Checkin
			result := util.Db.Where("user_id = ? AND checkin_date = ?", u.ID, date.Format("2006-01-02")).First(&existing)
			if result.Error == nil {
				continue
			}
			checkin := model.Checkin{
				UserID:      u.ID,
				CheckinDate: date,
				ExpGained:   10 + i,
				CoinsGained: 5 + i,
			}
			util.Db.Create(&checkin)
		}
	}
	println("打卡数据初始化完成！")
}
