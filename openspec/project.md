# Project Context

## Purpose
一键生成日报，提升工作效率，减少重复劳动。

## Tech Stack
- Go（模块名：`daily_report`）
- CLI：标准库 `flag`
- 配置：`gopkg.in/yaml.v3` + `${ENV_VAR}` 展开
- 数据采集：本地 Git CLI（`git log`）
- 输出：Markdown（内置模板 + 自定义模板）

## Project Conventions

### Code Style
- 遵循 `gofmt` 与 Go idiomatic 风格
- 导出类型/函数保持清晰注释
- 错误优先返回并保留上下文（`fmt.Errorf(... %w ...)`）
- 业务实现放在 `internal/`，通用数据模型放在 `pkg/models`

### Architecture Patterns
- 分层结构：
  - `cmd/cli/main.go`：参数解析、组装依赖、执行流程
  - `internal/config`：配置加载与环境变量展开
  - `internal/timeutil`：日期范围解析与时区处理
  - `internal/collector`：`Collector` 接口、`MultiCollector` 聚合器、Git 实现
  - `internal/report`：模板渲染器
  - `pkg/models`：`Item`、`ReportData`、`SourceStatus` 等模型
- 扩展模式：新增数据源时优先实现 `Collector` 接口并接入 `MultiCollector`
- 当前实际状态：仅 Git Collector 落地；会议/Jira/Confluence 为规划能力

### Testing Strategy
- 使用 Go 原生测试框架（`go test ./...`）
- 已覆盖模块：
  - `internal/config`（配置解析、必填校验、环境变量）
  - `internal/timeutil`（today/yesterday/单日/区间解析）
  - `internal/collector`（Git 仓库扫描与提交解析）
  - `internal/report`（默认模板与自定义模板渲染）
- 新增功能要求：
  - 至少补充对应包的单元测试
  - 变更 CLI 参数或报告格式时需更新 README/示例

### Git Workflow
- 当前仓库已在 Git 管理中
- 暂无强制分支命名或 Commit 规范（建议使用清晰的动词开头提交信息）
- 对行为变更，优先保证“代码 + 测试 + 文档”同一变更集提交

## Domain Context
- 目标用户：需要快速汇总“当天工作产出”的开发者/团队成员
- 产出聚焦“我在当天的活动”，不是全局统计报表
- 默认日报维度：
  - Git 提交
  - 会议
  - Jira 任务
  - Confluence 文档
- 当前实现只保证 Git 数据可用，其他数据源在 OpenSpec 历史变更中定义但未落地

## Important Constraints
- `git.author` 为必填配置；缺失应立即报错退出
- 时区依赖 `time.timezone`（默认 `Asia/Shanghai`）
- 配置中的敏感信息必须支持环境变量注入，不应硬编码入仓库
- 工具定位为本地/离线优先的 CLI 体验，避免引入复杂运行时依赖
- 报告生成默认模板渲染；LLM 模式目前仅保留配置与接口位置

## External Dependencies
- 本地 Git 可执行程序（`git` 命令可用）
- Go 依赖：
  - `gopkg.in/yaml.v3`
- 预留外部系统（未实现）：
  - 飞书/钉钉/企业微信会议 API
  - Jira API
  - Confluence API
  - LLM Provider（OpenAI/Anthropic/Local）
