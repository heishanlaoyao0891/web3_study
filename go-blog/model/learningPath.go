package model

import (
	"time"

	"gorm.io/gorm"
)

type LearningPath struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Cover       string         `gorm:"size:255" json:"cover"`
	Difficulty  int            `gorm:"default:1" json:"difficulty"`
	Duration    string         `gorm:"size:50" json:"duration"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Chapters    []LearningChapter `gorm:"foreignKey:PathID" json:"chapters"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LearningPath) TableName() string {
	return "learning_paths"
}

type LearningChapter struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	PathID    uint           `gorm:"not null;index" json:"path_id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LearningChapter) TableName() string {
	return "learning_chapters"
}
