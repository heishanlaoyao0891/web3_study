package model

import (
	"time"

	"gorm.io/gorm"
)

// ContentSource 抓取源配置（管理员管理）
type ContentSource struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`       // 源名称，如 "Hacker News"
	Type      string         `gorm:"size:50;not null" json:"type"`       // 适配器类型：hackernews / juejin / zhihu
	URL       string         `gorm:"size:500" json:"url"`                 // 源地址（API 或页面URL）
	DomainID  uint           `gorm:"index" json:"domain_id"`             // 抓回的文章归属哪个技术领域
	CategoryID uint          `gorm:"index" json:"category_id"`           // 抓回的文章归属哪个子分类（0=不指定）
	Cron      string         `gorm:"size:50;default:'0 * * * *'" json:"cron"` // 调度表达式（cron 格式）
	Enabled   bool           `gorm:"default:true" json:"enabled"`        // 是否启用
	LastRunAt *time.Time     `json:"last_run_at"`                         // 上次执行时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContentSource) TableName() string {
	return "content_sources"
}

// CrawlLog 抓取日志（每次执行一条）
type CrawlLog struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	SourceID      uint           `gorm:"not null;index" json:"source_id"`  // 关联 ContentSource
	SourceName    string         `gorm:"size:100" json:"source_name"`      // 冗余源名称，方便查询
	Status        string         `gorm:"size:20;index" json:"status"`      // success / failed / partial
	FetchedCount  int            `json:"fetched_count"`                     // 抓取条数
	SavedCount    int            `json:"saved_count"`                       // 入库条数（去重后）
	DuplicateCount int           `json:"duplicate_count"`                   // 重复跳过条数
	Error         string         `gorm:"type:text" json:"error"`            // 错误信息
	Duration      int            `json:"duration"`                          // 耗时（毫秒）
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CrawlLog) TableName() string {
	return "crawl_logs"
}