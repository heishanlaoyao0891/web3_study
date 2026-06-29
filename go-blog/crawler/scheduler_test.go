package crawler

import (
	"testing"

	"go-blog/model"

	"github.com/robfig/cron/v3"
)

// ============================================================
// Scheduler 任务管理测试
// ============================================================

func TestSchedulerAddRemoveSource(t *testing.T) {
	s := GetScheduler()

	// 清空已有的 job
	s.mu.Lock()
	for _, id := range s.jobIDs {
		s.cron.Remove(id)
	}
	s.jobIDs = make(map[uint]cron.EntryID)
	s.mu.Unlock()

	// 添加一个源
	src := model.ContentSource{
		ID:      1,
		Name:    "Test",
		Type:    "rss",
		URL:     "https://example.com/feed",
		Cron:    "0 * * * *",
		Enabled: true,
	}
	s.AddSource(src)

	// 检查是否注册
	s.mu.Lock()
	_, exists := s.jobIDs[1]
	s.mu.Unlock()

	if !exists {
		t.Error("AddSource 后 jobIDs[1] 不存在")
	}

	// 移除
	s.RemoveSource(1)

	s.mu.Lock()
	_, exists = s.jobIDs[1]
	s.mu.Unlock()

	if exists {
		t.Error("RemoveSource 后 jobIDs[1] 仍然存在")
	}
}

func TestSchedulerRemoveNonExistent(t *testing.T) {
	s := GetScheduler()
	// 移除不存在的源不应 panic
	s.RemoveSource(99999)
}

func TestSchedulerAddJobWithEmptyCron(t *testing.T) {
	s := GetScheduler()

	src := model.ContentSource{
		ID:      9998,
		Name:    "Test No Cron",
		Type:    "rss",
		Enabled: true,
	}

	// 清空已有
	s.mu.Lock()
	s.jobIDs = make(map[uint]cron.EntryID)
	s.mu.Unlock()

	// 空 Cron 不应 panic
	s.AddSource(src)

	s.mu.Lock()
	_, exists := s.jobIDs[9998]
	s.mu.Unlock()

	if !exists {
		t.Error("空 Cron 的源未注册")
	}

	s.RemoveSource(9998)
}
