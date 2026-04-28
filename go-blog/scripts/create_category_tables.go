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

	// 创建文章分类关联表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS article_categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			article_id INT NOT NULL,
			category_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_article_category (article_id, category_id),
			INDEX idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	if err != nil {
		fmt.Println("创建article_categories表失败:", err)
	} else {
		fmt.Println("✓ 创建article_categories表成功")
	}

	// 创建面试题分类关联表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interview_question_categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			interview_question_id INT NOT NULL,
			category_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_question_category (interview_question_id, category_id),
			INDEX idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	if err != nil {
		fmt.Println("创建interview_question_categories表失败:", err)
	} else {
		fmt.Println("✓ 创建interview_question_categories表成功")
	}
}
