package model

import (
	"time"

	"gorm.io/gorm"
)

type ArticleCategory struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ArticleID  uint           `gorm:"not null;index:idx_article_category" json:"article_id"`
	CategoryID uint           `gorm:"not null;index:idx_article_category" json:"category_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	
	Article  Article  `gorm:"foreignKey:ArticleID" json:"article"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}

type InterviewQuestionCategory struct {
	ID                   uint           `gorm:"primarykey" json:"id"`
	InterviewQuestionID  uint           `gorm:"not null;index:idx_question_category" json:"interview_question_id"`
	CategoryID           uint           `gorm:"not null;index:idx_question_category" json:"category_id"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	
	InterviewQuestion InterviewQuestion `gorm:"foreignKey:InterviewQuestionID" json:"interview_question"`
	Category          Category          `gorm:"foreignKey:CategoryID" json:"category"`
}

func (InterviewQuestionCategory) TableName() string {
	return "interview_question_categories"
}
