package model

import (
	"time"

	"gorm.io/gorm"
)

type CodeSnippet struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Code        string         `gorm:"type:text;not null" json:"code"`
	Language    string         `gorm:"size:50;default:solidity" json:"language"`
	Category    string         `gorm:"size:50" json:"category"`
	Tags        string         `gorm:"size:255" json:"tags"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CodeSnippet) TableName() string {
	return "code_snippets"
}

type ContractTemplate struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	Code          string         `gorm:"type:text;not null" json:"code"`
	Category      string         `gorm:"size:50" json:"category"`
	Tags          string         `gorm:"size:255" json:"tags"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContractTemplate) TableName() string {
	return "contract_templates"
}
