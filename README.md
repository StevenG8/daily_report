# Daily Report Generator

一个用 Go 语言开发的日报生成工具，支持从多个数据源（Git、会议、Jira、Confluence）收集当日工作产出，并生成结构化的 Markdown 报告。

## 功能特性

- ✅ **Git 提交收集** - 支持多仓库扫描，按作者邮箱过滤
- ✅ **模板渲染** - 生成格式化的 Markdown 报告
- ✅ **自定义模板** - 支持使用自定义 Markdown 模板
- ✅ **时间范围灵活** - 支持今天、昨天、自定义日期范围
- ✅ **环境变量支持** - 敏感信息通过环境变量配置
- ✅ **CLI 工具** - 命令行操作，支持文件输出

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone <repository-url>
cd daily_report

# 构建
go build -o daily_report ./cmd/cli

# 验证安装
./daily_report --help
```

## 快速开始

### 1. 创建配置文件

复制示例配置文件：

```bash
cp examples/config.example.yaml config.yaml
```

编辑 `config.yaml`，配置你的 Git 信息：

```yaml
git:
  author: "your.email@example.com"  # 你的 Git 作者（名字或邮箱）
  repos: []  # 可选：指定具体仓库路径
  repo_dirs:
    - "/path/to/your/projects"  # 扫描该目录下的所有 Git 仓库

report:
  mode: "template"

time:
  timezone: "Asia/Shanghai"
```

### 2. 生成日报

```bash
# 生成今天的日报
./daily_report

# 生成昨天的日报
./daily_report --date yesterday

# 生成指定日期的日报
./daily_report --date 2026-02-10

# 生成日期范围的日报
./daily_report --date 2026-02-10,2026-02-12

# 输出到文件
./daily_report --output report.md
```

## 配置说明

### Git 配置

```yaml
git:
  author_email: "your.email@example.com"  # 必填：Git 作者邮箱
  repos:                                  # 可选：指定具体仓库路径
    - "/path/to/repo1"
    - "/path/to/repo2"
  repo_dirs:                              # 可选：扫描目录下的所有仓库
    - "/path/to/projects"
```

### 报告配置

```yaml
report:
  mode: "template"           # 生成模式：template 或 llm
  template_path: ""          # 可选：自定义模板文件路径
  llm:
    provider: "openai"       # LLM 提供商
    model: "gpt-4o"
    api_key: "${LLM_API_KEY}"
```

### 时间配置

```yaml
time:
  timezone: "Asia/Shanghai"  # 时区设置
```

## 自定义模板

### 1. 创建模板文件

创建 `my_template.md`：

```markdown
# 工作日报 - {{date}}

## 💻 代码提交
{{git_count}} 次提交

{{git_section}}

## 📅 会议
{{meeting_count}} 场会议

{{meeting_section}}

---
生成时间: {{generate_time}}
```

### 2. 使用自定义模板

```bash
./daily_report --template my_template.md
```

### 3. 可用的模板变量

| 变量 | 描述 |
|------|------|
| `{{date}}` | 日期（中文格式，如：2026年2月11日） |
| `{{date_en}}` | 日期（英文格式，如：2026-02-11） |
| `{{generate_time}}` | 报告生成时间 |
| `{{source_status}}` | 数据源状态 |
| `{{git_count}}` | Git 提交数量 |
| `{{meeting_count}}` | 会议数量 |
| `{{jira_count}}` | Jira 任务数量 |
| `{{confluence_count}}` | Confluence 文档数量 |
| `{{git_section}}` | Git 提交详情 |
| `{{meeting_section}}` | 会议详情 |
| `{{jira_section}}` | Jira 任务详情 |
| `{{confluence_section}}` | Confluence 文档详情 |

## 命令行参数

```
Usage:
  daily_report [options]

Options:
  -config string
        Path to config file (default "config.yaml")
  -date string
        Date range: today, yesterday, or YYYY-MM-DD,YYYY-MM-DD (default "today")
  -mode string
        Report mode: template or llm (default "template")
  -output string
        Output file path (default: stdout)
  -template string
        Path to custom Markdown template file

Examples:
  daily_report                          # Generate today's report
  daily_report --date yesterday        # Generate yesterday's report
  daily_report --output report.md      # Save to file
  daily_report --template custom.tmpl  # Use custom template
```

## 环境变量

支持在配置文件中使用环境变量：

```yaml
jira:
  api_token: "${JIRA_API_TOKEN}"  # 从环境变量读取

report:
  llm:
    api_key: "${OPENAI_API_KEY}"
```

设置环境变量：

```bash
export JIRA_API_TOKEN="your_token"
export OPENAI_API_KEY="your_key"
```

## 示例输出

```markdown
# 日报 - 2026年2月11日

## 📊 汇总统计
- Git 提交: 5 次
- 会议: 3 场
- Jira 任务: 4 个
- Confluence 文档: 2 篇

## 💻 代码提交

### daily_report
- [feat] 添加用户过滤功能 (14:30)
  commit: abc1234
- [fix] 修复配置加载问题 (16:45)
  commit: def5678

---

生成于: 2026-02-11 18:00:00
数据源状态: ✅ Git
```

## 项目结构

```
daily_report/
├── cmd/cli/           # CLI 入口
├── internal/
│   ├── collector/     # 数据收集器
│   ├── config/        # 配置管理
│   ├── report/        # 报告生成器
│   └── timeutil/      # 时间处理工具
├── pkg/models/        # 数据模型
├── examples/          # 配置和模板示例
└── config.yaml        # 配置文件
```

## 开发

### 运行测试

```bash
go test ./...
```

### 构建

```bash
go build -o daily_report ./cmd/cli
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！