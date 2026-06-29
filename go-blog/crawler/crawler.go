package crawler

import (
	"context"
	"fmt"
	"time"

	"go-blog/model"
	"go-blog/util"
)

// RawItem 适配器返回的原始条目（与具体站点无关的统一结构）
type RawItem struct {
	Title       string
	URL         string    // 原文链接
	Content     string    // 摘要或正文
	RawID       string    // 源站唯一ID，用于去重
	PublishedAt time.Time // 发布时间
}

// Crawler 适配器接口：每个站点实现一个
// 框架负责调度 → 调 Fetch → 去重 → 转 Article → 入库 → 写 CrawlLog
// 适配器只管：从某站点把数据取回来装进 []RawItem
type Crawler interface {
	Name() string // 适配器类型标识，如 "hackernews"
	Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error)
}

// Registry 适配器注册表
var registry = make(map[string]Crawler)

// Register 注册适配器（在 init 中调用）
func Register(c Crawler) {
	registry[c.Name()] = c
}

// Get 按类型获取适配器
func Get(typeName string) (Crawler, error) {
	c, ok := registry[typeName]
	if !ok {
		return nil, fmt.Errorf("crawler not found: %s", typeName)
	}
	return c, nil
}

// RunOnce 执行一次抓取（供调度器或手动触发调用）
// 完整流程：Fetch → 去重 → 入库 → 写日志
func RunOnce(src model.ContentSource) (*model.CrawlLog, error) {
	start := time.Now()
	log := &model.CrawlLog{
		SourceID:   src.ID,
		SourceName: src.Name,
		Status:     "failed",
	}

	// 获取适配器
	c, err := Get(src.Type)
	if err != nil {
		log.Error = err.Error()
		saveLog(log, start)
		return log, err
	}

	// 抓取
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	items, err := c.Fetch(ctx, src)
	if err != nil {
		log.Error = err.Error()
		log.FetchedCount = len(items)
		saveLog(log, start)
		return log, err
	}

	log.FetchedCount = len(items)

	// 去重 + 入库
	saved, dup := saveItems(items, src)
	log.SavedCount = saved
	log.DuplicateCount = dup
	log.Status = "success"

	saveLog(log, start)
	return log, nil
}

// saveItems 将 RawItem 转为 Article 入库，按 Source+RawID 去重
// 返回：成功入库条数，重复跳过条数
func saveItems(items []RawItem, src model.ContentSource) (int, int) {
	saved, dup := 0, 0
	for _, item := range items {
		if item.RawID == "" || item.Title == "" {
			continue
		}

		// 去重：同 source + raw_id 的不重复入库
		var count int64
		util.Db.Model(&model.Article{}).
			Where("source = ? AND raw_id = ?", src.Type, item.RawID).
			Count(&count)
		if count > 0 {
			dup++
			continue
		}

		article := model.Article{
			Title:      item.Title,
			Content:    item.Content,
			Status:     1,
			Visibility: 1, // 抓取的文章默认公开
			UserID:     1, // admin 用户
			CategoryID: src.CategoryID,
			DomainID:   src.DomainID,
			Source:     src.Type,
			SourceURL:  item.URL,
			RawID:      item.RawID,
		}
		if err := util.Db.Create(&article).Error; err != nil {
			continue
		}
		saved++
	}
	return saved, dup
}

// saveLog 保存抓取日志并更新 ContentSource 的 LastRunAt
func saveLog(log *model.CrawlLog, start time.Time) {
	log.Duration = int(time.Since(start).Milliseconds())
	util.Db.Create(log)

	// 更新源的 LastRunAt
	util.Db.Model(&model.ContentSource{}).
		Where("id = ?", log.SourceID).
		Update("last_run_at", time.Now())
}