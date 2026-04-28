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

	// 检查面试题数量
	var count int
	db.QueryRow("SELECT COUNT(*) FROM interview_questions").Scan(&count)
	fmt.Printf("面试题总数: %d\n", count)

	// 检查有deleted_at的记录
	var deletedCount int
	db.QueryRow("SELECT COUNT(*) FROM interview_questions WHERE deleted_at IS NOT NULL").Scan(&deletedCount)
	fmt.Printf("已删除: %d\n", deletedCount)

	// 检查表结构
	rows, _ := db.Query("DESCRIBE interview_questions")
	fmt.Println("\n表结构:")
	for rows.Next() {
		var field, typ, null, key string
		var def, extra sql.NullString
		rows.Scan(&field, &typ, &null, &key, &def, &extra)
		fmt.Printf("  %s: %s\n", field, typ)
	}
}
