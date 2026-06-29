package crawler

import (
	"log/slog"
	"sync"

	"go-blog/model"
	"go-blog/util"

	"github.com/robfig/cron/v3"
)

// Scheduler 抓取调度器
type Scheduler struct {
	cron   *cron.Cron
	mu     sync.Mutex
	jobIDs map[uint]cron.EntryID
}

var scheduler *Scheduler

// GetScheduler 获取全局调度器实例
func GetScheduler() *Scheduler {
	if scheduler == nil {
		scheduler = &Scheduler{
			cron:   cron.New(cron.WithSeconds()),
			jobIDs: make(map[uint]cron.EntryID),
		}
	}
	return scheduler
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.LoadFromDB()
	s.cron.Start()
	slog.Info("爬虫调度器已启动")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	slog.Info("爬虫调度器已停止")
}

// LoadFromDB 从数据库加载所有启用的 ContentSource
func (s *Scheduler) LoadFromDB() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range s.jobIDs {
		s.cron.Remove(id)
	}
	s.jobIDs = make(map[uint]cron.EntryID)

	var sources []model.ContentSource
	util.Db.Where("enabled = ?", true).Find(&sources)

	for _, src := range sources {
		s.addJob(src)
	}

	slog.Info("爬虫调度器已加载抓取源", "count", len(s.jobIDs))
}

// AddSource 添加抓取源
func (s *Scheduler) AddSource(src model.ContentSource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addJob(src)
}

// RemoveSource 移除抓取源
func (s *Scheduler) RemoveSource(sourceID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.jobIDs[sourceID]; ok {
		s.cron.Remove(id)
		delete(s.jobIDs, sourceID)
	}
}

// addJob 注册 cron 任务（不加锁）
func (s *Scheduler) addJob(src model.ContentSource) {
	cronSpec := src.Cron
	if cronSpec == "" {
		cronSpec = "0 * * * *"
	}

	spec := "0 " + cronSpec

	srcCopy := src
	id, err := s.cron.AddFunc(spec, func() {
		slog.Info("抓取任务开始", "source", srcCopy.Name, "type", srcCopy.Type)
		result, err := RunOnce(srcCopy)
		if err != nil {
			slog.Error("抓取任务失败", "source", srcCopy.Name, "error", err)
		} else {
			slog.Info("抓取任务完成",
				"source", srcCopy.Name,
				"fetched", result.FetchedCount,
				"saved", result.SavedCount,
				"duplicates", result.DuplicateCount,
			)
		}
	})

	if err != nil {
		slog.Error("注册抓取任务失败",
			"source", src.Name,
			"cron", spec,
			"error", err,
		)
		return
	}

	s.jobIDs[src.ID] = id
}
