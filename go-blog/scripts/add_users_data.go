package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("mysql", "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("=== 添加测试用户 ===")
	addTestUsers(db)

	fmt.Println("\n=== 添加Web3技术文章 ===")
	addWeb3Articles(db)

	fmt.Println("\n=== 添加面试题 ===")
	addInterviewQuestions(db)

	fmt.Println("\n✅ 完成！")
}

func addTestUsers(db *sql.DB) {
	users := []struct {
		Username, Password, Nickname, Bio string
		Level, Exp, Coins int
	}{
		{"solidity_dev", "123456", "Solidity开发者", "3年智能合约开发经验，专注DeFi协议开发", 4, 600, 400},
		{"defi_researcher", "123456", "DeFi研究员", "研究DeFi协议机制，AMM和借贷协议分析", 3, 400, 250},
		{"nft_creator", "123456", "NFT创作者", "NFT开发与设计，探索数字艺术与区块链结合", 2, 200, 150},
	}

	for _, u := range users {
		var id int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", u.Username).Scan(&id)
		if err == sql.ErrNoRows {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			_, err := db.Exec(`INSERT INTO users (username, password, nickname, bio, level, exp, coins, status, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, 1, NOW(), NOW())`,
				u.Username, string(hashedPassword), u.Nickname, u.Bio, u.Level, u.Exp, u.Coins)
			if err != nil {
				fmt.Printf("创建用户 %s 失败: %v\n", u.Username, err)
			} else {
				fmt.Printf("✓ 创建用户: %s / 123456\n", u.Username)
			}
		} else {
			fmt.Printf("用户 %s 已存在\n", u.Username)
		}
	}
}

func addWeb3Articles(db *sql.DB) {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	var catIDs = make(map[string]int)
	rows, _ := db.Query("SELECT id, name FROM categories")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		catIDs[name] = id
	}

	articles := []struct {
		Title, Content, Category string
		Views int
	}{
		{"EIP-4844 Blob交易详解", "# EIP-4844 Blob交易详解\n\n## 背景\nEIP-4844引入Blob数据类型，为Layer2降低Gas费用90%。\n\n## 核心机制\n\n### Blob特点\n1. 临时存储，18天后删除\n2. 独立Gas计费\n3. 每区块目标3个Blob\n\n### 代码示例\n```solidity\ntype BlobTransaction struct {\n    BlobFeeCap      *big.Int\n    BlobHashes      []common.Hash\n}\n```", "Layer2", 2856},
		{"Foundry开发框架指南", "# Foundry开发框架指南\n\n## 优势\n- 编译速度快10-100倍\n- 支持Solidity写测试\n- 内置模糊测试\n\n## 安装\n```bash\ncurl -L https://foundry.paradigm.xyz | bash\nfoundryup\n```\n\n## 测试示例\n```solidity\nfunction testFuzz_SetValue(uint256 value) public {\n    store.setValue(value);\n    assertEq(store.getValue(), value);\n}\n```", "智能合约", 3421},
		{"Aave V3借贷原理", "# Aave V3借贷原理\n\n## 核心概念\n\n### 健康因子\n```\nHealthFactor = (抵押品 × 清算阈值) / 总债务\n```\n\n### 清算机制\nHF < 1时可被清算，清算人获得5-10%奖金。\n\n## V3新特性\n1. 隔离模式\n2. 效率模式\n3. 跨链门户", "DeFi协议", 3156},
		{"ERC-4337账户抽象", "# ERC-4337账户抽象\n\n## 解决问题\n1. 私钥丢失恢复\n2. 无Gas交易\n3. 批量操作\n4. 交易限制\n\n## 核心组件\n- UserOperation\n- EntryPoint\n- Paymaster\n- Bundler", "智能合约", 2678},
		{"跨链桥安全实践", "# 跨链桥安全实践\n\n## 常见攻击\n1. 签名验证漏洞\n2. 多签配置不当\n3. 预言机操控\n\n## 防护措施\n1. 严格签名验证\n2. 限制转账金额\n3. 时间锁延迟\n4. 紧急暂停\n\n## 审计清单\n- 签名验证\n- 重放攻击防护\n- 多签配置", "安全审计", 2345},
	}

	for _, a := range articles {
		var existingID int
		err := db.QueryRow("SELECT id FROM articles WHERE title = ?", a.Title).Scan(&existingID)
		if err == sql.ErrNoRows {
			_, err := db.Exec(`INSERT INTO articles (title, content, user_id, category_id, view_count, like_count, status, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, 1, NOW(), NOW())`,
				a.Title, a.Content, adminID, catIDs[a.Category], a.Views, a.Views/15)
			if err != nil {
				fmt.Printf("创建文章失败 %s: %v\n", a.Title, err)
			} else {
				fmt.Printf("✓ 添加文章: %s\n", a.Title)
			}
		} else {
			fmt.Printf("文章已存在: %s\n", a.Title)
		}
	}
}

func addInterviewQuestions(db *sql.DB) {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	questions := []struct {
		Title, Content, Answer, Category string
		Difficulty int
	}{
		{"解释EIP-4844的作用", "EIP-4844核心机制是什么？", "引入Blob数据类型，为L2降低Gas费用80-90%。Blob临时存储18天，独立Gas计费。", "Layer2", 2},
		{"Foundry vs Hardhat", "两个框架如何选择？", "Foundry: 编译快，Solidity测试，内置模糊测试\nHardhat: JS生态，插件丰富\n纯Solidity选Foundry，需要JS集成选Hardhat", "智能合约", 1},
		{"Aave健康因子计算", "如何计算健康因子？", "HF = (抵押品 × 清算阈值) / 总债务\nHF > 1安全，HF < 1可被清算", "DeFi协议", 2},
		{"账户抽象解决了什么？", "ERC-4337的意义？", "1. 私钥丢失恢复\n2. 无Gas交易\n3. 批量操作\n4. 交易限制\n5. 多签管理", "智能合约", 2},
		{"跨链桥安全风险", "跨链桥主要风险？", "1. 签名验证漏洞\n2. 多签配置不当\n3. 预言机操控\n4. 智能合约漏洞\n5. 私钥泄露", "安全审计", 3},
		{"代理模式选择", "透明代理 vs UUPS？", "透明代理：兼容好，每次检查调用者\nUUPS：Gas省，升级逻辑在实现合约\nBeacon：多个代理共享升级", "智能合约", 2},
		{"MEV是什么？", "解释MEV及其影响", "矿工/验证者可提取价值。包括：三明治攻击、套利、清算抢跑。\n\n防护：Flashbots、私有交易池、时间分散交易", "DeFi协议", 2},
		{"zkSync开发注意点", "zkSync开发有什么特殊之处？", "1. 不支持某些Solidity特性\n2. 账户抽象原生支持\n3. 使用专有SDK部署\n4. 编译器不同\n5. Gas计算不同", "Layer2", 3},
	}

	for _, q := range questions {
		var existingID int
		err := db.QueryRow("SELECT id FROM interview_questions WHERE title = ?", q.Title).Scan(&existingID)
		if err == sql.ErrNoRows {
			_, err := db.Exec(`INSERT INTO interview_questions (title, content, answer, category, difficulty, user_id, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
				q.Title, q.Content, q.Answer, q.Category, q.Difficulty, adminID)
			if err != nil {
				fmt.Printf("创建面试题失败 %s: %v\n", q.Title, err)
			} else {
				fmt.Printf("✓ 添加面试题: %s\n", q.Title)
			}
		} else {
			fmt.Printf("面试题已存在: %s\n", q.Title)
		}
	}
}
