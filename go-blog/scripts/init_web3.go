package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	dsn := "root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("=== 清理旧数据 ===")
	cleanOldData()

	fmt.Println("\n=== 初始化Web3数据 ===")
	initAdminUser()
	initCategories()
	initTags()
	initArticles()
	initQuestions()
	initInterviewQuestions()
	initCodeSnippets()
	initContractTemplates()
	initResources()
	initLearningPaths()
	initCourses()

	fmt.Println("\n✅ Web3数据初始化完成！")
}

func cleanOldData() {
	tables := []string{
		"article_tags", "likes", "favorites", "answers", "questions",
		"course_chapters", "learning_records", "courses", "checkins",
		"comments", "articles", "interview_questions", "code_snippets",
		"contract_templates", "resources", "learning_chapters", "learning_paths",
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}
	db.Exec("DELETE FROM users WHERE username != 'admin'")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM categories")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	
	fmt.Println("✓ 旧数据已清理")
}

func initAdminUser() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&id)
	if err == sql.ErrNoRows {
		_, err := db.Exec(`INSERT INTO users (username, password, nickname, bio, level, exp, coins, status, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
			"admin", string(hashedPassword), "Web3技术专家", "专注Web3开发，智能合约安全审计，DeFi协议研究", 5, 800, 500, 1)
		if err != nil {
			log.Println("创建admin失败:", err)
		}
	} else {
		db.Exec(`UPDATE users SET nickname=?, bio=?, level=?, exp=?, coins=? WHERE username=?`,
			"Web3技术专家", "专注Web3开发，智能合约安全审计，DeFi协议研究", 5, 800, 500, "admin")
	}
	fmt.Println("✓ 管理员账号: admin / 123456")
}

func initCategories() {
	categories := []struct {
		Name, Desc string
	}{
		{"智能合约", "Solidity开发、EVM原理、Gas优化、合约模式"},
		{"DeFi协议", "Uniswap、Aave、Compound、AMM机制、借贷协议"},
		{"NFT开发", "ERC-721、ERC-1155、元数据、版税、白名单机制"},
		{"Layer2", "Optimism、Arbitrum、zkSync、StarkNet、跨链桥"},
		{"安全审计", "重入攻击、闪电贷、价格操控、形式化验证"},
		{"公链研究", "以太坊、Solana、Polkadot、Cosmos生态"},
		{"面试求职", "Web3面试题、智能合约工程师、薪资行情"},
	}

	for _, c := range categories {
		db.Exec("INSERT INTO categories (name, `desc`, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", c.Name, c.Desc)
	}
	fmt.Println("✓ 分类数据")
}

func initTags() {
	tags := []string{
		"Solidity", "EVM", "Gas优化", "智能合约安全", "OpenZeppelin",
		"Uniswap", "Aave", "Compound", "AMM", "DeFi",
		"ERC-20", "ERC-721", "ERC-1155", "NFT", "版税",
		"Optimism", "Arbitrum", "zkSync", "StarkNet", "Layer2",
		"跨链", "预言机", "Chainlink", "DAO", "治理",
		"重入攻击", "闪电贷", "MEV", "空投", "MerkleTree",
		"Hardhat", "Foundry", "Truffle", "Web3.js", "Ethers.js",
		"以太坊", "Solana", "Rust", "Move", "Substrate",
	}

	for _, name := range tags {
		db.Exec("INSERT INTO tags (name, created_at, updated_at) VALUES (?, NOW(), NOW())", name)
	}
	fmt.Println("✓ 标签数据")
}

func initArticles() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	var catIDs map[string]int
	catIDs = make(map[string]int)
	rows, _ := db.Query("SELECT id, name FROM categories")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		catIDs[name] = id
	}

	articles := []struct {
		Title, Content string
		Category string
		Views int
	}{
		{"Solidity智能合约开发完整指南", "# Solidity智能合约开发完整指南\n\n## 一、环境搭建\n\n### 1. 安装Node.js\n推荐使用Node.js v18+版本。\n\n### 2. 安装Hardhat\n```bash\nnpm install --save-dev hardhat\nnpx hardhat init\n```\n\n## 二、第一个智能合约\n\n```solidity\npragma solidity ^0.8.19;\n\ncontract SimpleStorage {\n    uint256 private _value;\n    \n    function setValue(uint256 value) public {\n        _value = value;\n    }\n    \n    function getValue() public view returns (uint256) {\n        return _value;\n    }\n}\n```\n\n## 三、安全最佳实践\n\n使用ReentrancyGuard防止重入攻击。", "智能合约", 3521},
		{"Uniswap V3集中流动性原理详解", "# Uniswap V3集中流动性原理详解\n\n## 一、V2 vs V3对比\n\nUniswap V3允许LP在指定价格区间提供流动性，资本效率提升4000倍。\n\n## 二、核心机制\n\n价格空间被划分为离散的tick：price = 1.0001^tick", "DeFi协议", 4215},
		{"NFT智能合约开发实战", "# NFT智能合约开发实战\n\n## 一、ERC-721标准\n\n```solidity\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract MyNFT is ERC721 {\n    constructor() ERC721(\"MyNFT\", \"MNFT\") {}\n}\n```", "NFT开发", 2856},
		{"智能合约安全审计要点", "# 智能合约安全审计要点\n\n## 常见漏洞\n\n1. 重入攻击 - 使用Checks-Effects-Interactions模式\n2. 闪电贷价格操控 - 使用TWAP\n3. 访问控制漏洞 - 使用initializer", "安全审计", 2156},
		{"Layer2技术全景解析", "# Layer2技术全景解析\n\n## 技术路线\n\n### Optimistic Rollup\nOptimism、Arbitrum - 7天挑战期\n\n### ZK Rollup\nzkSync、StarkNet - 即时最终性", "Layer2", 1892},
	}

	for _, a := range articles {
		db.Exec(`INSERT INTO articles (title, content, user_id, category_id, view_count, like_count, status, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, 1, NOW(), NOW())`,
			a.Title, a.Content, adminID, catIDs[a.Category], a.Views, a.Views/15)
	}
	fmt.Println("✓ 文章数据")
}

func initQuestions() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	questions := []struct {
		Title, Content string
	}{
		{"Solidity中如何实现可升级合约？", "我想实现一个可以升级的智能合约，应该使用哪种方案？"},
		{"Uniswap V3如何计算无常损失？", "在Uniswap V3提供流动性时，无常损失的计算公式是什么？"},
		{"NFT白名单使用MerkleTree还是签名验证？", "两种方案各有什么优缺点？"},
		{"zkSync和StarkNet开发有什么区别？", "这两个ZK Rollup在开发体验、性能方面有什么差异？"},
		{"如何防止闪电贷攻击？", "DeFi协议如何防止闪电贷价格操控攻击？"},
	}

	for _, q := range questions {
		db.Exec(`INSERT INTO questions (title, content, user_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())`,
			q.Title, q.Content, adminID)
	}
	fmt.Println("✓ 问答数据")
}

func initInterviewQuestions() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	questions := []struct {
		Title, Content, Answer, Category string
		Difficulty int
	}{
		{"什么是重入攻击？如何防御？", "解释重入攻击原理和防御方法", "Checks-Effects-Interactions模式，ReentrancyGuard修饰符", "安全审计", 2},
		{"解释Uniswap V2的AMM机制", "Uniswap V2如何实现自动做市？", "恒定乘积公式 x*y=k，LP收取0.3%手续费", "DeFi协议", 2},
		{"ERC-721和ERC-1155的区别？", "比较两种NFT标准", "ERC-721每个Token唯一，ERC-1155支持批量操作", "NFT开发", 1},
		{"storage、memory、calldata区别？", "解释三种数据存储位置", "storage永久存储，memory临时存储，calldata只读参数", "智能合约", 1},
		{"什么是闪电贷？", "解释闪电贷原理和应用场景", "在一个交易内完成借款、使用、还款，无需抵押", "DeFi协议", 2},
	}

	for _, q := range questions {
		db.Exec(`INSERT INTO interview_questions (title, content, answer, category, difficulty, user_id, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
			q.Title, q.Content, q.Answer, q.Category, q.Difficulty, adminID)
	}
	fmt.Println("✓ 面试题数据")
}

func initCodeSnippets() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	snippets := []struct {
		Title, Desc, Code, Category string
	}{
		{"ERC20代币合约模板", "标准ERC20代币合约", "pragma solidity ^0.8.19;\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract MyToken is ERC20 {\n    constructor(uint256 supply) ERC20(\"MyToken\", \"MTK\") {\n        _mint(msg.sender, supply);\n    }\n}", "代币"},
		{"ReentrancyGuard使用", "防止重入攻击", "import \"@openzeppelin/contracts/security/ReentrancyGuard.sol\";\n\ncontract SafeContract is ReentrancyGuard {\n    function withdraw() public nonReentrant {}\n}", "安全"},
		{"Ownable权限控制", "简单的所有者权限", "contract Ownable {\n    address public owner;\n    modifier onlyOwner() { require(msg.sender == owner); _; }\n}", "权限"},
	}

	for _, s := range snippets {
		db.Exec(`INSERT INTO code_snippets (title, description, code, language, category, user_id, created_at, updated_at) 
			VALUES (?, ?, ?, 'solidity', ?, ?, NOW(), NOW())`,
			s.Title, s.Desc, s.Code, s.Category, adminID)
	}
	fmt.Println("✓ 代码片段数据")
}

func initContractTemplates() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	templates := []struct {
		Name, Desc, Code, Category string
	}{
		{"ERC20代币模板", "标准ERC20代币合约", "pragma solidity ^0.8.19;\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract TokenTemplate is ERC20 {\n    constructor() ERC20(\"Token\", \"TKN\") {\n        _mint(msg.sender, 1000000 * 10 ** decimals());\n    }\n}", "代币"},
		{"ERC721 NFT模板", "标准NFT合约", "pragma solidity ^0.8.19;\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract NFTTemplate is ERC721 {\n    constructor() ERC721(\"NFT\", \"NFT\") {}\n}", "NFT"},
	}

	for _, t := range templates {
		db.Exec(`INSERT INTO contract_templates (name, description, code, category, user_id, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
			t.Name, t.Desc, t.Code, t.Category, adminID)
	}
	fmt.Println("✓ 合约模板数据")
}

func initResources() {
	resources := []struct {
		Title, URL, Desc, Category string
	}{
		{"Solidity官方文档", "https://docs.soliditylang.org/", "Solidity编程语言官方文档", "文档"},
		{"OpenZeppelin文档", "https://docs.openzeppelin.com/", "智能合约安全库文档", "文档"},
		{"Hardhat文档", "https://hardhat.org/docs", "以太坊开发环境", "工具"},
		{"Foundry Book", "https://book.getfoundry.sh/", "快速Solidity开发工具", "工具"},
		{"Ethereum官方文档", "https://ethereum.org/developers", "以太坊开发者资源", "文档"},
		{"Web3.js文档", "https://web3js.readthedocs.io/", "JavaScript以太坊库", "文档"},
		{"Ethers.js文档", "https://docs.ethers.org/", "轻量级以太坊库", "文档"},
		{"Remix IDE", "https://remix.ethereum.org/", "在线Solidity IDE", "工具"},
	}

	for _, r := range resources {
		db.Exec(`INSERT INTO resources (title, url, description, category, created_at, updated_at) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`,
			r.Title, r.URL, r.Desc, r.Category)
	}
	fmt.Println("✓ 资源导航数据")
}

func initLearningPaths() {
	paths := []struct {
		Title, Desc string
		Difficulty int
		Duration string
	}{
		{"Web3开发入门路径", "从零开始学习区块链开发和智能合约编程", 1, "2个月"},
		{"DeFi协议开发路径", "深入学习DeFi协议原理和开发实战", 2, "3个月"},
		{"智能合约安全审计路径", "成为专业的智能合约安全审计师", 3, "4个月"},
	}

	for _, p := range paths {
		db.Exec(`INSERT INTO learning_paths (title, description, difficulty, duration, created_at, updated_at) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`,
			p.Title, p.Desc, p.Difficulty, p.Duration)
	}
	fmt.Println("✓ 学习路径数据")
}

func initCourses() {
	var adminID int
	db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	courses := []struct {
		Title, Desc, Category string
	}{
		{"Solidity从入门到精通", "系统学习Solidity智能合约开发", "智能合约"},
		{"DeFi协议开发实战", "学习开发AMM、借贷协议等DeFi应用", "DeFi"},
		{"NFT开发完整指南", "从零开发NFT合约，包含白名单、版税等功能", "NFT"},
		{"智能合约安全审计", "学习智能合约安全漏洞与审计技术", "安全"},
	}

	for _, c := range courses {
		db.Exec(`INSERT INTO courses (title, description, category, author_id, created_at, updated_at) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`,
			c.Title, c.Desc, c.Category, adminID)
	}
	fmt.Println("✓ 教程数据")
}
