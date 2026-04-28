package model

import "time"

type LearningRecord struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"not null;uniqueIndex:uk_user_course" json:"user_id"`
	CourseID    uint       `gorm:"not null;uniqueIndex:uk_user_course;index" json:"course_id"`
	ChapterID   *uint      `json:"chapter_id"`
	Progress    int        `gorm:"default:0" json:"progress"`
	LastLearnAt *time.Time `json:"last_learn_at"`
	IsCompleted int        `gorm:"default:0" json:"is_completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
