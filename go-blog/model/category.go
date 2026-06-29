package model

import (
	"gorm.io/gorm"
	"time"
)

// Category 文章分类实体
// 顶层 Category（ParentID 为空）= 技术领域（Web3 / Java / Go / Python / AI …）
// 子 Category（ParentID 指向顶层）= 主题（如 Web3 → Solidity/DeFi；Java → Spring/并发）
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"` // 分类名，唯一
	Desc      string         `gorm:"size:200" json:"desc"`                // 分类描述
	ParentID  *uint          `gorm:"index" json:"parent_id"`              // 父分类ID，nil 表示顶层（技术领域）
	SortOrder int            `gorm:"default:0" json:"sort_order"`         // 排序，越大越靠前
	Icon      string         `gorm:"size:50" json:"icon"`                 // 图标（emoji 或 CSS class）
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Articles  []Article      `gorm:"foreignKey:CategoryID" json:"articles"`  // 一对多：一个分类多篇文章
	Children  []Category     `gorm:"foreignKey:ParentID" json:"children"`    // 二级子分类
	Parent    *Category      `gorm:"foreignKey:ParentID" json:"-"`           // 父分类
}

// IsDomain 是否为技术领域（顶层分类）
func (c *Category) IsDomain() bool {
	return c.ParentID == nil
}
