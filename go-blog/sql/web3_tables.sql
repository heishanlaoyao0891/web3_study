-- Web3平台扩展表

-- 学习路径表
CREATE TABLE IF NOT EXISTS learning_paths (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL COMMENT '路径标题',
    description TEXT COMMENT '路径描述',
    cover VARCHAR(255) COMMENT '封面图',
    difficulty INT DEFAULT 1 COMMENT '难度：1入门 2中级 3高级',
    duration VARCHAR(50) COMMENT '预计时长',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学习路径';

-- 学习章节表
CREATE TABLE IF NOT EXISTS learning_chapters (
    id INT PRIMARY KEY AUTO_INCREMENT,
    path_id INT NOT NULL COMMENT '路径ID',
    title VARCHAR(100) NOT NULL COMMENT '章节标题',
    content TEXT COMMENT '章节内容',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_path_id (path_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学习章节';

-- 代码片段表
CREATE TABLE IF NOT EXISTS code_snippets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL COMMENT '标题',
    description TEXT COMMENT '描述',
    code TEXT NOT NULL COMMENT '代码内容',
    language VARCHAR(50) DEFAULT 'solidity' COMMENT '编程语言',
    category VARCHAR(50) COMMENT '分类',
    tags VARCHAR(255) COMMENT '标签，逗号分隔',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    user_id INT NOT NULL COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_category (category),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代码片段';

-- 合约模板表
CREATE TABLE IF NOT EXISTS contract_templates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    description TEXT COMMENT '模板描述',
    code TEXT NOT NULL COMMENT '合约代码',
    category VARCHAR(50) COMMENT '分类',
    tags VARCHAR(255) COMMENT '标签',
    view_count INT DEFAULT 0,
    download_count INT DEFAULT 0 COMMENT '下载/使用次数',
    user_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_category (category),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='合约模板';

-- 资源导航表
CREATE TABLE IF NOT EXISTS resources (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL COMMENT '资源标题',
    url VARCHAR(500) NOT NULL COMMENT '链接地址',
    description TEXT COMMENT '描述',
    category VARCHAR(50) COMMENT '分类：文档/工具/教程/社区',
    icon VARCHAR(255) COMMENT '图标',
    sort_order INT DEFAULT 0,
    view_count INT DEFAULT 0,
    click_count INT DEFAULT 0 COMMENT '点击次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_category (category),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源导航';

-- 面试题表
CREATE TABLE IF NOT EXISTS interview_questions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '问题标题',
    content TEXT NOT NULL COMMENT '问题描述',
    answer TEXT COMMENT '参考答案',
    category VARCHAR(50) COMMENT '分类：智能合约/DeFi/NFT/安全/Layer2',
    difficulty INT DEFAULT 1 COMMENT '难度：1简单 2中等 3困难',
    company VARCHAR(100) COMMENT '出自哪家公司',
    tags VARCHAR(255) COMMENT '标签',
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    user_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_category (category),
    INDEX idx_difficulty (difficulty),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='面试题库';
