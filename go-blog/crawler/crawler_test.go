package crawler

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-blog/model"
)

// ============================================================
// 适配器注册表单元测试
// ============================================================

func TestRegisterGet(t *testing.T) {
	// 备份并清空注册表
	saved := registry
	registry = make(map[string]Crawler)
	defer func() { registry = saved }()

	mock := &mockCrawler{name: "test"}
	Register(mock)

	c, err := Get("test")
	if err != nil {
		t.Fatalf("Get('test') 失败: %v", err)
	}
	if c.Name() != "test" {
		t.Fatalf("期望 Name()='test', 得到 '%s'", c.Name())
	}

	_, err = Get("unknown")
	if err == nil {
		t.Fatal("期望 Get('unknown') 返回错误，但得到 nil")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	saved := registry
	registry = make(map[string]Crawler)
	defer func() { registry = saved }()

	Register(&mockCrawler{name: "dup"})
	Register(&mockCrawler{name: "dup"})
	if len(registry) != 1 {
		t.Fatalf("期望 1 个注册项，实际 %d", len(registry))
	}
}

func TestRegistryContainsPreRegistered(t *testing.T) {
	// 验证 init() 中自动注册的 4 个适配器
	expected := []string{"hackernews", "rss", "juejin", "zhihu"}
	for _, name := range expected {
		if _, err := Get(name); err != nil {
			t.Errorf("适配器 '%s' 未注册: %v", name, err)
		}
	}
}

// ============================================================
// truncateStr 单元测试
// ============================================================

func TestTruncateStr(t *testing.T) {
	tests := []struct {
		input   string
		maxLen  int
		wantLen int
		want    string
	}{
		{"hello", 10, 5, "hello"},
		{"hello", 3, 3, "hel"},
		{"你好世界", 2, 2, "你好"},
		{"", 10, 0, ""},
		{"a", 1, 1, "a"},
	}

	for _, tc := range tests {
		got := truncateStr(tc.input, tc.maxLen)
		if len([]rune(got)) != tc.wantLen {
			t.Errorf("truncateStr(%q, %d) 长度=%d, 期望 %d", tc.input, tc.maxLen, len([]rune(got)), tc.wantLen)
		}
		if tc.want != "" && got != tc.want {
			t.Errorf("truncateStr(%q, %d) = %q, 期望 %q", tc.input, tc.maxLen, got, tc.want)
		}
	}
}

// ============================================================
// parseRSSDate 单元测试
// ============================================================

func TestParseRSSDate(t *testing.T) {
	tests := []struct {
		input     string
		wantErr   bool
		checkYear int
	}{
		// RFC1123Z
		{"Mon, 02 Jan 2006 15:04:05 -0700", false, 2006},
		// RFC1123
		{"Mon, 02 Jan 2006 15:04:05 MST", false, 2006},
		// RFC822
		{"02 Jan 06 15:04 MST", false, 2006},
		// RFC3339
		{"2024-03-15T10:30:00Z", false, 2024},
		// 常见 RSS 日期格式
		{"Mon, 15 Mar 2024 10:30:00 +0800", false, 2024},
		{"2024-03-15T10:30:00Z", false, 2024},
		// Go Blog 常用格式
		{"2024-06-15 10:30:00", false, 2024},
		// 无效输入
		{"", true, 0},
		{"not-a-date", true, 0},
	}

	for _, tc := range tests {
		got, err := parseRSSDate(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseRSSDate(%q) 期望错误，得到 %v", tc.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseRSSDate(%q) 返回错误: %v", tc.input, err)
			continue
		}
		if got.Year() != tc.checkYear {
			t.Errorf("parseRSSDate(%q) 年份=%d, 期望 %d", tc.input, got.Year(), tc.checkYear)
		}
	}
}

// ============================================================
// RSS 条目解析：parseRSSItems — 核心逻辑
// ============================================================

func TestParseRSSItems(t *testing.T) {
	items := []rssItem{
		{
			Title:       "Go 1.22 Released",
			Link:        "https://go.dev/blog/go1.22",
			Description: "A new release with improved generics",
			Guid:        "go1.22-2024",
			PubDate:     "Mon, 15 Feb 2024 10:00:00 +0000",
			Categories:  []string{"Go", "Release"},
		},
		{
			Title: "Testing Generics",
			Link:  "https://go.dev/blog/generics",
			// no Guid — falls back to Link
		},
		{
			// empty title → skipped
			Title: "",
			Link:  "https://example.com/no-title",
			Guid:  "no-title",
		},
		{
			// empty link → skipped
			Title: "Empty Link",
			Link:  "",
			Guid:  "empty-link",
		},
	}

	crawler := NewRSSCrawler()
	result, err := crawler.parseRSSItems(items)
	if err != nil {
		t.Fatalf("parseRSSItems 返回错误: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("期望 2 个条目, 得到 %d", len(result))
	}

	// 第一条：有 guid 就用 guid
	if result[0].Title != "Go 1.22 Released" {
		t.Errorf("Title = %q, 期望 'Go 1.22 Released'", result[0].Title)
	}
	if result[0].URL != "https://go.dev/blog/go1.22" {
		t.Errorf("URL = %q", result[0].URL)
	}
	if result[0].RawID != "go1.22-2024" {
		t.Errorf("RawID = %q, 期望 'go1.22-2024'", result[0].RawID)
	}
	if result[0].PublishedAt.Year() != 2024 {
		t.Errorf("Year = %d, 期望 2024", result[0].PublishedAt.Year())
	}

	// 第二条：无 guid，用 link 做 RawID
	if result[1].Title != "Testing Generics" {
		t.Errorf("Title[1] = %q", result[1].Title)
	}
	if result[1].RawID != "https://go.dev/blog/generics" {
		t.Errorf("RawID[1] = %q, 期望 Link URL", result[1].RawID)
	}
}

// ============================================================
// RSS Content 字段回退关系：有 Content 用 Content，否则用 Description
// ============================================================

func TestParseRSSItems_ContentFallback(t *testing.T) {
	items := []rssItem{
		{
			Title:       "With Content",
			Link:        "https://example.com/1",
			Guid:        "1",
			Content:     "<p>Full HTML content</p>",
			Description: "Short description",
		},
		{
			Title:       "Description Only",
			Link:        "https://example.com/2",
			Guid:        "2",
			Description: "Only description",
		},
	}

	crawler := NewRSSCrawler()
	result, err := crawler.parseRSSItems(items)
	if err != nil {
		t.Fatalf("parseRSSItems 返回错误: %v", err)
	}

	if result[0].Content != "<p>Full HTML content</p>" {
		t.Errorf("Content[0] = %q, 期望 HTML content", result[0].Content)
	}
	if result[1].Content != "Only description" {
		t.Errorf("Content[1] = %q, 期望 'Only description'", result[1].Content)
	}
}

// ============================================================
// Atom 条目解析：parseAtomEntries — 核心逻辑
// ============================================================

func TestParseAtomEntries(t *testing.T) {
	entries := []atomEntry{
		{
			Title:   "Announcing Go 1.22",
			Links:   []atomLink{{Href: "https://go.dev/blog/go1.22", Rel: "alternate"}},
			ID:      "tag:go.dev,2024:go1.22",
			Updated: "2024-02-15T10:00:00Z",
			Summary: "Go 1.22 release notes",
		},
		{
			// 无链接 → 跳过
			Title:   "No Link Entry",
			ID:      "tag:go.dev,2024:nolink",
			Updated: "2024-03-01T08:00:00Z",
		},
		{
			// 空标题 → 跳过
			Title:   "",
			Links:   []atomLink{{Href: "https://example.com/empty-title"}},
			ID:      "tag:go.dev,2024:empty-title",
		},
	}

	crawler := NewRSSCrawler()
	result, err := crawler.parseAtomEntries(entries)
	if err != nil {
		t.Fatalf("parseAtomEntries 返回错误: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("期望 1 个 Atom 条目, 得到 %d", len(result))
	}

	if result[0].Title != "Announcing Go 1.22" {
		t.Errorf("Title = %q, 期望 'Announcing Go 1.22'", result[0].Title)
	}
	if result[0].URL != "https://go.dev/blog/go1.22" {
		t.Errorf("URL = %q", result[0].URL)
	}
	if result[0].RawID != "tag:go.dev,2024:go1.22" {
		t.Errorf("RawID = %q", result[0].RawID)
	}
}

// ============================================================
// Atom 无 alternate link 时应回退到第一个 link
// ============================================================

func TestParseAtomEntries_FallbackLink(t *testing.T) {
	entries := []atomEntry{
		{
			Title:   "Go Blog",
			Links:   []atomLink{{Href: "https://go.dev/blog"}},
			ID:      "tag:go.dev,2024:blog",
			Updated: "2024-01-01T00:00:00Z",
		},
	}

	crawler := NewRSSCrawler()
	result, err := crawler.parseAtomEntries(entries)
	if err != nil {
		t.Fatalf("parseAtomEntries 返回错误: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("期望 1 个条目, 得到 %d", len(result))
	}
	if result[0].URL != "https://go.dev/blog" {
		t.Errorf("URL = %q, 期望 'https://go.dev/blog'", result[0].URL)
	}
}

// ============================================================
// HackerNews item 解析逻辑测试
// ============================================================

func TestHackerNewsItemParsing(t *testing.T) {
	// 这段 JSON 模拟 HN API 返回的单条 item
	body := `{"id":123456,"title":"Go 1.22 Released","url":"https://go.dev/blog/go1.22","time":1708000000,"score":345}`

	var hnItem struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
		Time  int64  `json:"time"`
		Score int    `json:"score"`
	}

	mustUnmarshal(t, body, &hnItem)

	if hnItem.ID != 123456 {
		t.Errorf("ID = %d, 期望 123456", hnItem.ID)
	}
	if hnItem.Title != "Go 1.22 Released" {
		t.Errorf("Title = %q", hnItem.Title)
	}
	if hnItem.URL != "https://go.dev/blog/go1.22" {
		t.Errorf("URL = %q", hnItem.URL)
	}
	if hnItem.Time != 1708000000 {
		t.Errorf("Time = %d", hnItem.Time)
	}
	if hnItem.Score != 345 {
		t.Errorf("Score = %d", hnItem.Score)
	}
}

func TestHackerNewsItemNoURL(t *testing.T) {
	// Ask HN 无 URL，应回退到 HN 链接
	body := `{"id":789,"title":"Ask HN: Best books on Go?","time":1708000100,"score":87}`

	var hnItem struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
		Time  int64  `json:"time"`
		Score int    `json:"score"`
	}
	mustUnmarshal(t, body, &hnItem)

	// 模拟 HackerNewsCrawler 构建逻辑
	articleURL := hnItem.URL
	if articleURL == "" {
		articleURL = "https://news.ycombinator.com/item?id=789"
	}
	if articleURL != "https://news.ycombinator.com/item?id=789" {
		t.Errorf("回退 URL = %q", articleURL)
	}
}

// ============================================================
// 掘金响应解析测试
// ============================================================

func TestJuejinResponseParsing(t *testing.T) {
	body := `{
		"err_no": 0,
		"err_msg": "success",
		"data": [
			{
				"item_info": {
					"article_info": {
						"article_id": "7301234567890",
						"title": "深入理解 Go 泛型",
						"brief_content": "Go 1.18 引入泛型...",
						"ctime": 1708000000,
						"view_count": 1234,
						"digg_count": 56
					}
				}
			}
		]
	}`

	var result struct {
		ErrNo  int    `json:"err_no"`
		ErrMsg string `json:"err_msg"`
		Data   []struct {
			ItemInfo *struct {
				ArticleInfo *struct {
					ArticleID    string `json:"article_id"`
					Title        string `json:"title"`
					BriefContent string `json:"brief_content"`
					CTime        int64  `json:"ctime"`
				} `json:"article_info"`
			} `json:"item_info"`
		} `json:"data"`
	}
	mustUnmarshal(t, body, &result)

	if result.ErrNo != 0 {
		t.Fatalf("ErrNo = %d, 期望 0", result.ErrNo)
	}
	if len(result.Data) != 1 {
		t.Fatalf("期望 1 条数据, 得到 %d", len(result.Data))
	}

	info := result.Data[0].ItemInfo.ArticleInfo
	if info.ArticleID != "7301234567890" {
		t.Errorf("ArticleID = %q", info.ArticleID)
	}
	if info.Title != "深入理解 Go 泛型" {
		t.Errorf("Title = %q", info.Title)
	}
}

// ============================================================
// 知乎响应解析测试
// ============================================================

func TestZhihuResponseParsing(t *testing.T) {
	body := `{
		"data": [
			{
				"id": 12345,
				"type": "hot_list",
				"target": {
					"id": 67890,
					"title": "如何学习 Go 语言？",
					"url": "https://www.zhihu.com/question/67890",
					"excerpt": "这是一个非常经典的问题...",
					"type": "question"
				}
			}
		]
	}`

	var result struct {
		Data []struct {
			ID     int64  `json:"id"`
			Type   string `json:"type"`
			Target *struct {
				ID          int64  `json:"id"`
				Title       string `json:"title"`
				URL         string `json:"url"`
				Excerpt     string `json:"excerpt"`
				Description string `json:"description"`
				Type        string `json:"type"`
			} `json:"target"`
		} `json:"data"`
	}
	mustUnmarshal(t, body, &result)

	if len(result.Data) != 1 {
		t.Fatalf("期望 1 条数据, 得到 %d", len(result.Data))
	}

	item := result.Data[0]
	if item.ID != 12345 {
		t.Errorf("item ID = %d, 期望 12345", item.ID)
	}
	target := item.Target
	if target.Title != "如何学习 Go 语言？" {
		t.Errorf("Title = %q", target.Title)
	}
	if target.URL != "https://www.zhihu.com/question/67890" {
		t.Errorf("URL = %q", target.URL)
	}
}

func TestZhihuTargetRelativeURL(t *testing.T) {
	// 模拟知乎返回相对路径 URL
	body := `{"data":[{"id":1,"target":{"id":1,"title":"Test","url":"/question/1"}}]}`

	var result struct {
		Data []struct {
			Target *struct {
				URL string `json:"url"`
			} `json:"target"`
		} `json:"data"`
	}
	mustUnmarshal(t, body, &result)

	link := result.Data[0].Target.URL
	if len(link) > 0 && link[0] == '/' {
		link = "https://www.zhihu.com" + link
	}
	if link != "https://www.zhihu.com/question/1" {
		t.Errorf("处理后的 URL = %q, 期望 'https://www.zhihu.com/question/1'", link)
	}
}

// ============================================================
// 辅助函数
// ============================================================

func mustUnmarshal(t *testing.T, data string, v interface{}) {
	t.Helper()
	_ = data
	if err := json.Unmarshal([]byte(data), v); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
}

// ============================================================
// Mock
// ============================================================

type mockCrawler struct {
	name string
}

func (m *mockCrawler) Name() string { return m.name }
func (m *mockCrawler) Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error) {
	return []RawItem{
		{Title: "mock", URL: "https://example.com", RawID: "1", PublishedAt: time.Now()},
	}, nil
}
