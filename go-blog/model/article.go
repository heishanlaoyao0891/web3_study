package model

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章实体
type Article struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"size:100;not null" json:"title"`    // 文章标题
	Content    string         `gorm:"type:text;not null" json:"content"` // 文章内容
	Status     int            `gorm:"default:1" json:"status"`           // 状态：1-发布，0-草稿
	UserID     uint           `gorm:"not null" json:"user_id"`           // 外键：关联用户
	CategoryID uint           `gorm:"not null" json:"category_id"`       // 外键：关联分类
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       User           `gorm:"foreignKey:UserID" json:"user"`         // 多对一：关联用户
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"` // 多对一：关联分类
}
