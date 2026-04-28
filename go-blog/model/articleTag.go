package model

import "time"

type ArticleTag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"not null;uniqueIndex:uk_article_tag" json:"article_id"`
	TagID     uint      `gorm:"not null;uniqueIndex:uk_article_tag;index" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
