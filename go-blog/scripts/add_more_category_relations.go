package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	defer db.Close()

	// 获取分类ID
	catIDs := make(map[string]int)
	rows, _ := db.Query("SELECT id, name FROM categories WHERE deleted_at IS NULL")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		catIDs[name] = id
	}

	// 为某些文章添加额外的分类
	// Solidity开发指南 -> 同时属于"安全审计"
	db.Exec("INSERT IGNORE INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (29, ?, NOW(), NOW())", catIDs["安全审计"])
	fmt.Println("✓ 文章[29] 添加 '安全审计' 分类")

	// Uniswap V3 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (30, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 文章[30] 添加 '智能合约' 分类")

	// NFT开发实战 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (31, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 文章[31] 添加 '智能合约' 分类")

	// Layer2技术解析 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (33, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 文章[33] 添加 '智能合约' 分类")

	// 跨链桥安全实践 -> 同时属于"Layer2"
	db.Exec("INSERT IGNORE INTO article_categories (article_id, category_id, created_at, updated_at) VALUES (39, ?, NOW(), NOW())", catIDs["Layer2"])
	fmt.Println("✓ 文章[39] 添加 'Layer2' 分类")

	// 面试题多分类
	// 重入攻击 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO interview_question_categories (interview_question_id, category_id, created_at, updated_at) VALUES (6, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 面试题[6] 添加 '智能合约' 分类")

	// Uniswap AMM -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO interview_question_categories (interview_question_id, category_id, created_at, updated_at) VALUES (7, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 面试题[7] 添加 '智能合约' 分类")

	// Aave健康因子 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO interview_question_categories (interview_question_id, category_id, created_at, updated_at) VALUES (13, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 面试题[13] 添加 '智能合约' 分类")

	// zkSync开发 -> 同时属于"智能合约"
	db.Exec("INSERT IGNORE INTO interview_question_categories (interview_question_id, category_id, created_at, updated_at) VALUES (18, ?, NOW(), NOW())", catIDs["智能合约"])
	fmt.Println("✓ 面试题[18] 添加 '智能合约' 分类")

	fmt.Println("\n✅ 完成！")
}
