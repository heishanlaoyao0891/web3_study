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

	// 检查面试题和用户关联
	fmt.Println("=== 面试题和用户关联 ===")
	rows, _ := db.Query(`
		SELECT iq.id, iq.title, iq.user_id, u.username
		FROM interview_questions iq
		JOIN users u ON iq.user_id = u.id
		ORDER BY iq.id DESC
	`)
	for rows.Next() {
		var id, userID int
		var title, username string
		rows.Scan(&id, &title, &userID, &username)
		fmt.Printf("ID:%d user_id:%d (%s) %s\n", id, userID, username, title)
	}

	// 检查ID=18的记录详情
	fmt.Println("\n=== ID=18详情 ===")
	var title, content, answer, category string
	var viewCount int
	db.QueryRow("SELECT title, content, answer, category, view_count FROM interview_questions WHERE id = 18").Scan(&title, &content, &answer, &category, &viewCount)
	fmt.Printf("Title: %s\nCategory: %s\nViews: %d\n", title, category, viewCount)
}
