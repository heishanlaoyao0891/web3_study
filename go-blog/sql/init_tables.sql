-- ========================================
-- 1. 用户表扩展字段
-- ========================================
ALTER TABLE users ADD COLUMN avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL';
ALTER TABLE users ADD COLUMN bio TEXT COMMENT '个人简介';
ALTER TABLE users ADD COLUMN level INT DEFAULT 1 COMMENT '等级 1-10';
ALTER TABLE users ADD COLUMN exp INT DEFAULT 0 COMMENT '经验值';
ALTER TABLE users ADD COLUMN coins INT DEFAULT 0 COMMENT '金币/积分';
ALTER TABLE users ADD COLUMN checkin_days INT DEFAULT 0 COMMENT '连续打卡天数';
ALTER TABLE users ADD COLUMN last_checkin_at DATETIME COMMENT '最后打卡时间';

-- ========================================
-- 2. 标签表
-- ========================================
CREATE TABLE tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '标签名',
    color VARCHAR(20) DEFAULT '#667eea' COMMENT '标签颜色',
    icon VARCHAR(100) DEFAULT '' COMMENT '标签图标',
    sort_order INT DEFAULT 0 COMMENT '排序',
    use_count INT DEFAULT 0 COMMENT '使用次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_name (name),
    KEY idx_use_count (use_count)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

-- ========================================
-- 3. 文章标签关联表
-- ========================================
CREATE TABLE article_tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    article_id BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    tag_id BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_article_tag (article_id, tag_id),
    KEY idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签关联表';

-- ========================================
-- 4. 点赞表（通用）
-- ========================================
CREATE TABLE likes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    target_type TINYINT NOT NULL COMMENT '目标类型: 1文章 2评论 3问答',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_target (user_id, target_type, target_id),
    KEY idx_target (target_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

-- ========================================
-- 5. 收藏表
-- ========================================
CREATE TABLE favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    target_type TINYINT NOT NULL COMMENT '目标类型: 1文章 2教程',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_target (user_id, target_type, target_id),
    KEY idx_target (target_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';

-- ========================================
-- 6. 问答表
-- ========================================
CREATE TABLE questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '提问者ID',
    title VARCHAR(200) NOT NULL COMMENT '问题标题',
    content TEXT NOT NULL COMMENT '问题描述',
    view_count INT DEFAULT 0 COMMENT '浏览数',
    answer_count INT DEFAULT 0 COMMENT '回答数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    status TINYINT DEFAULT 0 COMMENT '状态: 0待解决 1已解决 2已关闭',
    best_answer_id BIGINT UNSIGNED DEFAULT NULL COMMENT '最佳答案ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_user_id (user_id),
    KEY idx_status (status),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='问答表';

-- ========================================
-- 7. 回答表
-- ========================================
CREATE TABLE answers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    question_id BIGINT UNSIGNED NOT NULL COMMENT '问题ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '回答者ID',
    content TEXT NOT NULL COMMENT '回答内容',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    is_best TINYINT DEFAULT 0 COMMENT '是否最佳答案',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_question_id (question_id),
    KEY idx_user_id (user_id),
    KEY idx_is_best (is_best)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='回答表';

-- ========================================
-- 8. 教程表
-- ========================================
CREATE TABLE courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '教程标题',
    description TEXT COMMENT '教程描述',
    cover VARCHAR(255) DEFAULT '' COMMENT '封面图',
    category VARCHAR(50) NOT NULL COMMENT '分类: Java/Go/Python/前端/后端等',
    author_id BIGINT UNSIGNED NOT NULL COMMENT '作者ID',
    view_count INT DEFAULT 0 COMMENT '浏览数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    favorite_count INT DEFAULT 0 COMMENT '收藏数',
    chapter_count INT DEFAULT 0 COMMENT '章节数',
    is_free TINYINT DEFAULT 1 COMMENT '是否免费',
    priority INT DEFAULT 0 COMMENT '优先级/排序',
    status TINYINT DEFAULT 1 COMMENT '状态: 0草稿 1发布',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_category (category),
    KEY idx_author_id (author_id),
    KEY idx_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教程表';

-- ========================================
-- 9. 教程章节表
-- ========================================
CREATE TABLE course_chapters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    course_id BIGINT UNSIGNED NOT NULL COMMENT '教程ID',
    title VARCHAR(200) NOT NULL COMMENT '章节标题',
    content LONGTEXT COMMENT '章节内容',
    sort_order INT DEFAULT 0 COMMENT '排序',
    view_count INT DEFAULT 0 COMMENT '浏览数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_course_id (course_id),
    KEY idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教程章节表';

-- ========================================
-- 10. 学习记录表
-- ========================================
CREATE TABLE learning_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    course_id BIGINT UNSIGNED NOT NULL COMMENT '教程ID',
    chapter_id BIGINT UNSIGNED DEFAULT NULL COMMENT '章节ID',
    progress INT DEFAULT 0 COMMENT '进度百分比',
    last_learn_at DATETIME COMMENT '最后学习时间',
    is_completed TINYINT DEFAULT 0 COMMENT '是否完成',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_course (user_id, course_id),
    KEY idx_course_id (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学习记录表';

-- ========================================
-- 11. 打卡记录表
-- ========================================
CREATE TABLE checkins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    checkin_date DATE NOT NULL COMMENT '打卡日期',
    exp_gained INT DEFAULT 10 COMMENT '获得经验值',
    coins_gained INT DEFAULT 5 COMMENT '获得金币',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_date (user_id, checkin_date),
    KEY idx_checkin_date (checkin_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='打卡记录表';

-- ========================================
-- 12. 文章表扩展字段
-- ========================================
ALTER TABLE articles ADD COLUMN cover VARCHAR(255) DEFAULT '' COMMENT '封面图';
ALTER TABLE articles ADD COLUMN view_count INT DEFAULT 0 COMMENT '浏览数';
ALTER TABLE articles ADD COLUMN like_count INT DEFAULT 0 COMMENT '点赞数';
ALTER TABLE articles ADD COLUMN favorite_count INT DEFAULT 0 COMMENT '收藏数';
ALTER TABLE articles ADD COLUMN comment_count INT DEFAULT 0 COMMENT '评论数';

-- ========================================
-- 13. 初始化标签数据
-- ========================================
INSERT INTO tags (name, color, sort_order) VALUES
('Java', '#f89820', 1),
('Go', '#00add8', 2),
('Python', '#3776ab', 3),
('前端', '#42b883', 4),
('后端', '#6c757d', 5),
('数据库', '#336791', 6),
('AI', '#ff6b6b', 7),
('求职', '#667eea', 8),
('面经', '#764ba2', 9),
('学习打卡', '#28a745', 10);
