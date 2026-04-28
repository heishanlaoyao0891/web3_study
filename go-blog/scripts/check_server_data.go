package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 直接查询
	fmt.Println("=== 直接SQL查询 ===")
	rows, err := db.Query("SELECT id, title, category FROM interview_questions ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var title, category string
		rows.Scan(&id, &title, &category)
		count++
		fmt.Printf("%d. [%s] %s\n", id, category, title)
	}
	fmt.Printf("\n共 %d 条\n", count)

	// 检查user_id是否有效
	fmt.Println("\n=== 检查user_id ===")
	rows2, _ := db.Query(`
		SELECT iq.id, iq.title, iq.user_id, u.id as user_exists
		FROM interview_questions iq
		LEFT JOIN users u ON iq.user_id = u.id
		WHERE u.id IS NULL
	`)
	defer rows2.Close()
	
	invalidCount := 0
	for rows2.Next() {
		var id, userID int
		var userExists sql.NullInt64
		var title string
		rows2.Scan(&id, &title, &userID, &userExists)
		invalidCount++
		fmt.Printf("无效user_id: ID=%d, user_id=%d, %s\n", id, userID, title)
	}
	if invalidCount == 0 {
		fmt.Println("所有user_id都有效")
	}
}
