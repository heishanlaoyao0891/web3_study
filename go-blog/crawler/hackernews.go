package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go-blog/model"
)

// HackerNewsCrawler 适配器：从 Hacker News API 抓取热门文章
// API 文档：https://github.com/HackerNews/API
type HackerNewsCrawler struct{}

func (h *HackerNewsCrawler) Name() string {
	return "hackernews"
}

// Fetch 从 Hacker News 抓取 top stories
func (h *HackerNewsCrawler) Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error) {
	// 获取 Top Stories IDs
	topIDs, err := h.fetchTopIDs(ctx, 30) // 抓 top 30
	if err != nil {
		return nil, fmt.Errorf("获取 top stories 失败: %w", err)
	}

	items := make([]RawItem, 0, len(topIDs))
	for _, id := range topIDs {
		item, err := h.fetchItem(ctx, id)
		if err != nil {
			continue // 单条失败跳过
		}
		if item.Title == "" || item.URL == "" {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// fetchTopIDs 获取 Top Stories 的 ID 列表
func (h *HackerNewsCrawler) fetchTopIDs(ctx context.Context, limit int) ([]int, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	if limit > 0 && limit < len(ids) {
		ids = ids[:limit]
	}
	return ids, nil
}

// fetchItem 获取单条文章详情
func (h *HackerNewsCrawler) fetchItem(ctx context.Context, id int) (RawItem, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)

	resp, err := http.Get(url)
	if err != nil {
		return RawItem{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RawItem{}, err
	}

	var hnItem struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
		Time  int64  `json:"time"` // Unix timestamp
		Score int    `json:"score"`
	}
	if err := json.Unmarshal(body, &hnItem); err != nil {
		return RawItem{}, err
	}

	// 如果没有 URL（Ask HN / Show HN），用 HN 链接
	articleURL := hnItem.URL
	if articleURL == "" {
		articleURL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", id)
	}

	content := fmt.Sprintf("Score: %d | Source: Hacker News", hnItem.Score)

	return RawItem{
		Title:       hnItem.Title,
		URL:         articleURL,
		Content:     content,
		RawID:       strconv.Itoa(hnItem.ID),
		PublishedAt: time.Unix(hnItem.Time, 0),
	}, nil
}

func init() {
	slog.Info("注册 Hacker News 适配器")
	Register(&HackerNewsCrawler{})
}