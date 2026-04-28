package main

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	godotenv.Load()
}

func main() {
	if err := util.CreateTestDB(".env_local"); err != nil {
		panic(err)
	}

	fmt.Println("开始初始化Web3博客数据...")

	initUsers()
	initCategories()
	initTags()
	initArticles()
	initQuestions()
	initCourses()
	initLearningPaths()
	initCodeSnippets()
	initContractTemplates()
	initResources()
	initInterviewQuestions()

	fmt.Println("Web3数据初始化完成！")
}

func initUsers() {
	users := []struct {
		Username string
		Password string
		Nickname string
		Bio      string
		Level    int
		Exp      int
		Coins    int
	}{
		{"admin", "123456", "Web3管理员", "Web3技术社区管理员，专注区块链技术开发与推广", 5, 800, 500},
		{"vitalik_fan", "123456", "V神粉丝", "以太坊爱好者，研究智能合约安全", 4, 600, 350},
		{"solidity_dev", "123456", "Solidity开发者", "3年智能合约开发经验，专注DeFi协议", 5, 900, 600},
		{"defi_master", "123456", "DeFi大师", "DeFi协议研究员，AMM机制深度分析", 4, 700, 400},
		{"nft_artist", "123456", "NFT艺术家", "NFT创作者与开发者，探索数字艺术", 3, 400, 250},
		{"layer2_dev", "123456", "L2开发者", "专注Layer2扩容方案，Optimism/Arbitrum研究", 4, 650, 380},
		{"security_auditor", "123456", "安全审计师", "智能合约安全审计，漏洞挖掘专家", 6, 1200, 800},
		{"blockchain_newbie", "123456", "区块链新手", "正在学习Web3开发，希望成为智能合约工程师", 1, 50, 30},
	}

	for _, u := range users {
		var existing model.User
		result := util.Db.Where("username = ?", u.Username).First(&existing)
		if result.Error == nil {
			continue
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		user := model.User{
			Username: u.Username,
			Password: string(hashedPassword),
			Nickname: u.Nickname,
			Bio:      u.Bio,
			Level:    u.Level,
			Exp:      u.Exp,
			Coins:    u.Coins,
			Status:   1,
		}
		util.Db.Create(&user)
	}
	fmt.Println("用户数据初始化完成！")
}

func initCategories() {
	categories := []struct {
		Name string
		Desc string
	}{
		{"智能合约", "Solidity智能合约开发、安全审计、Gas优化"},
		{"DeFi", "去中心化金融协议、AMM、借贷、衍生品"},
		{"NFT", "NFT铸造、交易、元数据标准、版税机制"},
		{"Layer2", "Rollup、状态通道、侧链等扩容方案"},
		{"公链研究", "以太坊、Solana、Polkadot等公链技术"},
		{"安全审计", "智能合约漏洞、安全最佳实践、审计工具"},
		{"求职面试", "Web3岗位面试题、薪资行情、职业发展"},
	}

	for _, c := range categories {
		var existing model.Category
		result := util.Db.Where("name = ?", c.Name).First(&existing)
		if result.Error == nil {
			continue
		}

		category := model.Category{
			Name: c.Name,
			Desc: c.Desc,
		}
		util.Db.Create(&category)
	}
	fmt.Println("分类数据初始化完成！")
}

func initTags() {
	tags := []string{
		"Solidity", "EVM", "Gas优化", "智能合约安全", "OpenZeppelin",
		"Uniswap", "Aave", "Compound", "Curve", "DeFi",
		"ERC-20", "ERC-721", "ERC-1155", "NFT", "元数据",
		"Optimism", "Arbitrum", "zkSync", "StarkNet", "Layer2",
		"跨链", "预言机", "Chainlink", "DAO", "治理",
		"重入攻击", "闪电贷", "MEV", "空投", "白名单",
		"Hardhat", "Foundry", "Truffle", "Web3.js", "Ethers.js",
		"以太坊", "Solana", "Polkadot", "Cosmos", "Avalanche",
	}

	for _, name := range tags {
		var existing model.Tag
		result := util.Db.Where("name = ?", name).First(&existing)
		if result.Error == nil {
			continue
		}

		tag := model.Tag{Name: name}
		util.Db.Create(&tag)
	}
	fmt.Println("标签数据初始化完成！")
}

func initArticles() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	var catSmartContract, catDeFi, catNFT, catLayer2, catSecurity model.Category
	util.Db.Where("name = ?", "智能合约").First(&catSmartContract)
	util.Db.Where("name = ?", "DeFi").First(&catDeFi)
	util.Db.Where("name = ?", "NFT").First(&catNFT)
	util.Db.Where("name = ?", "Layer2").First(&catLayer2)
	util.Db.Where("name = ?", "安全审计").First(&catSecurity)

	articles := []struct {
		Title      string
		Content    string
		UserID     uint
		CategoryID uint
		ViewCount  int
		LikeCount  int
	}{
		{
			"Solidity智能合约开发入门指南",
			"# Solidity智能合约开发入门指南\n\n## 一、环境搭建\n\n### 安装Node.js\n推荐使用Node.js v18+\n\n### 安装Hardhat\n```bash\nnpm install --save-dev hardhat\nnpx hardhat init\n```\n\n## 二、第一个合约\n\n```solidity\n// SPDX-License-Identifier: MIT\npragma solidity ^0.8.19;\n\ncontract SimpleStorage {\n    uint256 private storedData;\n    \n    function set(uint256 x) public {\n        storedData = x;\n    }\n    \n    function get() public view returns (uint256) {\n        return storedData;\n    }\n}\n```\n\n## 三、编译部署\n\n```bash\nnpx hardhat compile\nnpx hardhat run scripts/deploy.js\n```\n\n## 四、测试\n\n```javascript\ndescribe(\"SimpleStorage\", function () {\n  it(\"Should store value\", async function () {\n    const Storage = await ethers.getContractFactory(\"SimpleStorage\");\n    const storage = await Storage.deploy();\n    await storage.set(42);\n    expect(await storage.get()).to.equal(42);\n  });\n});\n```",
			admin.ID, catSmartContract.ID, 2856, 186,
		},
		{
			"Uniswap V3 AMM机制解析",
			"# Uniswap V3 AMM机制解析\n\n## 核心概念\n\n### 集中流动性\nV3允许LP在指定价格区间提供流动性，资本效率提升4000倍。\n\n### Tick机制\n价格空间划分为离散的tick：\n```\nprice = 1.0001 ** tick\n```\n\n## 流动性计算\n\n当价格在区间内时：\n```\nL = sqrt(x * y)\n```\n\n## 实际应用\n\n1. 稳定币交易：在0.99-1.01区间集中流动性\n2. 限价单：在单点提供流动性\n3. 套利策略：跨DEX价格差套利",
			admin.ID, catDeFi.ID, 3421, 256,
		},
		{
			"NFT智能合约开发实战",
			"# NFT智能合约开发实战\n\n## ERC-721标准\n\n```solidity\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract MyNFT is ERC721 {\n    constructor() ERC721(\"MyNFT\", \"MNFT\") {}\n    \n    function safeMint(address to, uint256 tokenId) public {\n        _safeMint(to, tokenId);\n    }\n}\n```\n\n## 元数据标准\n\n支持链上元数据和IPFS存储。\n\n## 高级功能\n\n- 可升级NFT\n- 批量铸造\n- 版税机制\n- 白名单预售",
			admin.ID, catNFT.ID, 2156, 178,
		},
		{
			"智能合约安全审计要点",
			"# 智能合约安全审计要点\n\n## 常见漏洞\n\n### 1. 重入攻击\n使用Checks-Effects-Interactions模式或ReentrancyGuard。\n\n### 2. 整数溢出\nSolidity 0.8+自动检查，但unchecked块中仍需注意。\n\n### 3. 访问控制\n使用initializer修饰符保护初始化函数。\n\n### 4. 价格操控\n使用TWAP而非即时价格。\n\n## 审计工具\n\n- Slither\n- Mythril\n- Foundry",
			admin.ID, catSecurity.ID, 1892, 145,
		},
		{
			"Layer2技术全景解析",
			"# Layer2技术全景解析\n\n## 为什么需要Layer2\n\n- 以太坊TPS低（15-45）\n- Gas费用高\n- 用户体验差\n\n## 技术路线\n\n### Optimistic Rollup\n代表：Optimism, Arbitrum\n- 欺诈证明\n- 7天挑战期\n- EVM兼容\n\n### ZK Rollup\n代表：zkSync, StarkNet\n- 零知识证明\n- 即时最终性\n- 数学保证安全",
			admin.ID, catLayer2.ID, 2543, 198,
		},
	}

	for _, a := range articles {
		var existing model.Article
		result := util.Db.Where("title = ?", a.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		article := model.Article{
			Title:      a.Title,
			Content:    a.Content,
			UserID:     a.UserID,
			CategoryID: a.CategoryID,
			ViewCount:  a.ViewCount,
			LikeCount:  a.LikeCount,
		}
		util.Db.Create(&article)
	}
	fmt.Println("文章数据初始化完成！")
}

func initQuestions() {
	var zhangsan, lisi model.User
	util.Db.Where("username = ?", "solidity_dev").First(&zhangsan)
	util.Db.Where("username = ?", "defi_master").First(&lisi)

	questions := []struct {
		Title    string
		Content  string
		UserID   uint
	}{
		{"Solidity中如何实现可升级合约？", "我想实现一个可以升级的智能合约，应该怎么做？有哪些方案？", zhangsan.ID},
		{"Uniswap V3的无常损失如何计算？", "在Uniswap V3中提供流动性，无常损失的计算公式是什么？", lisi.ID},
		{"NFT白名单如何实现？", "我想做一个NFT白名单预售功能，推荐用什么方案？Merkle Tree还是签名验证？", zhangsan.ID},
		{"zkSync和StarkNet有什么区别？", "这两个ZK Rollup方案的技术路线有什么不同？开发者应该选择哪个？", lisi.ID},
		{"如何防止闪电贷攻击？", "DeFi协议如何防止闪电贷价格操控攻击？", zhangsan.ID},
	}

	for _, q := range questions {
		var existing model.Question
		result := util.Db.Where("title = ?", q.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		question := model.Question{
			Title:   q.Title,
			Content: q.Content,
			UserID:  q.UserID,
		}
		util.Db.Create(&question)
	}
	fmt.Println("问答数据初始化完成！")
}

func initCourses() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	courses := []struct {
		Title       string
		Description string
		UserID      uint
	}{
		{"Solidity从入门到精通", "系统学习Solidity智能合约开发，从基础语法到高级特性", admin.ID},
		{"DeFi协议开发实战", "学习开发AMM、借贷协议等DeFi应用", admin.ID},
		{"NFT开发完整指南", "从零开发NFT合约，包含白名单、版税等功能", admin.ID},
		{"智能合约安全审计", "学习智能合约安全漏洞与审计技术", admin.ID},
	}

	for _, c := range courses {
		var existing model.Course
		result := util.Db.Where("title = ?", c.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		course := model.Course{
			Title:       c.Title,
			Description: c.Description,
			UserID:      c.UserID,
		}
		util.Db.Create(&course)
	}
	fmt.Println("教程数据初始化完成！")
}

func initLearningPaths() {
	paths := []struct {
		Title       string
		Description string
		Difficulty  int
		Duration    string
	}{
		{"Web3开发入门路径", "从零开始学习区块链开发和智能合约编程", 1, "2个月"},
		{"DeFi协议开发路径", "深入学习DeFi协议原理和开发实战", 2, "3个月"},
		{"智能合约安全审计路径", "成为专业的智能合约安全审计师", 3, "4个月"},
	}

	for _, p := range paths {
		var existing model.LearningPath
		result := util.Db.Where("title = ?", p.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		path := model.LearningPath{
			Title:       p.Title,
			Description: p.Description,
			Difficulty:  p.Difficulty,
			Duration:    p.Duration,
		}
		util.Db.Create(&path)
	}
	fmt.Println("学习路径数据初始化完成！")
}

func initCodeSnippets() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	snippets := []struct {
		Title       string
		Description string
		Code        string
		Category    string
	}{
		{"ERC20代币合约模板", "标准的ERC20代币合约实现", "pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract MyToken is ERC20 {\n    constructor(uint256 initialSupply) ERC20(\"MyToken\", \"MTK\") {\n        _mint(msg.sender, initialSupply);\n    }\n}", "代币"},
		{"ReentrancyGuard使用", "防止重入攻击的修饰符", "import \"@openzeppelin/contracts/security/ReentrancyGuard.sol\";\n\ncontract SafeContract is ReentrancyGuard {\n    function withdraw() public nonReentrant {\n        // 安全的提款逻辑\n    }\n}", "安全"},
		{"Ownable权限控制", "简单的所有者权限控制", "contract Ownable {\n    address public owner;\n    \n    constructor() {\n        owner = msg.sender;\n    }\n    \n    modifier onlyOwner() {\n        require(msg.sender == owner, \"Not owner\");\n        _;\n    }\n}", "权限"},
	}

	for _, s := range snippets {
		var existing model.CodeSnippet
		result := util.Db.Where("title = ?", s.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		snippet := model.CodeSnippet{
			Title:       s.Title,
			Description: s.Description,
			Code:        s.Code,
			Category:    s.Category,
			UserID:      admin.ID,
		}
		util.Db.Create(&snippet)
	}
	fmt.Println("代码片段数据初始化完成！")
}

func initContractTemplates() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	templates := []struct {
		Name        string
		Description string
		Code        string
		Category    string
	}{
		{"ERC20代币模板", "标准ERC20代币合约", "pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract TokenTemplate is ERC20 {\n    constructor() ERC20(\"Token\", \"TKN\") {\n        _mint(msg.sender, 1000000 * 10 ** decimals());\n    }\n}", "代币"},
		{"ERC721 NFT模板", "标准NFT合约模板", "pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract NFTTemplate is ERC721 {\n    constructor() ERC721(\"NFT\", \"NFT\") {}\n}", "NFT"},
	}

	for _, t := range templates {
		var existing model.ContractTemplate
		result := util.Db.Where("name = ?", t.Name).First(&existing)
		if result.Error == nil {
			continue
		}

		template := model.ContractTemplate{
			Name:        t.Name,
			Description: t.Description,
			Code:        t.Code,
			Category:    t.Category,
			UserID:      admin.ID,
		}
		util.Db.Create(&template)
	}
	fmt.Println("合约模板数据初始化完成！")
}

func initResources() {
	resources := []struct {
		Title       string
		URL         string
		Description string
		Category    string
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
		var existing model.Resource
		result := util.Db.Where("title = ?", r.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		resource := model.Resource{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Category:    r.Category,
		}
		util.Db.Create(&resource)
	}
	fmt.Println("资源导航数据初始化完成！")
}

func initInterviewQuestions() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	questions := []struct {
		Title      string
		Content    string
		Answer     string
		Category   string
		Difficulty int
		Company    string
	}{
		{
			"什么是重入攻击？如何防止？",
			"请解释重入攻击的原理和防御方法。",
			"重入攻击是指合约在更新状态之前调用外部合约，导致外部合约可以重新进入原合约执行恶意操作。\n\n防御方法：\n1. 使用Checks-Effects-Interactions模式\n2. 使用ReentrancyGuard修饰符\n3. 先更新状态再进行外部调用",
			"安全审计", 2, "多家公司",
		},
		{
			"解释Uniswap V2的AMM机制",
			"Uniswap V2如何实现自动做市？",
			"Uniswap V2使用恒定乘积公式 x*y=k 实现自动做市。\n\n核心特点：\n1. 流动性均匀分布在[0,∞]\n2. 资金池由LP提供\n3. 交易手续费0.3%\n4. 价格由储备比例决定",
			"DeFi", 2, "Uniswap",
		},
		{
			"ERC-721和ERC-1155的区别？",
			"比较两种NFT标准的特点和使用场景。",
			"ERC-721:\n- 每个Token唯一\n- 适合收藏品\n- 批量操作Gas高\n\nERC-1155:\n- 支持同质化和非同质化\n- 批量操作高效\n- 适合游戏道具",
			"NFT", 1, "多家公司",
		},
		{
			"什么是闪电贷？有什么应用场景？",
			"解释闪电贷原理和应用。",
			"闪电贷是在一个交易内借款并还款，无需抵押。\n\n原理：利用以太坊原子性，借款和还款在同一交易中完成。\n\n应用场景：\n1. 套利交易\n2. 清算清算\n3. 置换抵押品",
			"DeFi", 2, "Aave",
		},
		{
			"Solidity中storage、memory、calldata的区别？",
			"解释三种数据存储位置的区别。",
			"storage:\n- 永久存储在链上\n- Gas消耗最高\n- 状态变量默认位置\n\nmemory:\n- 临时存储，函数调用后清除\n- Gas消耗中等\n- 引用类型参数默认位置\n\ncalldata:\n- 只读，用于外部函数参数\n- Gas消耗最低\n- 不可修改",
			"智能合约", 1, "基础题",
		},
	}

	for _, q := range questions {
		var existing model.InterviewQuestion
		result := util.Db.Where("title = ?", q.Title).First(&existing)
		if result.Error == nil {
			continue
		}

		question := model.InterviewQuestion{
			Title:      q.Title,
			Content:    q.Content,
			Answer:     q.Answer,
			Category:   q.Category,
			Difficulty: q.Difficulty,
			Company:    q.Company,
			UserID:     admin.ID,
		}
		util.Db.Create(&question)
	}
	fmt.Println("面试题数据初始化完成！")
}
