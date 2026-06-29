package crawler

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go-blog/model"
)

// ========================================
// RSS/Atom XML 结构（Go 标准库 encoding/xml 解析）
// ========================================

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Content     string   `xml:"encoded"`
	PubDate     string   `xml:"pubDate"`
	Guid        string   `xml:"guid"`
	Categories  []string `xml:"category"`
}

type atomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title   string     `xml:"title"`
	Links   []atomLink `xml:"link"`
	Content string     `xml:"content"`
	Summary string     `xml:"summary"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// RSSCrawler 通用 RSS/Atom 订阅源爬虫适配器
// 通过 ContentSource.URL 配置订阅地址，支持 RSS 2.0 和 Atom 格式
type RSSCrawler struct {
	client *http.Client
}

func NewRSSCrawler() *RSSCrawler {
	return &RSSCrawler{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (r *RSSCrawler) Name() string {
	return "rss"
}

// Fetch 从 ContentSource.URL 读取 RSS/Atom 订阅源，返回 []RawItem
func (r *RSSCrawler) Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error) {
	url := src.URL
	if url == "" {
		return nil, fmt.Errorf("RSS URL 为空")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HermesCrawler/1.0; +https://github.com/heishanlaoyao0891/web3_study)")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 尝试解析为 RSS 2.0
	var rss rssFeed
	if err := xml.Unmarshal(body, &rss); err == nil && len(rss.Channel.Items) > 0 {
		return r.parseRSSItems(rss.Channel.Items)
	}

	// 尝试解析为 Atom
	var atom atomFeed
	if err := xml.Unmarshal(body, &atom); err == nil && len(atom.Entries) > 0 {
		return r.parseAtomEntries(atom.Entries)
	}

	return nil, fmt.Errorf("无法解析为 RSS 或 Atom 格式: %s", url)
}

func (r *RSSCrawler) parseRSSItems(items []rssItem) ([]RawItem, error) {
	result := make([]RawItem, 0, len(items))
	for _, item := range items {
		if item.Title == "" || item.Link == "" {
			continue
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		rawID := item.Guid
		if rawID == "" {
			rawID = item.Link
		}

		published := time.Time{}
		if item.PubDate != "" {
			if t, err := parseRSSDate(item.PubDate); err == nil {
				published = t
			}
		}

		result = append(result, RawItem{
			Title:       strings.TrimSpace(item.Title),
			URL:         item.Link,
			Content:     strings.TrimSpace(content),
			RawID:       rawID,
			PublishedAt: published,
		})
	}
	return result, nil
}

func (r *RSSCrawler) parseAtomEntries(entries []atomEntry) ([]RawItem, error) {
	result := make([]RawItem, 0, len(entries))
	for _, entry := range entries {
		if entry.Title == "" {
			continue
		}

		// 取链接
		link := ""
		for _, l := range entry.Links {
			if l.Rel == "alternate" || l.Rel == "" {
				link = l.Href
				break
			}
		}
		if link == "" && len(entry.Links) > 0 {
			link = entry.Links[0].Href
		}
		if link == "" {
			continue
		}

		content := entry.Content
		if content == "" {
			content = entry.Summary
		}

		published := time.Time{}
		if entry.Updated != "" {
			if t, err := time.Parse(time.RFC3339, entry.Updated); err == nil {
				published = t
			}
		}

		result = append(result, RawItem{
			Title:       strings.TrimSpace(entry.Title),
			URL:         link,
			Content:     strings.TrimSpace(content),
			RawID:       entry.ID,
			PublishedAt: published,
		})
	}
	return result, nil
}

// parseRSSDate 尝试多种 RSS 日期格式
func parseRSSDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析日期: %s", s)
}

func init() {
	log.Println("[crawler] 注册 RSS 适配器")
	Register(NewRSSCrawler())
}
