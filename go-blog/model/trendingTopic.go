package model

import (
	"time"

	"gorm.io/gorm"
)

// TrendingTopic 风口话题（M3.1）
// 管理员手工录入或后续从热点站抓取，首页"风口话题"专区展示
type TrendingTopic struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`      // 话题标题
	Summary    string         `gorm:"type:text" json:"summary"`            // 简要描述
	URL        string         `gorm:"size:500" json:"url"`                 // 相关链接
	Source     string         `gorm:"size:50" json:"source"`               // 来源标识
	DomainID   uint           `gorm:"index" json:"domain_id"`             // 归属技术领域
	HeatScore  int            `gorm:"default:0" json:"heat_score"`         // 热度分值
	Status     int            `gorm:"default:1;index" json:"status"`      // 1=展示中 0=下线
	ExpireAt   *time.Time     `json:"expire_at"`                            // 过期时间，nil=不过期
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TrendingTopic) TableName() string {
	return "trending_topics"
}