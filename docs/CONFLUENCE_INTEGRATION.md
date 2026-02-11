# Confluence 集成指南

本文档介绍如何配置和使用 Confluence 数据收集器。

## 配置说明

### 基本配置

在 `config.yaml` 中配置 Confluence 数据源：

```yaml
confluence:
  username: "your_confluence_username"  # Confluence 用户名或邮箱
  url: "https://confluence.company.com"  # Confluence 服务器地址
  api_token: "${CONFLUENCE_API_TOKEN}"   # Confluence API Token
  space_key: "SPACE"                     # 可选：筛选特定空间
```

### 配置项说明

| 配置项 | 必填 | 说明 | 示例 |
|--------|------|------|------|
| `username` | 是 | Confluence 用户名或邮箱 | `"zhangsan"` 或 `"zhangsan@company.com"` |
| `url` | 是 | Confluence 服务器地址 | `"https://confluence.company.com"` |
| `api_token` | 是 | Confluence API Token | 从 Atlassian 账户获取 |
| `space_key` | 否 | 空间 Key，用于筛选特定空间 | `"DOC"` |

## 获取 API Token

### Confluence Cloud

1. 登录 Atlassian 账户：https://id.atlassian.com/manage-profile/security/api-tokens
2. 点击 "Create API token"
3. 输入标签名称（如 "Daily Report"）
4. 复制生成的 Token

### Confluence Server / Data Center

1. 登录 Confluence 服务器
2. 进入用户设置 → 个人访问令牌
3. 创建新令牌
4. 复制生成的 Token

## 配置示例

### Confluence Cloud 配置

```yaml
confluence:
  username: "zhangsan@company.com"
  url: "https://yourcompany.atlassian.net/wiki"
  api_token: "${CONFLUENCE_API_TOKEN}"
  space_key: "DOC"  # 可选
```

### Confluence Server 配置

```yaml
confluence:
  username: "zhangsan"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_API_TOKEN}"
  space_key: "DOC"  # 可选
```

## 环境变量配置

为了安全，建议使用环境变量存储 API Token：

```bash
export CONFLUENCE_API_TOKEN="your_api_token_here"
```

然后在配置文件中引用：

```yaml
confluence:
  api_token: "${CONFLUENCE_API_TOKEN}"
```

## 工作原理

### 查询逻辑

Confluence 收集器使用并集逻辑：

```sql
(creator = {username} AND created >= '{start}')
OR
(lastModifier = {username} AND lastUpdated >= '{start}')
```

**说明：**
- 查询今天创建的文档（创建人是我）
- 或今天修改过的文档（最后修改人是我）
- 满足任一条件即视为我的产出

### 查询参数

- `username`: 用户名或邮箱
- `start`: 开始时间
- `space_key`: 空间 Key（可选）

### 输出数据

每个 Confluence 文档包含以下信息：

| 字段 | 说明 | 示例 |
|------|------|------|
| `type` | 数据源类型 | `"confluence"` |
| `title` | 文档标题 | `"用户指南"` |
| `time` | 更新时间 | `2026-02-11 10:30:00` |
| `content` | 作者信息 | `"张三"` |
| `link` | 文档链接 | `https://confluence.company.com/pages/123456` |
| `metadata.space_key` | 空间 Key | `"DOC"` |
| `metadata.page_id` | 页面 ID | `"123456"` |
| `metadata.creator` | 创建人 | `"zhangsan"` |
| `metadata.lastModifier` | 最后修改人 | `"zhangsan"` |

## 常见问题

### 1. 认证失败

**可能原因：**
- API Token 无效或过期
- 用户名或邮箱错误
- Confluence 服务器地址错误

**解决方法：**
```bash
# 测试 API Token
curl -u "username:token" \
  "https://confluence.company.com/rest/api/user/current"
```

### 2. 没有返回文档

**可能原因：**
- 用户没有空间访问权限
- 时间范围内没有文档更新
- 空间 Key 配置错误

**解决方法：**
```bash
# 手动查询测试
curl -u "username:token" \
  "https://confluence.company.com/rest/api/search?cql=creator=currentUser()"
```

### 3. 跨时区问题

Confluence 服务器时区可能与本地不同。在 `config.yaml` 中配置时区：

```yaml
time:
  timezone: "Asia/Shanghai"
```

## 高级配置

### 多空间支持

如果需要监控多个空间，可以不配置 `space_key`，或创建多个配置：

```yaml
# 方式 1：不配置 space_key（查询所有空间）
confluence:
  username: "zhangsan@company.com"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_API_TOKEN}"

# 方式 2：指定特定空间
confluence:
  username: "zhangsan@company.com"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_API_TOKEN}"
  space_key: "DOC"
```

### 自定义 CQL

当前版本使用固定的 CQL 查询。如需自定义查询逻辑，需要修改源代码。

## 性能优化

### 减少返回字段

Confluence API 默认返回所有字段，可以配置只返回需要的字段：

```go
// 当前实现默认返回以下字段：
// - title, creator, lastModifier, created, lastUpdated, space
```

### 分页处理

Confluence API 支持分页，当前版本自动处理分页查询。

## 故障排除

### 调试模式

```bash
./daily_report --config config.yaml --date today
```

查看数据源状态：
- `✅ confluence` - 收集成功
- `❌ confluence (错误信息)` - 收集失败

### 常见错误信息

| 错误信息 | 原因 | 解决方法 |
|----------|------|----------|
| `401 Unauthorized` | 认证失败 | 检查用户名和 API Token |
| `403 Forbidden` | 权限不足 | 检查用户是否有空间访问权限 |
| `404 Not Found` | URL 错误 | 检查 Confluence 服务器地址 |
| `invalid space key` | 空间 Key 无效 | 检查 space_key 配置 |

### API 测试

使用 curl 测试 Confluence API：

```bash
# 测试认证
curl -u "username:token" \
  "https://confluence.company.com/rest/api/user/current"

# 测试查询
curl -u "username:token" \
  "https://confluence.company.com/rest/api/search?cql=creator=currentUser()&limit=10"
```

## 示例配置

### 完整配置示例

```yaml
confluence:
  username: "zhangsan@company.com"
  url: "https://confluence.company.com"
  api_token: "${CONFLUENCE_API_TOKEN}"
  space_key: "DOC"

time:
  timezone: "Asia/Shanghai"
```

### 多环境配置

```bash
# 开发环境
export CONFLUENCE_API_TOKEN="dev_token"
export CONFLUENCE_URL="https://confluence-dev.company.com"

# 生产环境
export CONFLUENCE_API_TOKEN="prod_token"
export CONFLUENCE_URL="https://confluence.company.com"
```

```yaml
confluence:
  username: "zhangsan@company.com"
  url: "${CONFLUENCE_URL:-https://confluence.company.com}"
  api_token: "${CONFLUENCE_API_TOKEN}"
```

## 查询逻辑详解

### 场景 1：今天创建的文档

```sql
creator = "zhangsan@company.com" AND created >= "2026-02-11 00:00"
```

返回今天由 `zhangsan@company.com` 创建的所有文档。

### 场景 2：今天修改的文档

```sql
lastModifier = "zhangsan@company.com" AND lastUpdated >= "2026-02-11 00:00"
```

返回今天由 `zhangsan@company.com` 修改的所有文档。

### 场景 3：创建或修改（并集）

```sql
(creator = "zhangsan@company.com" AND created >= "2026-02-11 00:00")
OR
(lastModifier = "zhangsan@company.com" AND lastUpdated >= "2026-02-11 00:00")
```

返回今天创建或修改的所有文档。

## 下一步

配置好 Confluence 后，可以：
1. 运行 `./daily_report` 生成包含 Confluence 文档的日报
2. 使用自定义模板调整文档信息的显示格式

其他数据源集成文档：
- [Git 集成指南](./GIT_INTEGRATION.md)
- [Jira 集成指南](./JIRA_INTEGRATION.md)
- [会议平台集成指南](./MEETING_INTEGRATION.md)