# Git 数据源集成指南

本文档介绍如何配置和使用 Git 数据收集器。

## 配置说明

### 基本配置

在 `config.yaml` 中配置 Git 数据源：

```yaml
git:
  author: "your.email@example.com"  # 必填：你的 Git 作者（名字或邮箱）
  repos:                                  # 可选：指定具体仓库路径
    - "/path/to/repo1"
    - "/path/to/repo2"
  repo_dirs:                              # 可选：扫描目录下的所有仓库
    - "/path/to/projects"
```

### 配置项说明

| 配置项 | 必填 | 说明 | 示例 |
|--------|------|------|------|
| `author` | 是 | Git 作者（名字或邮箱），用于过滤提交记录 | `"张三"` 或 `"user@example.com"` |
| `repos` | 否 | 指定具体的 Git 仓库路径列表 | `["/home/user/repo1"]` |
| `repo_dirs` | 否 | 扫描目录下的所有 Git 仓库 | `["/home/user/projects"]` |

## 工作原理

### 作者匹配

Git 收集器使用 `git log --author=<author>` 命令匹配作者。

```bash
git log --author="张三" --since="2026-02-11 00:00:00" --until="2026-02-11 23:59:59"
```

**注意：** 
- 作者可以是名字或邮箱，Git 会进行模糊匹配
- 如果使用了多个名字或邮箱，建议使用其中一个
- 匹配时会同时匹配作者名和邮箱字段

### 多仓库支持

Git 收集器支持两种配置方式：

**方式 1：指定具体仓库（`repos`）**
```yaml
git:
  author: "user@example.com"
  repos:
    - "/home/user/project1"
    - "/home/user/project2"
```

**方式 2：扫描目录（`repo_dirs`）**
```yaml
git:
  author: "user@example.com"
  repo_dirs:
    - "/home/user/projects"  # 自动扫描该目录下所有包含 .git 的子目录
```

两种方式可以同时使用，工具会合并两个列表。

### 目录扫描逻辑

当使用 `repo_dirs` 时：
1. 遍历指定目录
2. 查找包含 `.git` 文件夹的子目录
3. 将发现的每个仓库加入收集列表
4. 对每个仓库执行 git log 命令

### 输出数据

每个 Git 提交包含以下信息：

| 字段 | 说明 | 示例 |
|------|------|------|
| `type` | 数据源类型 | `"git"` |
| `title` | 提交消息 | `"feat: 添加新功能"` |
| `time` | 提交时间 | `2026-02-11 14:30:00 +0800` |
| `link` | 提交链接 | `/path/to/repo/commit/abc123` |
| `metadata.repo` | 仓库名称 | `"project1"` |
| `metadata.commit` | 提交哈希 | `"abc123456..."` |
| `metadata.author` | 作者姓名 | `"张三"` |
| `metadata.author_email` | 作者邮箱 | `"zhangsan@example.com"` |

## 常见问题

### 1. 找不到提交记录

**可能原因：**
- Git 配置的作者与提交记录中的作者不匹配
- 仓库路径配置错误
- 仓库不是 Git 仓库

**排查方法：**
```bash
# 检查提交记录中的作者
git log --format="%an <%ae>"

# 检查是否是 Git 仓库
cd /path/to/repo
git rev-parse --git-dir
```

### 2. 没有扫描到仓库

**可能原因：**
- `repo_dirs` 路径错误
- 目录下没有包含 `.git` 文件夹的仓库

**排查方法：**
```bash
# 查找目录下的所有 Git 仓库
find /path/to/projects -name ".git" -type d
```

### 3. 提交时间不正确

**可能原因：**
- 系统时区配置不正确
- Git 提交时区与系统时区不一致

**解决方法：**
```yaml
time:
  timezone: "Asia/Shanghai"  # 配置正确的时区
```

### 4. 权限问题

**可能原因：**
- 对仓库目录没有读取权限
- Git 命令执行权限不足

**解决方法：**
```bash
# 检查目录权限
ls -la /path/to/repo

# 测试 Git 命令
cd /path/to/repo
git log
```

## 性能优化

### 大量仓库

如果需要扫描大量仓库：
- 使用 `repo_dirs` 而非 `repos`（自动发现）
- 将仓库按项目分组，使用多个 `repo_dirs` 配置
- 考虑并行处理（当前版本为串行处理）

### 减少扫描时间

- 精确配置 `repo_dirs`，避免扫描不相关的目录
- 使用 `repos` 明确指定需要扫描的仓库

## 示例配置

### 开发者个人电脑

```yaml
git:
  author: "张三"
  repo_dirs:
    - "/home/developer/projects"
    - "/home/developer/work"
```

### 公司服务器

```yaml
git:
  author: "developer@company.com"
  repos:
    - "/opt/repos/frontend"
    - "/opt/repos/backend"
    - "/opt/repos/shared"
```

### 混合配置

```yaml
git:
  author: "张三"
  repos:
    - "/special/repo1"  # 特殊仓库
  repo_dirs:
    - "/home/developer/projects"  # 常规项目目录
```

## 技术细节

### Git 命令格式

```bash
git -C <repo_path> log \
  --author="<author>" \
  --since="<start_time>" \
  --until="<end_time>" \
  --pretty=format:%H|%an|%ae|%ai|%s
```

### 输出格式

`pretty=format` 输出格式：
```
<commit_hash>|<author_name>|<author_email>|<author_time>|<commit_message>
```

示例：
```
abc123|张三|zhangsan@example.com|2026-02-11 14:30:00 +0800|feat: 添加新功能
```

## 故障排除

### 调试模式

如果遇到问题，可以查看工具输出中的错误信息：

```bash
./daily_report --config config.yaml --date today
```

工具会在数据源状态中显示 Git 收集器的状态：
- `✅ git` - 收集成功
- `❌ git (错误信息)` - 收集失败，会显示具体错误

### 常见错误信息

| 错误信息 | 原因 | 解决方法 |
|----------|------|----------|
| `not a git repository: /path/to/repo` | 指定路径不是 Git 仓库 | 检查路径是否正确 |
| `failed to get git log` | Git 命令执行失败 | 检查 Git 是否安装，仓库是否有效 |
| `no repositories found` | 没有找到任何仓库 | 检查 `repos` 和 `repo_dirs` 配置 |

## 下一步

配置好 Git 数据源后，可以：
1. 运行 `./daily_report` 生成日报
2. 使用 `--date yesterday` 生成昨天的日报
3. 使用自定义模板美化输出格式

其他数据源集成文档：
- [Jira 集成指南](./JIRA_INTEGRATION.md)
- [Confluence 集成指南](./CONFLUENCE_INTEGRATION.md)
- [会议平台集成指南](./MEETING_INTEGRATION.md)
