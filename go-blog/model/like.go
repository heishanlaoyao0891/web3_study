package model

import "time"

type Like struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex:uk_user_target" json:"user_id"`
	TargetType int       `gorm:"not null;uniqueIndex:uk_user_target" json:"target_type"`
	TargetID   uint      `gorm:"not null;uniqueIndex:uk_user_target;index:idx_target" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
