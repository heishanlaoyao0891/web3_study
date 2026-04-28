package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Username      string         `gorm:"size:50;not null;unique" json:"username"`
	Password      string         `gorm:"size:100;not null" json:"-"`
	Nickname      string         `gorm:"size:50" json:"nickname"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	Bio           string         `gorm:"type:text" json:"bio"`
	Status        int            `gorm:"default:1" json:"status"`
	DisableUntil  *time.Time     `json:"disable_until"`
	Level         int            `gorm:"default:1" json:"level"`
	Exp           int            `gorm:"default:0" json:"exp"`
	Coins         int            `gorm:"default:0" json:"coins"`
	CheckinDays   int            `gorm:"default:0" json:"checkin_days"`
	LastCheckinAt *time.Time     `json:"last_checkin_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Articles      []Article      `gorm:"foreignKey:UserID" json:"articles"`
}

func (u *User) LevelName() string {
	names := []string{"新手", "入门", "初级", "中级", "高级", "资深", "专家", "大师", "宗师", "传奇"}
	if u.Level >= 1 && u.Level <= 10 {
		return names[u.Level-1]
	}
	return "新手"
}

func (u *User) ExpToNextLevel() int {
	return u.Level * 100
}
