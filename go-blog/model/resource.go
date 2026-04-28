package model

import (
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"size:100;not null" json:"title"`
	URL        string         `gorm:"size:500;not null" json:"url"`
	Description string        `gorm:"type:text" json:"description"`
	Category   string         `gorm:"size:50" json:"category"`
	Icon       string         `gorm:"size:255" json:"icon"`
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	ClickCount int            `gorm:"default:0" json:"click_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Resource) TableName() string {
	return "resources"
}

type InterviewQuestion struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Answer     string         `gorm:"type:text" json:"answer"`
	Category   string         `gorm:"size:50" json:"category"`
	Difficulty int            `gorm:"default:1" json:"difficulty"`
	Company    string         `gorm:"size:100" json:"company"`
	Tags       string         `gorm:"size:255" json:"tags"`
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user"`
	Categories []Category     `gorm:"many2many:interview_question_categories;" json:"categories"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (InterviewQuestion) TableName() string {
	return "interview_questions"
}
