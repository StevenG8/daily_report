## Context
这是一个日报生成工具，需要从多个异构数据源（Git、飞书/钉钉/企业微信、Jira、Confluence）收集当日工作产出，并生成结构化的 Markdown 报告。用户希望通过 CLI 或 Web 界面使用该工具，支持自动时间检测。

## Goals / Non-Goals
- Goals:
  - 提供统一的日报生成能力
  - 支持多种数据源的灵活集成
  - 输出格式化的 Markdown 报告
  - CLI 和 Web 双模式支持
  - 自动检测当日时间范围

- Non-Goals:
  - 实时数据推送（基于定时或手动触发）
  - 复杂的报表分析和可视化
  - 多用户权限管理
  - 历史日报的持久化存储

## Decisions

### 架构设计
采用插件化架构，将数据收集和报告生成分离：

```
daily_report/
├── cmd/
│   ├── cli/          # CLI 入口
│   └── server/       # Web 服务入口
├── internal/
│   ├── collector/    # 数据收集器接口和实现
│   ├── report/       # 报告生成器
│   └── config/       # 配置管理
├── pkg/
│   ├── models/       # 通用数据模型
│   └── sources/      # 具体数据源客户端（git、feishu、jira 等）
└── config/
    └── config.yaml   # 配置文件
```

### 数据收集器接口
定义统一的 Collector 接口，采用方案 B（构造时注入用户配置）：
```go
type Collector interface {
    Name() string
    Collect(ctx context.Context, start, end time.Time) ([]Item, error)
}
```

每个 Collector 在创建时已经持有用户配置，Collect 方法内部实现过滤逻辑。

### 用户过滤策略
核心原则：**今天有我的操作 = 我的产出**

配置采用全部单独配置的方式，无优先级概念：

```yaml
# config.yaml
git:
  author_email: "zhangsan@company.com"
  repos:  # 可选，指定具体仓库
    - "/path/to/specific/repo"
  repo_dirs:  # 可选，扫描目录下的所有仓库
    - "/path/to/projects"

meetings:
  platform: "feishu"  # feishu, dingtalk, wecom
  user_id: "ou_xxx"
  app_id: "${FEISHU_APP_ID}"
  app_secret: "${FEISHU_APP_SECRET}"

jira:
  username: "zhangsan"
  url: "https://jira.company.com"
  api_token: "${JIRA_TOKEN}"
  project_key: "PROJ"  # 可选，筛选特定项目

confluence:
  username: "zhangsan"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_TOKEN}"
  space_key: "SPACE"  # 可选，筛选特定空间

time:
  timezone: "Asia/Shanghai"
```

各数据源的过滤逻辑：

- **Git**: 精确匹配邮箱，收集今天我有提交的记录
  ```bash
  git log --author="<email>" --since="{start}" --until="{end}"
  ```
  支持两种配置方式：
  - `repos`: 指定具体的仓库路径
  - `repo_dirs`: 指定目录，自动扫描该目录下所有包含 .git 的子目录

- **Jira**: 采用简化方案，收集分配给我或我发起、且今天有更新的任务
  ```sql
  (assignee = {username} OR reporter = {username})
  AND updated >= '{start}' AND updated <= '{end}'
  ```

- **Confluence**: 并集，收集今天创建或今天修改的文档（只要其中一次操作人是我）
  ```sql
  (creator = {username} AND created >= '{start}')
  OR (lastModifier = {username} AND lastUpdated >= '{start}')
  ```

- **Meetings**: 收集我参与的会议
  ```http
  GET /events?user_id={user_id}&start={start}&end={end}
  ```

### 时间范围处理
- 默认使用系统时区的当天 00:00:00 到 23:59:59
- CLI 可选支持 --from/--to 参数覆盖默认行为
- Web 界面通过查询参数支持自定义范围

### 报告生成
采用混合模式，支持两种生成方式：

**默认模式：模板渲染**
- 按数据源分类输出（Git 提交、会议、任务、文档）
- 使用 Markdown 模板渲染，数据完全可控
- 支持自定义模板扩展
- 生成速度快，适合日常使用

**可选模式：LLM 美化**
- 通过命令行参数或配置开启
- 使用 LLM 将结构化数据美化成自然语言日报
- 可以总结、归纳、生成工作亮点
- 适合正式汇报场景

配置结构：
```yaml
report:
  mode: "template"  # template | llm
  llm:
    provider: "openai"  # openai, anthropic, local
    model: "gpt-4o"
    api_key: "${LLM_API_KEY}"
    system_prompt: "你是一个专业的日报助手..."
```

命令行参数：
```bash
daily-report              # 默认模板模式
daily-report --mode llm   # LLM 美化模式
```

实现策略：
1. 优先实现模板渲染（核心功能）
2. 预留 LLM 接口设计
3. 后续可选实现 LLM 美化功能

### 配置管理
- 使用 YAML 配置文件管理各数据源的认证信息和配置
- 敏感信息（API Token）支持环境变量覆盖

## Risks / Trade-offs

### 数据源 API 限制
- **Risk**: 飞书、钉钉等平台的 API 可能有调用频率限制
- **Mitigation**: 实现请求缓存和限流控制，支持批量查询优化

### 认证安全
- **Risk**: API Token 等敏感信息需要安全存储
- **Mitigation**: 配置文件中支持加密存储，推荐使用环境变量

### 跨平台会议记录获取
- **Risk**: 不同会议平台的 API 差异较大
- **Mitigation**: 定义抽象的数据模型，各平台适配器转换为统一格式

## Migration Plan
1. 实现核心架构（Collector 接口、基础数据模型）
2. 依次实现各数据源收集器
3. 实现报告生成器
4. 实现 CLI 入口
5. 实现 Web 服务
6. 编写测试和文档

## Open Questions
- Web 界面是否需要用户认证？（初期可考虑仅本地使用）
- 是否需要支持定时任务自动生成？（可预留接口）
- 是否需要支持多语言输出？（初期仅支持中文）