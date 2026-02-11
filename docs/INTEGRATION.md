# 数据源集成指南索引

本文档提供了 Daily Report Generator 支持的所有数据源的集成指南。

## 支持的数据源

| 数据源 | 状态 | 集成文档 |
|--------|------|----------|
| Git | ✅ 已实现 | [Git 集成指南](./GIT_INTEGRATION.md) |
| 飞书会议 | ⏳ 计划中 | [会议平台集成指南](./MEETING_INTEGRATION.md) |
| 钉钉会议 | ⏳ 计划中 | [会议平台集成指南](./MEETING_INTEGRATION.md) |
| 企业微信会议 | ⏳ 计划中 | [会议平台集成指南](./MEETING_INTEGRATION.md) |
| Jira | ⏳ 计划中 | [Jira 集成指南](./JIRA_INTEGRATION.md) |
| Confluence | ⏳ 计划中 | [Confluence 集成指南](./CONFLUENCE_INTEGRATION.md) |

## 快速开始

### 1. Git 数据源（已实现）

Git 数据源已完全实现，配置简单：

```yaml
git:
  author_email: "your.email@example.com"
  repo_dirs:
    - "/path/to/projects"
```

详细配置请参考：[Git 集成指南](./GIT_INTEGRATION.md)

### 2. 其他数据源（计划中）

其他数据源的集成框架已经设计完成，可以按照各集成指南进行配置。

## 配置文件结构

完整的 `config.yaml` 示例：

```yaml
# Git 配置（已实现）
git:
  author_email: "your.email@example.com"
  repo_dirs:
    - "/path/to/projects"

# 会议配置（计划中）
meetings:
  platform: "feishu"
  user_id: "your_user_id"
  app_id: "${FEISHU_APP_ID}"
  app_secret: "${FEISHU_APP_SECRET}"

# Jira 配置（计划中）
jira:
  username: "your_jira_username"
  url: "https://jira.company.com"
  api_token: "${JIRA_API_TOKEN}"
  project_key: "PROJ"

# Confluence 配置（计划中）
confluence:
  username: "your_confluence_username"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_API_TOKEN}"
  space_key: "SPACE"

# 报告配置
report:
  mode: "template"
  template_path: ""

# 时间配置
time:
  timezone: "Asia/Shanghai"
```

## 环境变量

为了安全，建议使用环境变量存储敏感信息：

```bash
# 会议平台
export FEISHU_APP_ID="your_app_id"
export FEISHU_APP_SECRET="your_app_secret"

# Jira
export JIRA_API_TOKEN="your_jira_token"

# Confluence
export CONFLUENCE_API_TOKEN="your_confluence_token"
```

## 通用配置

### 时间配置

所有数据源都使用统一的时间配置：

```yaml
time:
  timezone: "Asia/Shanghai"
```

支持的时区格式：
- `"Asia/Shanghai"` - 上海时间
- `"UTC"` - 世界标准时间
- `"America/New_York"` - 纽约时间

### 时间范围

工具支持多种时间范围指定方式：

```bash
# 今天
./daily_report

# 昨天
./daily_report --date yesterday

# 特定日期
./daily_report --date 2026-02-10

# 日期范围
./daily_report --date 2026-02-10,2026-02-12
```

## 故障排除

### 通用问题

**问题：数据源状态显示失败**

解决方法：
1. 检查配置文件路径是否正确
2. 查看具体的错误信息
3. 根据对应集成指南进行排查

**问题：没有返回数据**

解决方法：
1. 检查用户标识是否正确
2. 确认时间范围内是否有数据
3. 检查权限配置

### 调试模式

```bash
# 显示详细输出
./daily_report --config config.yaml --date today
```

工具会在报告末尾显示各数据源的状态：
- `✅ 数据源名称` - 收集成功
- `❌ 数据源名称 (错误信息)` - 收集失败

## 下一步

1. 配置数据源
2. 运行 `./daily_report` 生成日报
3. 根据需要自定义模板
4. 定期运行（可配合 cron 定时任务）

## 相关文档

- [README.md](../README.md) - 主文档
- [配置示例](../examples/config.example.yaml) - 配置文件示例
- [模板示例](../examples/template.example.md) - 自定义模板示例

## 贡献

如需添加新的数据源支持：
1. 实现 Collector 接口
2. 添加配置项
3. 编写集成文档
4. 添加测试

参考现有的 [Git 收集器实现](../internal/collector/git.go)。