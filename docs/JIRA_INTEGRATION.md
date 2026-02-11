# Jira 集成指南

本文档介绍如何配置和使用 Jira 数据收集器。

## 配置说明

### 基本配置

在 `config.yaml` 中配置 Jira 数据源：

```yaml
jira:
  username: "your_jira_username"  # Jira 用户名或邮箱
  url: "https://jira.company.com"  # Jira 服务器地址
  api_token: "${JIRA_API_TOKEN}"   # Jira API Token
  project_key: "PROJ"              # 可选：筛选特定项目
```

### 配置项说明

| 配置项 | 必填 | 说明 | 示例 |
|--------|------|------|------|
| `username` | 是 | Jira 用户名或邮箱 | `"zhangsan"` 或 `"zhangsan@company.com"` |
| `url` | 是 | Jira 服务器地址 | `"https://jira.company.com"` |
| `api_token` | 是 | Jira API Token | 从 Jira 账户设置中获取 |
| `project_key` | 否 | 项目 Key，用于筛选特定项目 | `"PROJ"` |

## 获取 API Token

### Jira Cloud

1. 登录 Jira Cloud：https://id.atlassian.com/manage-profile/security/api-tokens
2. 点击 "Create API token"
3. 输入标签名称（如 "Daily Report"）
4. 复制生成的 Token

### Jira Server / Data Center

1. 登录 Jira 服务器
2. 进入用户设置 → 个人访问令牌
3. 创建新令牌
4. 复制生成的 Token

## 配置示例

### Jira Cloud 配置

```yaml
jira:
  username: "zhangsan@company.com"
  url: "https://yourcompany.atlassian.net"
  api_token: "${JIRA_API_TOKEN}"
  project_key: "PROJ"  # 可选
```

### Jira Server 配置

```yaml
jira:
  username: "zhangsan"
  url: "https://jira.company.com"
  api_token: "${JIRA_API_TOKEN}"
  project_key: "PROJ"  # 可选
```

## 环境变量配置

为了安全，建议使用环境变量存储 API Token：

```bash
export JIRA_API_TOKEN="your_api_token_here"
```

然后在配置文件中引用：

```yaml
jira:
  api_token: "${JIRA_API_TOKEN}"
```

## 工作原理

### 查询逻辑

Jira 收集器使用简化方案：

```sql
(assignee = {username} OR reporter = {username})
AND updated >= '{start}' AND updated <= '{end}'
```

**说明：**
- 查询分配给当前用户或由当前用户创建的任务
- 筛选出在指定时间范围内有更新的任务
- 使用 `updated` 字段判断是否有更新

### 查询参数

- `username`: 用户名或邮箱
- `start`: 开始时间
- `end`: 结束时间
- `project_key`: 项目 Key（可选）

### 输出数据

每个 Jira 任务包含以下信息：

| 字段 | 说明 | 示例 |
|------|------|------|
| `type` | 数据源类型 | `"jira"` |
| `title` | 任务标题 | `"[PROJ-123] 添加新功能"` |
| `time` | 更新时间 | `2026-02-11 16:45:00` |
| `link` | 任务链接 | `https://jira.company.com/browse/PROJ-123` |
| `metadata.issue_key` | 任务 Key | `"PROJ-123"` |
| `metadata.status` | 任务状态 | `"In Progress"` |
| `metadata.assignee` | 分配人 | `"zhangsan"` |
| `metadata.reporter` | 创建人 | `"lisi"` |

## 常见问题

### 1. 认证失败

**可能原因：**
- API Token 无效或过期
- 用户名或邮箱错误
- Jira 服务器地址错误

**解决方法：**
```bash
# 测试 API Token
curl -u "username:token" \
  "https://jira.company.com/rest/api/2/myself"
```

### 2. 没有返回任务

**可能原因：**
- 用户没有任务权限
- 时间范围内没有任务更新
- 项目 Key 配置错误

**解决方法：**
```bash
# 手动查询测试
curl -u "username:token" \
  "https://jira.company.com/rest/api/2/search?jql=assignee=username"
```

### 3. 跨时区问题

Jira 服务器时区可能与本地不同。在 `config.yaml` 中配置时区：

```yaml
time:
  timezone: "Asia/Shanghai"
```

## 高级配置

### 多项目支持

如果需要监控多个项目，可以不配置 `project_key`，或使用多个配置：

```yaml
# 方式 1：不配置 project_key（查询所有项目）
jira:
  username: "zhangsan@company.com"
  url: "https://jira.company.com"
  api_token: "${JIRA_API_TOKEN}"

# 方式 2：指定特定项目
jira:
  username: "zhangsan@company.com"
  url: "https://jira.company.com"
  api_token: "${JIRA_API_TOKEN}"
  project_key: "PROJ"
```

### 自定义 JQL

当前版本使用固定的 JQL 查询。如需自定义查询逻辑，需要修改源代码。

## 性能优化

### 减少返回字段

Jira API 默认返回所有字段，可以配置只返回需要的字段：

```go
// 当前实现默认返回以下字段：
// - key, summary, status, assignee, reporter, updated
```

### 分页处理

Jira API 支持分页，当前版本自动处理分页查询。

## 故障排除

### 调试模式

```bash
./daily_report --config config.yaml --date today
```

查看数据源状态：
- `✅ jira` - 收集成功
- `❌ jira (错误信息)` - 收集失败

### 常见错误信息

| 错误信息 | 原因 | 解决方法 |
|----------|------|----------|
| `401 Unauthorized` | 认证失败 | 检查用户名和 API Token |
| `403 Forbidden` | 权限不足 | 检查用户是否有任务访问权限 |
| `404 Not Found` | URL 错误 | 检查 Jira 服务器地址 |
| `invalid project key` | 项目 Key 无效 | 检查 project_key 配置 |

### API 测试

使用 curl 测试 Jira API：

```bash
# 测试认证
curl -u "username:token" \
  "https://jira.company.com/rest/api/2/myself"

# 测试查询
curl -u "username:token" \
  "https://jira.company.com/rest/api/2/search?jql=assignee=username&fields=key,summary,status"
```

## 示例配置

### 完整配置示例

```yaml
jira:
  username: "zhangsan@company.com"
  url: "https://jira.company.com"
  api_token: "${JIRA_API_TOKEN}"
  project_key: "PROJ"

time:
  timezone: "Asia/Shanghai"
```

### 多环境配置

```bash
# 开发环境
export JIRA_API_TOKEN="dev_token"
export JIRA_URL="https://jira-dev.company.com"

# 生产环境
export JIRA_API_TOKEN="prod_token"
export JIRA_URL="https://jira.company.com"
```

```yaml
jira:
  username: "zhangsan@company.com"
  url: "${JIRA_URL:-https://jira.company.com}"
  api_token: "${JIRA_API_TOKEN}"
```

## 下一步

配置好 Jira 后，可以：
1. 运行 `./daily_report` 生成包含 Jira 任务的日报
2. 使用自定义模板调整任务信息的显示格式

其他数据源集成文档：
- [Git 集成指南](./GIT_INTEGRATION.md)
- [Confluence 集成指南](./CONFLUENCE_INTEGRATION.md)
- [会议平台集成指南](./MEETING_INTEGRATION.md)