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
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Tags          []Tag          `gorm:"many2many:article_tags;" json:"tags"`
}
