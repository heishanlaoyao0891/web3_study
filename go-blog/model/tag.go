package model

import "time"

type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null;unique" json:"name"`
	Color     string    `gorm:"size:20;default:'#667eea'" json:"color"`
	Icon      string    `gorm:"size:100" json:"icon"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	UseCount  int       `gorm:"default:0" json:"use_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
