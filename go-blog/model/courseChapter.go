package model

import "time"

type CourseChapter struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CourseID   uint      `gorm:"not null;index" json:"course_id"`
	Title      string    `gorm:"size:200;not null" json:"title"`
	Content    string    `gorm:"type:longtext" json:"content"`
	SortOrder  int       `gorm:"default:0;index" json:"sort_order"`
	ViewCount  int       `gorm:"default:0" json:"view_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
