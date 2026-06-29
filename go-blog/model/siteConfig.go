package model

import (
	"time"

	"gorm.io/gorm"
)

// SiteConfig 站点配置（KV 结构），用于首页文案/品牌/副标题等动态配置
// 管理员通过后台修改，无需改 HTML 硬编码
type SiteConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Key       string         `gorm:"column:config_key;size:50;not null;uniqueIndex" json:"key"` // config_key（key 是 MySQL 保留字）
	Value     string         `gorm:"type:text" json:"value"`                                    // 配置值
	Category  string         `gorm:"size:50;index" json:"category"`                              // 分组：brand / nav / home / footer
	Desc      string         `gorm:"size:200" json:"desc"`                                       // 描述，方便管理员理解
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SiteConfig) TableName() string {
	return "site_configs"
}

// 默认配置键约定（初始化时写入）
const (
	CfgSiteTitle      = "site_title"
	CfgSiteSubtitle   = "site_subtitle"
	CfgSiteIntroTitle = "site_intro_title"
	CfgSiteIntroItems = "site_intro_items"
	CfgFooterText     = "footer_text"
	CfgFooterSubtext  = "footer_subtext"
)