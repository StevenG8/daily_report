# Change: 添加日报生成工具

## Why
开发者和团队成员需要高效地汇总每日工作产出，包括代码提交、会议记录、任务变更和文档编辑。目前这些信息分散在多个平台（Git、飞书/钉钉/企业微信、Jira、Confluence），手动整理耗时且容易遗漏。

## What Changes
- 创建一个 Go 语言的日报生成工具，支持 CLI 和 Web 两种使用方式
- 集成多个数据源收集器：Git、会议平台、Jira、Confluence
- 自动检测当天时间范围，支持 Markdown 格式输出
- 提供可扩展的数据源插件架构

## Impact
- Affected specs: data-collection, report-generation
- Affected code: 创建新的 cmd/、internal/、pkg/ 目录结构
- 外部依赖: 需要各数据源的 API 客户端库