package crawler

import (
	"log"
	"sync"

	"go-blog/model"
	"go-blog/util"

	"github.com/robfig/cron/v3"
)

// Scheduler 抓取调度器
// 从 DB 读取 ContentSource 配置，按 cron 表达式定时执行
type Scheduler struct {
	cron   *cron.Cron
	mu     sync.Mutex
	jobIDs map[uint]cron.EntryID // sourceID -> cron job id
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

// Start 启动调度器，从 DB 加载所有启用的抓取源
func (s *Scheduler) Start() {
	s.LoadFromDB()
	s.cron.Start()
	log.Println("[crawler] 调度器已启动")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("[crawler] 调度器已停止")
}

// LoadFromDB 从数据库加载所有启用的 ContentSource，注册定时任务
func (s *Scheduler) LoadFromDB() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 清除旧任务
	for _, id := range s.jobIDs {
		s.cron.Remove(id)
	}
	s.jobIDs = make(map[uint]cron.EntryID)

	var sources []model.ContentSource
	util.Db.Where("enabled = ?", true).Find(&sources)

	for _, src := range sources {
		s.addJob(src)
	}

	log.Printf("[crawler] 已加载 %d 个抓取源", len(s.jobIDs))
}

// AddSource 添加单个抓取源到调度器
func (s *Scheduler) AddSource(src model.ContentSource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addJob(src)
}

// RemoveSource 从调度器移除抓取源
func (s *Scheduler) RemoveSource(sourceID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.jobIDs[sourceID]; ok {
		s.cron.Remove(id)
		delete(s.jobIDs, sourceID)
	}
}

// addJob 注册一个 cron 任务（不加锁，调用方负责）
func (s *Scheduler) addJob(src model.ContentSource) {
	cronSpec := src.Cron
	if cronSpec == "" {
		cronSpec = "0 * * * *" // 默认每小时
	}

	// robfig/cron v3 格式：秒 分 时 日 月 周
	// ContentSource.Cron 存的是 5 段格式，补一个秒位
	spec := "0 " + cronSpec

	srcCopy := src
	id, err := s.cron.AddFunc(spec, func() {
		log.Printf("[crawler] 开始抓取: %s (type=%s)", srcCopy.Name, srcCopy.Type)
		result, err := RunOnce(srcCopy)
		if err != nil {
			log.Printf("[crawler] 抓取失败: %s - %v", srcCopy.Name, err)
		} else {
			log.Printf("[crawler] 抓取完成: %s - 抓取%d条, 入库%d条, 重复%d条",
				srcCopy.Name, result.FetchedCount, result.SavedCount, result.DuplicateCount)
		}
	})

	if err != nil {
		log.Printf("[crawler] 注册任务失败: %s (cron=%s) - %v", src.Name, spec, err)
		return
	}

	s.jobIDs[src.ID] = id
}