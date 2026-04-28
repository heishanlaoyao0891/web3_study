package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	counts := map[string]string{
		"users":              "SELECT COUNT(*) FROM users",
		"categories":         "SELECT COUNT(*) FROM categories",
		"tags":               "SELECT COUNT(*) FROM tags",
		"articles":           "SELECT COUNT(*) FROM articles",
		"questions":          "SELECT COUNT(*) FROM questions",
		"interview_questions": "SELECT COUNT(*) FROM interview_questions",
		"code_snippets":      "SELECT COUNT(*) FROM code_snippets",
		"contract_templates": "SELECT COUNT(*) FROM contract_templates",
		"resources":          "SELECT COUNT(*) FROM resources",
		"learning_paths":     "SELECT COUNT(*) FROM learning_paths",
		"courses":            "SELECT COUNT(*) FROM courses",
	}

	fmt.Println("=== 数据统计 ===")
	for name, query := range counts {
		var count int
		db.QueryRow(query).Scan(&count)
		fmt.Printf("%-20s: %d\n", name, count)
	}

	fmt.Println("\n=== 用户列表 ===")
	rows, _ := db.Query("SELECT id, username, nickname, level, exp, coins FROM users")
	for rows.Next() {
		var id, level, exp, coins int
		var username, nickname string
		rows.Scan(&id, &username, &nickname, &level, &exp, &coins)
		fmt.Printf("ID:%d 用户:%s 昵称:%s Lv:%d 经验:%d 金币:%d\n", id, username, nickname, level, exp, coins)
	}

	fmt.Println("\n=== 分类列表 ===")
	rows, _ = db.Query("SELECT id, name FROM categories")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("ID:%d %s\n", id, name)
	}

	fmt.Println("\n=== 文章列表 ===")
	rows, _ = db.Query("SELECT id, title, view_count FROM articles")
	for rows.Next() {
		var id, views int
		var title string
		rows.Scan(&id, &title, &views)
		fmt.Printf("ID:%d %s (浏览:%d)\n", id, title, views)
	}
}
