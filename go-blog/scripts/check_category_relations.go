package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	defer db.Close()

	// 检查文章分类关联
	fmt.Println("=== 文章多分类关联 ===")
	rows, _ := db.Query(`
		SELECT a.id, a.title, GROUP_CONCAT(c.name) as categories
		FROM articles a
		LEFT JOIN article_categories ac ON a.id = ac.article_id
		LEFT JOIN categories c ON ac.category_id = c.id
		WHERE a.deleted_at IS NULL
		GROUP BY a.id
		ORDER BY a.id DESC
		LIMIT 5
	`)
	for rows.Next() {
		var id int
		var title, categories string
		rows.Scan(&id, &title, &categories)
		fmt.Printf("文章[%d] %s\n  分类: %s\n", id, title, categories)
	}

	// 检查面试题分类关联
	fmt.Println("\n=== 面试题多分类关联 ===")
	rows, _ = db.Query(`
		SELECT iq.id, iq.title, GROUP_CONCAT(c.name) as categories
		FROM interview_questions iq
		LEFT JOIN interview_question_categories iqc ON iq.id = iqc.interview_question_id
		LEFT JOIN categories c ON iqc.category_id = c.id
		WHERE iq.deleted_at IS NULL
		GROUP BY iq.id
		ORDER BY iq.id DESC
		LIMIT 5
	`)
	for rows.Next() {
		var id int
		var title, categories string
		rows.Scan(&id, &title, &categories)
		fmt.Printf("面试题[%d] %s\n  分类: %s\n", id, title, categories)
	}
}
