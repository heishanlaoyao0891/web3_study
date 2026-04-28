package main

import (
	"fmt"
	"go-blog/model"
	"go-blog/util"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := util.CreateTestDB(".env_local"); err != nil {
		panic(err)
	}

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
	tables := []interface{}{
		&model.ArticleTag{},
		&model.Like{},
		&model.Favorite{},
		&model.Answer{},
		&model.Question{},
		&model.CourseChapter{},
		&model.LearningRecord{},
		&model.Course{},
		&model.Checkin{},
		&model.Comment{},
		&model.Article{},
		&model.InterviewQuestion{},
		&model.CodeSnippet{},
		&model.ContractTemplate{},
		&model.Resource{},
		&model.LearningChapter{},
		&model.LearningPath{},
	}

	for _, table := range tables {
		util.Db.Unscoped().Where("1 = 1").Delete(table)
	}

	util.Db.Unscoped().Where("username != ?", "admin").Delete(&model.User{})
	util.Db.Unscoped().Where("1 = 1").Delete(&model.Tag{})
	util.Db.Unscoped().Where("1 = 1").Delete(&model.Category{})

	fmt.Println("✓ 旧数据已清理")
}

func initAdminUser() {
	var admin model.User
	result := util.Db.Where("username = ?", "admin").First(&admin)

	if result.Error != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		admin = model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "Web3技术专家",
			Bio:      "专注Web3开发，智能合约安全审计，DeFi协议研究",
			Level:    5,
			Exp:      800,
			Coins:    500,
			Status:   1,
		}
		util.Db.Create(&admin)
	} else {
		admin.Nickname = "Web3技术专家"
		admin.Bio = "专注Web3开发，智能合约安全审计，DeFi协议研究"
		admin.Level = 5
		admin.Exp = 800
		admin.Coins = 500
		util.Db.Save(&admin)
	}
	fmt.Println("✓ 管理员账号: admin / 123456")
}

func initCategories() {
	categories := []struct {
		Name string
		Desc string
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
		category := model.Category{Name: c.Name, Desc: c.Desc}
		util.Db.Create(&category)
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
		util.Db.Create(&model.Tag{Name: name})
	}
	fmt.Println("✓ 标签数据")
}

func initArticles() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	var catSmartContract, catDeFi, catNFT, catLayer2, catSecurity model.Category
	util.Db.Where("name = ?", "智能合约").First(&catSmartContract)
	util.Db.Where("name = ?", "DeFi协议").First(&catDeFi)
	util.Db.Where("name = ?", "NFT开发").First(&catNFT)
	util.Db.Where("name = ?", "Layer2").First(&catLayer2)
	util.Db.Where("name = ?", "安全审计").First(&catSecurity)

	articles := []struct {
		Title      string
		Content    string
		CategoryID uint
		ViewCount  int
	}{
		{
			"Solidity智能合约开发完整指南",
			"# Solidity智能合约开发完整指南\n\n## 一、环境搭建\n\n### 1. 安装Node.js\n推荐使用Node.js v18+版本。\n\n### 2. 安装Hardhat\n```bash\nnpm install --save-dev hardhat\nnpx hardhat init\n```\n\n### 3. 安装Foundry（推荐）\n```bash\ncurl -L https://foundry.paradigm.xyz | bash\nfoundryup\n```\n\n## 二、第一个智能合约\n\n```solidity\n// SPDX-License-Identifier: MIT\npragma solidity ^0.8.19;\n\ncontract SimpleStorage {\n    uint256 private _value;\n    \n    event ValueChanged(uint256 newValue);\n    \n    function setValue(uint256 value) public {\n        _value = value;\n        emit ValueChanged(value);\n    }\n    \n    function getValue() public view returns (uint256) {\n        return _value;\n    }\n}\n```\n\n## 三、核心概念\n\n### 1. 数据类型\n- **值类型**: bool, int/uint, address, bytes\n- **引用类型**: array, struct, mapping\n\n### 2. 函数可见性\n- public: 内外部均可调用\n- external: 仅外部调用\n- internal: 当前合约及子合约\n- private: 仅当前合约\n\n## 四、安全最佳实践\n\n### 1. 防止重入攻击\n```solidity\nimport \"@openzeppelin/contracts/security/ReentrancyGuard.sol\";\n\ncontract SafeContract is ReentrancyGuard {\n    function withdraw() public nonReentrant {\n        // 安全的提款逻辑\n    }\n}\n```\n\n## 五、Gas优化技巧\n\n1. 使用calldata替代memory\n2. 循环中使用++i替代i++\n3. 缓存数组长度\n4. 使用事件存储历史数据",
			catSmartContract.ID, 3521,
		},
		{
			"Uniswap V3集中流动性原理详解",
			"# Uniswap V3集中流动性原理详解\n\n## 一、V2 vs V3对比\n\n### Uniswap V2\n- 流动性均匀分布在[0, ∞]\n- 资本效率低\n- 简单易用\n\n### Uniswap V3\n- 流动性可集中在指定区间\n- 资本效率提升4000倍\n- 适合专业做市商\n\n## 二、核心机制\n\n### 1. Tick系统\n价格空间被划分为离散的tick：\n```\nprice = 1.0001^tick\n```\n\n### 2. 流动性区间\nLP可选择在[Pa, Pb]区间提供流动性。\n\n## 三、数学原理\n\n### 流动性计算\n```solidity\n// 当价格在区间内\nL = sqrt(x * y)\n```\n\n## 四、实际应用\n\n1. 稳定币做市：在0.99-1.01区间集中流动性\n2. 限价单：在单点提供流动性\n3. 套利策略：监控不同池子价格差",
			catDeFi.ID, 4215,
		},
		{
			"NFT智能合约开发实战",
			"# NFT智能合约开发实战\n\n## 一、ERC-721标准\n\n### 核心接口\n```solidity\ninterface IERC721 {\n    function balanceOf(address owner) external view returns (uint256);\n    function ownerOf(uint256 tokenId) external view returns (address);\n    function transferFrom(address from, address to, uint256 tokenId) external;\n}\n```\n\n### 使用OpenZeppelin\n```solidity\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract MyNFT is ERC721 {\n    constructor() ERC721(\"MyNFT\", \"MNFT\") {}\n    \n    function safeMint(address to, uint256 tokenId) public {\n        _safeMint(to, tokenId);\n    }\n}\n```\n\n## 二、高级功能\n\n### 1. 白名单预售（MerkleTree）\n```solidity\nimport \"@openzeppelin/contracts/utils/cryptography/MerkleProof.sol\";\n\ncontract NFTWhitelist is ERC721 {\n    bytes32 public merkleRoot;\n    \n    function whitelistMint(bytes32[] calldata proof) external {\n        require(MerkleProof.verify(proof, merkleRoot, keccak256(abi.encodePacked(msg.sender))), \"Not whitelisted\");\n        _safeMint(msg.sender, tokenId);\n    }\n}\n```\n\n### 2. 版税机制（ERC-2981）\n支持NFT交易版税自动分成。",
			catNFT.ID, 2856,
		},
		{
			"智能合约安全审计要点",
			"# 智能合约安全审计要点\n\n## 一、常见漏洞类型\n\n### 1. 重入攻击（Reentrancy）\n\n**漏洞示例：**\n```solidity\n// 危险代码！\nfunction withdraw() public {\n    uint256 amount = balances[msg.sender];\n    (bool success,) = msg.sender.call{value: amount}(\"\");\n    balances[msg.sender] = 0; // 状态更新在转账之后\n}\n```\n\n**修复方案：**\n```solidity\nfunction withdraw() public {\n    uint256 amount = balances[msg.sender];\n    balances[msg.sender] = 0; // 先更新状态\n    (bool success,) = msg.sender.call{value: amount}(\"\");\n    require(success);\n}\n```\n\n### 2. 闪电贷价格操控\n使用TWAP（时间加权平均价格）防止价格操控。\n\n### 3. 访问控制漏洞\n使用initializer修饰符保护初始化函数。\n\n## 二、审计工具\n\n- Slither\n- Foundry\n- Mythril",
			catSecurity.ID, 2156,
		},
		{
			"Layer2技术全景解析",
			"# Layer2技术全景解析\n\n## 一、为什么需要Layer2\n\n### 以太坊困境\n- TPS低：约15-45 TPS\n- Gas费高：高峰期几十美元\n\n### Layer2优势\n- 继承以太坊安全性\n- 大幅降低Gas费\n- 提升TPS到数千\n\n## 二、技术路线\n\n### 1. Optimistic Rollup\n代表：Optimism、Arbitrum\n- 7天挑战期\n- 欺诈证明\n- EVM完全兼容\n\n### 2. ZK Rollup\n代表：zkSync、StarkNet\n- 零知识证明\n- 即时最终性\n- 数学保证安全\n\n## 三、主流L2对比\n\n| 特性 | Optimism | Arbitrum | zkSync |\n|------|----------|----------|--------|\n| 最终性 | 7天 | 7天 | 即时 |\n| EVM兼容 | 完全 | 完全 | 大部分 |",
			catLayer2.ID, 1892,
		},
	}

	for _, a := range articles {
		article := model.Article{
			Title:      a.Title,
			Content:    a.Content,
			UserID:     admin.ID,
			CategoryID: a.CategoryID,
			ViewCount:  a.ViewCount,
			LikeCount:  a.ViewCount / 15,
		}
		util.Db.Create(&article)
	}
	fmt.Println("✓ 文章数据")
}

func initQuestions() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	questions := []struct {
		Title   string
		Content string
	}{
		{"Solidity中如何实现可升级合约？", "我想实现一个可以升级的智能合约，应该使用哪种方案？代理模式还是透明代理？"},
		{"Uniswap V3如何计算无常损失？", "在Uniswap V3提供流动性时，无常损失的计算公式是什么？"},
		{"NFT白名单使用MerkleTree还是签名验证？", "NFT预售白名单方案选择：MerkleTree和签名验证各有什么优缺点？"},
		{"zkSync和StarkNet开发有什么区别？", "这两个ZK Rollup在开发体验、性能方面有什么差异？"},
		{"如何防止闪电贷攻击？", "DeFi协议如何防止闪电贷价格操控攻击？"},
	}

	for _, q := range questions {
		question := model.Question{
			Title:   q.Title,
			Content: q.Content,
			UserID:  admin.ID,
		}
		util.Db.Create(&question)
	}
	fmt.Println("✓ 问答数据")
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
	}{
		{
			"什么是重入攻击？如何防御？",
			"请解释重入攻击的原理和防御方法。",
			"重入攻击是指合约在更新状态之前调用外部合约，导致可以重复执行。\n\n防御方法：\n1. Checks-Effects-Interactions模式\n2. ReentrancyGuard修饰符\n3. 使用transfer而非call",
			"安全审计", 2,
		},
		{
			"解释Uniswap V2的AMM机制",
			"Uniswap V2如何实现自动做市？",
			"Uniswap V2使用恒定乘积公式 x*y=k 实现自动做市。\n\n核心公式：\n- 储备恒等式：x * y = k\n- 价格由储备比例决定\n- LP收取0.3%手续费",
			"DeFi协议", 2,
		},
		{
			"ERC-721和ERC-1155的区别？",
			"比较两种NFT标准的特点。",
			"ERC-721:\n- 每个Token唯一\n- 适合收藏品\n\nERC-1155:\n- 支持批量操作\n- 适合游戏道具",
			"NFT开发", 1,
		},
		{
			"Solidity中storage、memory、calldata的区别？",
			"解释三种数据存储位置的特点。",
			"storage: 永久存储，Gas最高\nmemory: 临时存储，可修改\ncalldata: 只读参数，Gas最低",
			"智能合约", 1,
		},
		{
			"什么是闪电贷？",
			"解释闪电贷原理和应用场景。",
			"在一个交易内完成借款、使用、还款，无需抵押。\n\n应用：套利、清算、置换抵押品",
			"DeFi协议", 2,
		},
	}

	for _, q := range questions {
		question := model.InterviewQuestion{
			Title:      q.Title,
			Content:    q.Content,
			Answer:     q.Answer,
			Category:   q.Category,
			Difficulty: q.Difficulty,
			UserID:     admin.ID,
		}
		util.Db.Create(&question)
	}
	fmt.Println("✓ 面试题数据")
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
		{
			"ERC20代币合约模板",
			"标准的ERC20代币合约实现",
			"pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract MyToken is ERC20 {\n    constructor(uint256 supply) ERC20(\"MyToken\", \"MTK\") {\n        _mint(msg.sender, supply);\n    }\n}",
			"代币",
		},
		{
			"ReentrancyGuard使用",
			"防止重入攻击的修饰符",
			"import \"@openzeppelin/contracts/security/ReentrancyGuard.sol\";\n\ncontract SafeContract is ReentrancyGuard {\n    function withdraw() public nonReentrant {\n        // 安全的提款逻辑\n    }\n}",
			"安全",
		},
		{
			"Ownable权限控制",
			"简单的所有者权限控制",
			"contract Ownable {\n    address public owner;\n    \n    constructor() { owner = msg.sender; }\n    \n    modifier onlyOwner() {\n        require(msg.sender == owner, \"Not owner\");\n        _;\n    }\n}",
			"权限",
		},
	}

	for _, s := range snippets {
		snippet := model.CodeSnippet{
			Title:       s.Title,
			Description: s.Description,
			Code:        s.Code,
			Category:    s.Category,
			Language:    "solidity",
			UserID:      admin.ID,
		}
		util.Db.Create(&snippet)
	}
	fmt.Println("✓ 代码片段数据")
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
		{
			"ERC20代币模板",
			"标准ERC20代币合约",
			"pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\n\ncontract TokenTemplate is ERC20 {\n    constructor() ERC20(\"Token\", \"TKN\") {\n        _mint(msg.sender, 1000000 * 10 ** decimals());\n    }\n}",
			"代币",
		},
		{
			"ERC721 NFT模板",
			"标准NFT合约模板",
			"pragma solidity ^0.8.19;\n\nimport \"@openzeppelin/contracts/token/ERC721/ERC721.sol\";\n\ncontract NFTTemplate is ERC721 {\n    constructor() ERC721(\"NFT\", \"NFT\") {}\n}",
			"NFT",
		},
	}

	for _, t := range templates {
		template := model.ContractTemplate{
			Name:        t.Name,
			Description: t.Description,
			Code:        t.Code,
			Category:    t.Category,
			UserID:      admin.ID,
		}
		util.Db.Create(&template)
	}
	fmt.Println("✓ 合约模板数据")
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
		resource := model.Resource{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Category:    r.Category,
		}
		util.Db.Create(&resource)
	}
	fmt.Println("✓ 资源导航数据")
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
		path := model.LearningPath{
			Title:       p.Title,
			Description: p.Description,
			Difficulty:  p.Difficulty,
			Duration:    p.Duration,
		}
		util.Db.Create(&path)
	}
	fmt.Println("✓ 学习路径数据")
}

func initCourses() {
	var admin model.User
	util.Db.Where("username = ?", "admin").First(&admin)

	courses := []struct {
		Title       string
		Description string
		Category    string
	}{
		{"Solidity从入门到精通", "系统学习Solidity智能合约开发，从基础语法到高级特性", "智能合约"},
		{"DeFi协议开发实战", "学习开发AMM、借贷协议等DeFi应用", "DeFi"},
		{"NFT开发完整指南", "从零开发NFT合约，包含白名单、版税等功能", "NFT"},
		{"智能合约安全审计", "学习智能合约安全漏洞与审计技术", "安全"},
	}

	for _, c := range courses {
		course := model.Course{
			Title:       c.Title,
			Description: c.Description,
			Category:    c.Category,
			AuthorID:    admin.ID,
		}
		util.Db.Create(&course)
	}
	fmt.Println("✓ 教程数据")
}
