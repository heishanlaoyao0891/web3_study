package model

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Cover         string         `gorm:"size:255" json:"cover"`
	Category      string         `gorm:"size:50;not null;index" json:"category"`
	AuthorID      uint           `gorm:"not null;index" json:"author_id"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	FavoriteCount int            `gorm:"default:0" json:"favorite_count"`
	ChapterCount  int            `gorm:"default:0" json:"chapter_count"`
	IsFree        int            `gorm:"default:1" json:"is_free"`
	Priority      int            `gorm:"default:0;index" json:"priority"`
	Status        int            `gorm:"default:1" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Author        User           `gorm:"foreignKey:AuthorID" json:"author"`
	Chapters      []CourseChapter `gorm:"foreignKey:CourseID" json:"chapters"`
}
