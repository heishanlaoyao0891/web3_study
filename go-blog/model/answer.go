package model

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	QuestionID  uint           `gorm:"not null;index" json:"question_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	IsBest      int            `gorm:"default:0;index" json:"is_best"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	Question    Question       `gorm:"foreignKey:QuestionID" json:"question"`
}
