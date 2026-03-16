package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户实体（对应Java的UserEntity）
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`                    // 主键
	Username  string         `gorm:"size:50;not null;unique" json:"username"` // 用户名，唯一非空
	Password  string         `gorm:"size:100;not null" json:"password"`       // 密码（建议加密存储）
	Nickname  string         `gorm:"size:50" json:"nickname"`                 // 昵称
	CreatedAt time.Time      `json:"created_at"`                              // 创建时间（Gorm自动维护）
	UpdatedAt time.Time      `json:"updated_at"`                              // 更新时间（Gorm自动维护）
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除（json:"-" 不返回给前端）
	Articles  []Article      `gorm:"foreignKey:UserID" json:"articles"`       // 一对多：一个用户多篇文章
}
