## ADDED Requirements

### Requirement: 数据收集器接口
系统 SHALL 定义统一的数据收集器接口，支持从不同数据源收集当日工作产出数据。

#### Scenario: 实现新的数据收集器
- **WHEN** 开发者需要添加新的数据源支持
- **THEN** 可以通过实现 Collector 接口来创建新的收集器
- **AND** 收集器必须返回统一的数据格式

#### Scenario: 配置数据收集器
- **WHEN** 系统初始化时
- **THEN** 可以从配置文件加载各数据收集器的配置
- **AND** 支持通过环境变量覆盖敏感配置

### Requirement: 用户过滤配置
系统 SHALL 支持为各数据源配置用户标识，用于过滤个人工作产出。

#### Scenario: 配置 Git 用户
- **WHEN** 配置 Git 收集器
- **THEN** 用户可以配置 author_email 用于精确匹配提交记录
- **AND** 支持配置具体的仓库路径（repos）
- **AND** 支持配置目录路径自动扫描（repo_dirs）

#### Scenario: 配置会议平台用户
- **WHEN** 配置会议收集器
- **THEN** 用户可以配置平台类型（feishu/dingtalk/wecom）和用户标识
- **AND** 支持配置 API 凭证（app_id, app_secret）

#### Scenario: 配置 Jira 用户
- **WHEN** 配置 Jira 收集器
- **THEN** 用户可以配置用户名、URL、API Token
- **AND** 可选配置项目 Key 以筛选特定项目

#### Scenario: 配置 Confluence 用户
- **WHEN** 配置 Confluence 收集器
- **THEN** 用户可以配置用户名、URL、API Token
- **AND** 可选配置空间 Key 以筛选特定空间

#### Scenario: 环境变量支持
- **WHEN** 配置文件中使用 ${VAR_NAME} 格式
- **THEN** 系统从环境变量中读取对应的值
- **AND** 用于保护敏感信息（如 API Token）

### Requirement: Git 提交记录收集
系统 MUST 能够从 Git 仓库收集指定时间范围内的提交记录，并按作者邮箱过滤。

#### Scenario: 收集当日 Git 提交
- **WHEN** 执行日报生成命令
- **THEN** 系统使用配置的 author_email 精确匹配提交记录
- **AND** 获取当天 00:00 到 23:59 的所有匹配提交
- **AND** 提取 commit hash、作者、时间、消息、变更文件数等信息

#### Scenario: 支持多仓库
- **WHEN** 配置文件中指定多个 Git 仓库路径
- **THEN** 系统依次收集所有仓库的提交记录
- **AND** 在报告中按仓库分类显示

#### Scenario: 扫描目录下的所有仓库
- **WHEN** 配置文件中指定了 repo_dirs 目录
- **THEN** 系统自动扫描该目录下所有包含 .git 文件夹的子目录
- **AND** 将发现的所有仓库加入收集列表
- **AND** 支持同时配置 repos 和 repo_dirs

#### Scenario: 邮箱精确匹配
- **WHEN** 配置了 author_email
- **THEN** 系统使用 git log --author 参数精确匹配邮箱
- **AND** 只返回作者邮箱匹配的提交记录

### Requirement: 会议记录收集
系统 MUST 能够从飞书、钉钉、企业微信等平台收集当日会议记录，并按用户参与过滤。

#### Scenario: 收集飞书会议记录
- **WHEN** 配置了飞书 API 凭证和用户 ID
- **THEN** 系统通过飞书 API 获取用户参与的当天会议日程
- **AND** 提取会议主题、时间、参会人员、会议链接等信息

#### Scenario: 收集钉钉会议记录
- **WHEN** 配置了钉钉 API 凭证和用户 ID
- **THEN** 系统通过钉钉 API 获取用户参与的当天会议日程
- **AND** 提取会议主题、时间、参会人员、会议链接等信息

#### Scenario: 收集企业微信会议记录
- **WHEN** 配置了企业微信 API 凭证和用户 ID
- **THEN** 系统通过企业微信 API 获取用户参与的当天会议日程
- **AND** 提取会议主题、时间、参会人员、会议链接等信息

#### Scenario: 平台类型配置
- **WHEN** 配置会议收集器
- **THEN** 用户可以指定平台类型（feishu/dingtalk/wecom）
- **AND** 系统根据平台类型调用对应的 API

### Requirement: Jira 任务变更收集
系统 MUST 能够从 Jira 收集指定时间范围内的任务变更记录，并按用户过滤。

#### Scenario: 收集当日任务变更（简化方案）
- **WHEN** 执行日报生成命令
- **THEN** 系统查询分配给当前用户或由当前用户创建的任务
- **AND** 筛选出在指定时间范围内有更新的任务（updated 字段）
- **AND** 提取任务 ID、标题、状态、分配人、更新时间等信息

#### Scenario: 筛选特定项目
- **WHEN** 配置文件中指定了 Jira 项目 Key
- **THEN** 系统仅收集指定项目的任务变更
- **AND** 忽略其他项目的数据

#### Scenario: JQL 查询构建
- **WHEN** 构建 Jira 查询
- **THEN** 使用 JQL: `(assignee = {username} OR reporter = {username}) AND updated >= '{start}' AND updated <= '{end}'`
- **AND** {username} 从配置文件中读取

### Requirement: Confluence 文档产出收集
系统 MUST 能够从 Confluence 收集指定时间范围内的文档，并按用户过滤。

#### Scenario: 收集当日文档产出（并集逻辑）
- **WHEN** 执行日报生成命令
- **THEN** 系统收集今天创建的文档（creator = me AND created >= start）
- **AND** 系统收集今天修改的文档（lastModifier = me AND lastUpdated >= start）
- **AND** 使用并集（OR）逻辑，满足任一条件即视为我的产出
- **AND** 提取文档标题、作者、更新时间、链接等信息

#### Scenario: 筛选特定空间
- **WHEN** 配置文件中指定了 Confluence 空间 Key
- **THEN** 系统仅收集指定空间的文档
- **AND** 忽略其他空间的文档

#### Scenario: CQL 查询构建
- **WHEN** 构建 Confluence 查询
- **THEN** 使用 CQL: `(creator = {username} AND created >= '{start}') OR (lastModifier = {username} AND lastUpdated >= '{start}')`
- **AND** {username} 从配置文件中读取

### Requirement: 时间范围处理
系统 MUST 支持自动检测当天时间范围，并支持自定义时间范围。

#### Scenario: 自动检测当天时间
- **WHEN** 用户未指定时间范围
- **THEN** 系统使用系统时区的当天 00:00:00 到 23:59:59 作为时间范围
- **AND** 将此时间范围传递给所有数据收集器

#### Scenario: 支持自定义时间范围
- **WHEN** 用户通过 CLI 参数指定 --from 和 --to
- **THEN** 系统使用用户指定的时间范围
- **AND** 将此时间范围传递给所有数据收集器

### Requirement: 错误处理和重试
系统 MUST 处理数据收集过程中的错误，并在失败时提供友好的错误信息。

#### Scenario: 数据源连接失败
- **WHEN** 某个数据源连接失败（如网络问题、认证失败）
- **THEN** 系统记录错误日志并继续执行其他数据收集器
- **AND** 在最终报告中注明该数据源收集失败

#### Scenario: API 限流处理
- **WHEN** 数据源 API 返回限流错误
- **THEN** 系统等待适当时间后重试
- **AND** 如果重试次数超过阈值，则跳过该数据源并记录错误