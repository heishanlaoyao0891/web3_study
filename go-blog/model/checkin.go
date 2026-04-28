package model

import "time"

type Checkin struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;uniqueIndex:uk_user_date" json:"user_id"`
	CheckinDate time.Time `gorm:"type:date;not null;uniqueIndex:uk_user_date;index" json:"checkin_date"`
	ExpGained   int       `gorm:"default:10" json:"exp_gained"`
	CoinsGained int       `gorm:"default:5" json:"coins_gained"`
	CreatedAt   time.Time `json:"created_at"`
}
