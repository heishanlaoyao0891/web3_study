package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	ViewCount    int            `gorm:"default:0" json:"view_count"`
	AnswerCount  int            `gorm:"default:0" json:"answer_count"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	Status       int            `gorm:"default:0;index" json:"status"`
	BestAnswerID *uint          `json:"best_answer_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	Answers      []Answer       `gorm:"foreignKey:QuestionID" json:"answers"`
}
