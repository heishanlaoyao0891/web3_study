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

	fmt.Println("=== 初始化文章多分类关联 ===")
	initArticleCategories(db)

	fmt.Println("\n=== 初始化面试题多分类关联 ===")
	initInterviewQuestionCategories(db)

	fmt.Println("\n✅ 完成！")
}

func initArticleCategories(db *sql.DB) {
	// 获取分类ID映射
	catIDs := make(map[string]int)
	rows, _ := db.Query("SELECT id, name FROM categories WHERE deleted_at IS NULL")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		catIDs[name] = id
	}

	// 获取所有文章
	rows, _ = db.Query("SELECT id, title, category_id FROM articles WHERE deleted_at IS NULL")
	
	count := 0
	for rows.Next() {
		var id, categoryID int
		var title string
		rows.Scan(&id, &title, &categoryID)
		
		// 检查是否已存在关联
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM article_categories WHERE article_id = ? AND category_id = ?", id, categoryID).Scan(&exists)
		
		if exists == 0 {
			_, err := db.Exec("INSERT INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", id, categoryID)
			if err == nil {
				count++
				fmt.Printf("✓ 文章[%d] %s 添加分类关联\n", id, title)
			}
		}
	}
	fmt.Printf("共添加 %d 条文章分类关联\n", count)
}

func initInterviewQuestionCategories(db *sql.DB) {
	// 获取分类ID映射
	catIDs := make(map[string]int)
	rows, _ := db.Query("SELECT id, name FROM categories WHERE deleted_at IS NULL")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		catIDs[name] = id
	}

	// 获取所有面试题
	rows, _ = db.Query("SELECT id, title, category FROM interview_questions WHERE deleted_at IS NULL")
	
	count := 0
	for rows.Next() {
		var id int
		var title, category string
		rows.Scan(&id, &title, &category)
		
		// 根据category名称找到对应的分类ID
		catID, ok := catIDs[category]
		if !ok {
			fmt.Printf("⚠ 面试题[%d] %s 的分类 '%s' 不存在\n", id, title, category)
			continue
		}
		
		// 检查是否已存在关联
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM interview_question_categories WHERE interview_question_id = ? AND category_id = ?", id, catID).Scan(&exists)
		
		if exists == 0 {
			_, err := db.Exec("INSERT INTO interview_question_categories (interview_question_id, category_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", id, catID)
			if err == nil {
				count++
				fmt.Printf("✓ 面试题[%d] %s 添加分类关联\n", id, title)
			}
		}
	}
	fmt.Printf("共添加 %d 条面试题分类关联\n", count)
}
