package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go-blog/model"
)

// JuejinCrawler 掘金热门文章爬虫适配器
// 调用掘金推荐 API 获取热门技术文章
type JuejinCrawler struct {
	client *http.Client
}

func NewJuejinCrawler() *JuejinCrawler {
	return &JuejinCrawler{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (j *JuejinCrawler) Name() string {
	return "juejin"
}

// Fetch 从掘金推荐 API 获取热门文章列表
func (j *JuejinCrawler) Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error) {
	// 掘金推荐 API（cate_id=0 表示全部，limit=20）
	apiURL := "https://api.juejin.cn/recommend_api/v1/article/recommend_all_feed"

	payload := strings.NewReader(`{"cate_id":"0","limit":20}`)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, payload)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HermesCrawler/1.0)")
	req.Header.Set("Origin", "https://juejin.cn")
	req.Header.Set("Referer", "https://juejin.cn/")

	resp, err := j.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result juejinResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	if result.ErrNo != 0 {
		return nil, fmt.Errorf("掘金 API 返回错误: %d %s", result.ErrNo, result.ErrMsg)
	}

	items := make([]RawItem, 0, len(result.Data))
	for _, d := range result.Data {
		item := d.ItemInfo
		if item == nil || item.ArticleInfo == nil {
			continue
		}
		a := item.ArticleInfo
		if a.Title == "" || a.ArticleID == "" {
			continue
		}

		// 掘金文章 URL
		articleURL := fmt.Sprintf("https://juejin.cn/post/%s", a.ArticleID)

		// 取内容摘要
		content := a.BriefContent
		if content == "" {
			content = a.Summary
		}

		items = append(items, RawItem{
			Title:       a.Title,
			URL:         articleURL,
			Content:     truncateStr(content, 500),
			RawID:       a.ArticleID,
			PublishedAt: time.UnixMilli(a.CTime * 1000),
		})
	}

	return items, nil
}

// ========================================
// 掘金 API 响应结构
// ========================================

type juejinResponse struct {
	ErrNo  int              `json:"err_no"`
	ErrMsg string           `json:"err_msg"`
	Data   []juejinDataItem `json:"data"`
}

type juejinDataItem struct {
	ItemInfo *juejinItemInfo `json:"item_info"`
}

type juejinItemInfo struct {
	ArticleInfo *juejinArticle `json:"article_info"`
}

type juejinArticle struct {
	ArticleID    string `json:"article_id"`
	Title        string `json:"title"`
	BriefContent string `json:"brief_content"`
	Summary      string `json:"summary"`
	CTime        int64  `json:"ctime"` // Unix 秒
	ViewCount    int    `json:"view_count"`
	DiggCount    int    `json:"digg_count"`
	CommentCount int    `json:"comment_count"`
	UserID       string `json:"user_id"`
	User         *juejinUser `json:"author_user_info"`
}

type juejinUser struct {
	UserName   string `json:"user_name"`
	Avatar     string `json:"avatar_large"`
}

func init() {
	log.Println("[crawler] 注册掘金适配器")
	Register(NewJuejinCrawler())
}

// truncateStr 截断字符串
func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
