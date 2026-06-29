package model

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Title         string         `gorm:"size:100;not null" json:"title"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Cover         string         `gorm:"size:255" json:"cover"`
	Status        int            `gorm:"default:1" json:"status"`
	Visibility    int            `json:"visibility"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	CategoryID    uint           `gorm:"not null" json:"category_id"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	FavoriteCount int            `gorm:"default:0" json:"favorite_count"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`
	// 抓取来源标记（M2.3）：手动创建的文章这些字段为零值
	Source    string `gorm:"size:50;index" json:"source"`   // 来源标识，如 "hackernews" / "juejin"
	SourceURL string `gorm:"size:500" json:"source_url"`    // 原文链接
	RawID     string `gorm:"size:100;index" json:"raw_id"`  // 源站唯一ID，去重用
	DomainID  uint   `gorm:"index" json:"domain_id"`        // 归属技术领域（顶层Category ID）
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Categories    []Category     `gorm:"many2many:article_categories;" json:"categories"`
	Tags          []Tag          `gorm:"many2many:article_tags;" json:"tags"`
}