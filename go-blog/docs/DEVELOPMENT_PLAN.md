# 开发计划：从 Web3 社区 → 通用技术学习平台

> 起草：资深开发（架构 + 起步打样） | 协作：中级开发（承接 🟡 卡片）
> 日期：2026-06-29

---

## 一、项目定位变更

- **旧**：Web3 技术社区（Solidity / DeFi / NFT 专属）
- **新**：通用技术学习平台，Web3 仅为一个技术领域标签
- 技术领域：Web3 / Java / Go / Python / AI / ... 管理员可自助配置
- 内容三态：**最新资讯** + **面试题** + **风口话题**
- 资讯/面试题后续由定时任务从热点站抓取入库

## 二、技术约束

- **不前后端分离**，保持 Gin 服务端模板渲染 (templates/*.html)
- 不引入前端框架（React/Vue），不引第三方 UI 库

---

## 三、总体设计决策

### D1. 领域建模：Category 二级树

- 给现有 `Category` 加 `ParentID *uint` + `SortOrder int` + `Icon string`
- 顶层 Category = 技术领域（Web3 / Java / Go / Python / AI …）
- 子 Category = 主题（如 Web3 → Solidity/DeFi；Java → Spring/并发）
- 文章仍挂叶子分类；领域层负责聚合展示
- 管理员复用现有 `/category` 页面增删，最小改动

### D2. 抓取系统：自建轻量框架 + 适配器模式

- 新增 model：`ContentSource`（抓取源配置）/ `CrawlLog`（抓取日志）
- 抽象接口：
  ```go
  type RawItem struct {
      Title       string
      URL         string
      Content     string
      DomainID    uint
      RawID       string    // 源站唯一ID，去重用
      PublishedAt time.Time
  }
  type Crawler interface {
      Name() string
      Fetch(ctx context.Context, src model.ContentSource) ([]RawItem, error)
  }
  ```
- 调度：`github.com/robfig/cron/v3`
- 幂等：Article 加 `Source/SourceURL/RawID/DomainID`，按 `RawID` 去重
- 抓回来的内容进**现有 Article 表**（标记 source），首页"最新资讯"零渲染成本

### D3. 内容三态

| 形态 | 复用/新增 | 说明 |
|---|---|---|
| 最新资讯 | 复用 Article（source 标记） | 抓取入库即上首页 |
| 面试题 | 复用 InterviewQuestion | 已是通用模型，加领域筛选 |
| 风口话题 | 新增 TrendingTopic | 标题/热度/来源/领域/过期 |

### D4. 技术债偿还（最小必要）

| 债务 | 处理 | 理由 |
|---|---|---|
| 全部 handler 重复 `user.(map[string]interface{})` switch 取 ID | 抽 `util.GetUser(c)` | 避免中级开发照抄 |
| admin 判断 `username == "admin"` 字符串硬编码 | 改 `user.Role == "admin"` | 权限应靠角色字段 |
| 配置/密钥明文进 git | 移到环境变量 | 安全红线 |
| 生产环境跑 AutoMigrate | 独立 `scripts/migrate` | 反模式 |
| 无结构化日志 | 引 slog 或 zap | 可观测性 |
| 无健康检查/优雅关停 | 加 `/health` | 部署需要 |

---

## 四、里程碑 & 任务卡片

> 🔴 = 资深（架构 + 起步打样 + CR） | 🟡 = 中级开发独立承接 | 🟢 = 结对

### 🚩 M1 — 通用化基础（领域体系 + 文案去硬编码） | 2-3 天

| # | 任务 | 谁 | 验收 |
|---|---|---|---|
| 1.1 | Category 加 ParentID/SortOrder/Icon，迁移脚本把旧"技术"升为领域 | 🔴 | 字段上线，旧数据兼容 |
| 1.2 | `/category` 管理页支持二级分类 + 父子选择 | 🟡 | admin 能建 "Java → Spring" |
| 1.3 | 新增 SiteConfig(KV) 表，首页标题/副标题/intro/footer 从 DB 读 | 🔴+🟡 | admin 改文案首页即时生效 |
| 1.4 | 首页热门标签从 DB 取 top N | 🟡 | 硬编码 Solidity 等移除 |
| 1.5 | util.GetUser(c) + User.Role 重构 | 🔴 | handler 不再 switch-type |

### 🚩 M2 — 抓取框架 + 首个适配器 | 3-4 天

| # | 任务 | 谁 | 验收 |
|---|---|---|---|
| 2.1 | 引入 robfig/cron + goquery，建 crawler 包骨架 + Crawler 接口 | 🔴 | 框架可单测 |
| 2.2 | ContentSource / CrawlLog model + admin CRUD 页面 | 🟡 | 管理员能增删启停抓取源 |
| 2.3 | Article 加 Source/SourceURL/RawID/DomainID 字段 | 🔴 | 兼容旧数据 |
| 2.4 | 首个适配器：Hacker News（JSON API） | 🟡 | 定时抓取→入库→首页展示 |
| 2.5 | 适配器单测 + 幂等测试 | 🟢 | 连跑两次不重复 |

### 🚩 M3 — 内容三态 + 风口话题 | 2-3 天

| # | 任务 | 谁 | 验收 |
|---|---|---|---|
| 3.1 | TrendingTopic model + admin 页 + 首页专区 | 🟡 | admin 可录入，首页展示 |
| 3.2 | 首页改版：领域 Tab 切换 + 资讯/面试/风口三卡 | 🟡 | 切 Tab 内容过滤 |
| 3.3 | 面试题库加"按领域筛选" | 🟢 | 面试题按领域可筛 |
| 3.4 | 掘金/知乎抓取适配器（goquery DOM） | 🟡 | 至少一源稳定抓取 |

### 🚩 M4 — 工程化收尾 | 2 天

| # | 任务 | 谁 | 验收 |
|---|---|---|---|
| 4.1 | 密钥移到 env，清 .env_prod 历史建议 | 🔴 | 无明文密码 |
| 4.2 | 生产关闭 AutoMigrate，独立 scripts/migrate | 🔴 | main 无 AutoMigrate |
| 4.3 | 引入结构化日志（slog/zap） | 🟡 | 日志带级别字段 |
| 4.4 | /health + 优雅关停 | 🟡 | 部署可用 |
| 4.5 | 核心路径单测 | 🟢 | 覆盖率达标 |

---

## 五、协作机制 & PR 约定

### 分工
- **资深**：定架构、写框架代码（M1.1 / M1.5 / M2.1 / M2.3 / M1.3 起步）、所有 PR CR
- **中级开发**：承接所有 🟡 卡片；卡壳超 1 小时群里 @ 不要闷头死磕

### PR 验收清单（每个 PR 必须满足）
1. 不得在 handler 直接 `util.Db` 做跨表查询 → 封装到 service 函数
2. 新增 handler 复用 `util.GetUser(c)`，**禁止**再写 `user.(map[string]interface{})` switch
3. 新增模板与 `index.html` 风格一致，CSS 复用，不引第三方 UI 库
4. PR 描述带截图 + `go build ./...` 通过

### 分支约定
- 分支命名：`m{里程碑}-{任务号}-{短描述}`，如 `m1-2-category-tree-ui`
- PR 标题：`[M{里程碑}] 任务号 简述`

### 协作节奏
- 每个里程碑结束做一次联调演示
- 不拖到最后一刻集成

---

## 六、环境注意

- `go.mod` 声明 `go 1.25.0`，本机若低于此版本需 `go env -w GOTOOLCHAIN=auto` 或升级
- 网络受限环境下 toolchain 下载可能超时，需配 GOPROXY（如 `https://goproxy.cn`）
- `.env_prod` 里的明文密码在 M4.1 清理前，**不要**把真实 .env 提交到新分支

---

## 七、待确认事项

1. M2 抓取来源清单是否只做 Hacker News，或有其他指定站点
2. 是否在正式开发前先清理 .env_prod 密钥泄漏（建议是）