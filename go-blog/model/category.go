package model

import (
	"gorm.io/gorm"
	"time"
)

// Category 文章分类实体
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"` // 分类名，唯一
	Desc      string         `gorm:"size:200" json:"desc"`                // 分类描述
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Articles  []Article      `gorm:"foreignKey:CategoryID" json:"articles"` // 一对多：一个分类多篇文章
}
