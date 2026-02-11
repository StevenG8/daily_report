# 会议平台集成指南

本文档介绍如何配置和使用会议平台数据收集器（飞书、钉钉、企业微信）。

## 支持的平台

- 飞书 (Feishu)
- 钉钉 (DingTalk)
- 企业微信 (WeCom)

## 配置说明

### 通用配置

在 `config.yaml` 中配置会议平台：

```yaml
meetings:
  platform: "feishu"  # 平台类型：feishu, dingtalk, wecom
  user_id: "your_user_id"
  app_id: "${MEETING_APP_ID}"
  app_secret: "${MEETING_APP_SECRET}"
```

### 配置项说明

| 配置项 | 必填 | 说明 | 示例 |
|--------|------|------|------|
| `platform` | 是 | 平台类型 | `"feishu"`, `"dingtalk"`, `"wecom"` |
| `user_id` | 是 | 用户标识 | 飞书用户ID、钉钉unionid、企业微信userid |
| `app_id` | 是 | 应用ID | 飞书App ID、钉钉AppKey、企业微信AgentId |
| `app_secret` | 是 | 应用密钥 | 飞书App Secret、钉钉AppSecret、企业微信Secret |

## 飞书集成

### 获取凭证

1. 登录飞书开放平台：https://open.feishu.cn/
2. 创建企业自建应用
3. 获取 `App ID` 和 `App Secret`
4. 在权限管理中开启相关权限：
   - `calendar:calendar:readonly` - 读取日历权限
5. 获取用户 ID

### 配置示例

```yaml
meetings:
  platform: "feishu"
  user_id: "ou_xxxxxxxxxxxxxx"  # 飞书用户ID
  app_id: "${FEISHU_APP_ID}"
  app_secret: "${FEISHU_APP_SECRET}"
```

### 获取用户ID

```bash
# 通过飞书开放平台 API 查询
curl -X GET "https://open.feishu.cn/open-apis/contact/v3/users/me" \
  -H "Authorization: Bearer <access_token>"
```

## 钉钉集成

### 获取凭证

1. 登录钉钉开放平台：https://open.dingtalk.com/
2. 创建企业内部应用
3. 获取 `AppKey` 和 `AppSecret`
4. 在权限管理中开启相关权限：
   - `日历: 读取日历信息`
5. 获取用户 unionid

### 配置示例

```yaml
meetings:
  platform: "dingtalk"
  user_id: "xxxxxxxxx"  # 钉钉用户unionid
  app_id: "${DINGTALK_APP_KEY}"
  app_secret: "${DINGTALK_APP_SECRET}"
```

### 获取用户unionid

```bash
# 通过钉钉 API 查询
curl -X GET "https://api.dingtalk.com/v1.0/contact/users/me" \
  -H "x-acs-dingtalk-access-token: <access_token>"
```

## 企业微信集成

### 获取凭证

1. 登录企业微信管理后台：https://work.weixin.qq.com/
2. 创建应用
3. 获取 `AgentId` 和 `Secret`
4. 在权限管理中开启相关权限：
   - `查看日历`
5. 获取用户 userid

### 配置示例

```yaml
meetings:
  platform: "wecom"
  user_id: "zhangsan"  # 企业微信userid
  app_id: "${WECOM_AGENT_ID}"
  app_secret: "${WECOM_SECRET}"
```

### 获取用户userid

```bash
# 通过企业微信 API 查询
curl "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=<access_token>&code=<code>"
```

## 工作原理

### API 调用流程

1. 使用 `app_id` 和 `app_secret` 获取访问令牌（Access Token）
2. 使用访问令牌调用日历 API
3. 查询用户在指定时间范围内的会议
4. 提取会议主题、时间、参会人等信息

### 查询参数

- `start_time`: 开始时间
- `end_time`: 结束时间
- `user_id`: 用户标识

### 输出数据

每个会议包含以下信息：

| 字段 | 说明 | 示例 |
|------|------|------|
| `type` | 数据源类型 | `"meeting"` |
| `title` | 会议主题 | `"技术方案评审"` |
| `time` | 会议时间 | `2026-02-11 14:00:00` |
| `content` | 参会人员 | `"张三、李四、王五"` |
| `link` | 会议链接 | `https://feishu.cn/meeting/xxx` |

## 环境变量配置

为了安全，建议使用环境变量存储敏感信息：

```bash
# 飞书
export FEISHU_APP_ID="cli_xxxxxxxxxxxxxx"
export FEISHU_APP_SECRET="your_secret_here"

# 钉钉
export DINGTALK_APP_KEY="dingxxxxxxxxxxxx"
export DINGTALK_APP_SECRET="your_secret_here"

# 企业微信
export WECOM_AGENT_ID="1000001"
export WECOM_SECRET="your_secret_here"
```

然后在配置文件中引用：

```yaml
meetings:
  platform: "feishu"
  app_id: "${FEISHU_APP_ID}"
  app_secret: "${FEISHU_APP_SECRET}"
```

## 常见问题

### 1. 获取会议失败

**可能原因：**
- 应用权限不足
- Access Token 过期
- 用户 ID 错误

**解决方法：**
- 检查应用权限配置
- 重新获取 Access Token
- 验证用户 ID 是否正确

### 2. 权限不足

**解决方法：**
1. 登录对应开放平台
2. 进入应用管理
3. 在权限管理中开启日历相关权限
4. 提交审核并等待通过

### 3. 用户 ID 获取

**方法 1：通过开放平台 API**
```bash
# 飞书
curl -X GET "https://open.feishu.cn/open-apis/contact/v3/users/me" \
  -H "Authorization: Bearer <access_token>"

# 钉钉
curl -X GET "https://api.dingtalk.com/v1.0/contact/users/me" \
  -H "x-acs-dingtalk-access-token: <access_token>"
```

**方法 2：通过企业管理后台**
在企业管理后台的成员管理页面查看用户 ID

## API 限流

各平台都有 API 调用频率限制：

| 平台 | 限制说明 |
|------|----------|
| 飞书 | 根据应用类型不同，通常 1000 次/分钟 |
| 钉钉 | 通常 1000 次/分钟 |
| 企业微信 | 通常 500 次/分钟 |

如遇到限流，工具会自动重试，超过阈值则跳过该数据源。

## 示例配置

### 飞书配置

```yaml
meetings:
  platform: "feishu"
  user_id: "ou_1234567890abcdef"
  app_id: "${FEISHU_APP_ID}"
  app_secret: "${FEISHU_APP_SECRET}"
```

### 钉钉配置

```yaml
meetings:
  platform: "dingtalk"
  user_id: "ding1234567890abcdef"
  app_id: "${DINGTALK_APP_KEY}"
  app_secret: "${DINGTALK_APP_SECRET}"
```

### 企业微信配置

```yaml
meetings:
  platform: "wecom"
  user_id: "zhangsan"
  app_id: "${WECOM_AGENT_ID}"
  app_secret: "${WECOM_SECRET}"
```

## 故障排除

### 调试模式

```bash
./daily_report --config config.yaml --date today
```

查看数据源状态：
- `✅ meetings` - 收集成功
- `❌ meetings (错误信息)` - 收集失败

### 常见错误信息

| 错误信息 | 原因 | 解决方法 |
|----------|------|----------|
| `invalid app_id` | 应用ID无效 | 检查 app_id 配置 |
| `invalid user_id` | 用户ID无效 | 检查 user_id 配置 |
| `permission denied` | 权限不足 | 检查应用权限配置 |
| `rate limit exceeded` | API 限流 | 稍后重试或联系平台提升限额 |

## 下一步

配置好会议平台后，可以：
1. 运行 `./daily_report` 生成包含会议的日报
2. 使用自定义模板调整会议信息的显示格式

其他数据源集成文档：
- [Git 集成指南](./GIT_INTEGRATION.md)
- [Jira 集成指南](./JIRA_INTEGRATION.md)
- [Confluence 集成指南](./CONFLUENCE_INTEGRATION.md)