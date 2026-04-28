package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ArticleID uint           `gorm:"not null;index" json:"article_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ParentID  *uint          `json:"parent_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Article   Article        `gorm:"foreignKey:ArticleID" json:"article"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies"`
}
