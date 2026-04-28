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

	fmt.Println("=== 面试题分类统计 ===")
	rows, _ := db.Query("SELECT category, COUNT(*) as cnt FROM interview_questions WHERE deleted_at IS NULL GROUP BY category")
	for rows.Next() {
		var category string
		var count int
		rows.Scan(&category, &count)
		fmt.Printf("分类: %s, 数量: %d\n", category, count)
	}

	fmt.Println("\n=== 所有面试题 ===")
	rows, _ = db.Query("SELECT id, title, category FROM interview_questions WHERE deleted_at IS NULL ORDER BY id")
	for rows.Next() {
		var id int
		var title, category string
		rows.Scan(&id, &title, &category)
		fmt.Printf("ID:%d [%s] %s\n", id, category, title)
	}

	fmt.Println("\n=== categories表 ===")
	rows, _ = db.Query("SELECT id, name FROM categories WHERE deleted_at IS NULL")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("ID:%d %s\n", id, name)
	}
}
