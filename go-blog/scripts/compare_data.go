package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	defer db.Close()

	// 比较ID=18和其他记录
	fmt.Println("=== ID=18 vs 其他记录 ===")
	
	var id18Content, id17Content string
	db.QueryRow("SELECT content FROM interview_questions WHERE id = 18").Scan(&id18Content)
	db.QueryRow("SELECT content FROM interview_questions WHERE id = 17").Scan(&id17Content)
	
	fmt.Printf("ID=18 content length: %d\n", len(id18Content))
	fmt.Printf("ID=17 content length: %d\n", len(id17Content))
	
	// 检查deleted_at
	fmt.Println("\n=== 检查deleted_at ===")
	rows, _ := db.Query("SELECT id, deleted_at FROM interview_questions ORDER BY id DESC LIMIT 5")
	for rows.Next() {
		var id int
		var deletedAt sql.NullTime
		rows.Scan(&id, &deletedAt)
		if deletedAt.Valid {
			fmt.Printf("ID:%d deleted_at: %v\n", id, deletedAt.Time)
		} else {
			fmt.Printf("ID:%d deleted_at: NULL\n", id)
		}
	}
	
	// 检查是否有重复ID
	fmt.Println("\n=== 检查重复记录 ===")
	rows, _ = db.Query("SELECT id, COUNT(*) as cnt FROM interview_questions GROUP BY id HAVING cnt > 1")
	count := 0
	for rows.Next() {
		var id, cnt int
		rows.Scan(&id, &cnt)
		fmt.Printf("重复ID: %d, 次数: %d\n", id, cnt)
		count++
	}
	if count == 0 {
		fmt.Println("没有重复ID")
	}
}
