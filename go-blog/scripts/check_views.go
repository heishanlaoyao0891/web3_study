package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	defer db.Close()

	rows, _ := db.Query("SELECT id, title, view_count FROM interview_questions ORDER BY id DESC")
	fmt.Println("面试题浏览量:")
	for rows.Next() {
		var id, views int
		var title string
		rows.Scan(&id, &title, &views)
		fmt.Printf("ID:%d views:%d - %s\n", id, views, title)
	}
}
