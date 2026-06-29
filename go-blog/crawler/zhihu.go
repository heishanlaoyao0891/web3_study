package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go-blog/model"
)

// ZhihuCrawler 知乎热榜爬虫适配器
// 调用知乎热榜 API 获取热门话题
type ZhihuCrawler struct {
	client *http.Client
}

func NewZhihuCrawler() *ZhihuCrawler {
	return &ZhihuCrawler{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (z *ZhihuCrawler) Name() string {
	return "zhihu"
}

// Fetch 获取知乎热榜
func (z *ZhihuCrawler) Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error) {
	apiURL := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total?limit=30"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HermesCrawler/1.0)")
	req.Header.Set("Accept", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result zhihuResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("知乎 API 返回空数据")
	}

	items := make([]RawItem, 0, len(result.Data))
	for _, item := range result.Data {
		target := item.Target
		if target == nil || target.Title == "" {
			continue
		}

		// 知乎问题/文章链接
		link := target.URL
		if link == "" {
			link = fmt.Sprintf("https://www.zhihu.com/question/%d", target.ID)
		}
		// 确保是完整 URL
		if len(link) > 0 && link[0] == '/' {
			link = "https://www.zhihu.com" + link
		}

		// 摘要
		content := target.Excerpt
		if content == "" {
			content = target.Description
		}

		// 热榜专用 ID 去重
		rawID := fmt.Sprintf("%d", item.ID)
		if rawID == "0" {
			rawID = fmt.Sprintf("%d", target.ID)
		}

		items = append(items, RawItem{
			Title:       target.Title,
			URL:         link,
			Content:     truncateStr(content, 300),
			RawID:       rawID,
			PublishedAt: time.Now(),
		})
	}

	return items, nil
}

// ========================================
// 知乎 API 响应结构
// ========================================

type zhihuResponse struct {
	Data []zhihuDataItem `json:"data"`
}

type zhihuDataItem struct {
	ID     int64          `json:"id"`
	Type   string         `json:"type"`
	Target *zhihuTarget   `json:"target"`
}

type zhihuTarget struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Excerpt     string `json:"excerpt"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func init() {
	log.Println("[crawler] 注册知乎适配器")
	Register(NewZhihuCrawler())
}
