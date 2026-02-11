## ADDED Requirements

### Requirement: 报告生成模式
系统 SHALL 支持两种报告生成模式：模板渲染和 LLM 美化。

#### Scenario: 默认模板渲染模式
- **WHEN** 用户未指定生成模式
- **THEN** 系统使用模板渲染生成结构化 Markdown 报告
- **AND** 数据完全可控，可追溯数据源

#### Scenario: LLM 美化模式
- **WHEN** 用户通过 --mode llm 参数指定
- **THEN** 系统将收集的数据发送给 LLM 进行美化
- **AND** 生成自然语言风格的日报，包含工作总结和亮点

#### Scenario: 模式配置
- **WHEN** 配置文件中指定 report.mode
- **THEN** 系统使用配置的模式生成报告
- **AND** 命令行参数可覆盖配置文件

### Requirement: Markdown 报告生成
系统 MUST 能够将收集到的数据生成为格式化的 Markdown 报告。

#### Scenario: 生成基础报告
- **WHEN** 数据收集完成
- **THEN** 系统生成包含所有数据源的 Markdown 报告
- **AND** 报告按数据源分类（Git、会议、任务、文档）
- **AND** 每个数据源下的条目按时间倒序排列

#### Scenario: 报告结构
- **WHEN** 生成报告
- **THEN** 报告包含标题、日期、汇总统计、各数据源章节
- **AND** 每个数据源章节包含该数据源的条目列表
- **AND** 每个条目包含时间、标题、链接等关键信息

#### Scenario: 汇总统计
- **WHEN** 生成报告
- **THEN** 报告顶部显示每个数据源的条目数量
- **AND** 包括 Git 提交数、会议数、任务变更数、文档数

### Requirement: LLM 配置
系统 SHALL 支持 LLM 美化模式的配置。

#### Scenario: 配置 LLM 提供商
- **WHEN** 配置 LLM 功能
- **THEN** 用户可以指定提供商（openai、anthropic、local）
- **AND** 配置模型名称和 API 密钥

#### Scenario: 配置系统提示词
- **WHEN** 配置 LLM 功能
- **THEN** 用户可以自定义 system_prompt
- **AND** 用于指导 LLM 生成特定风格的日报

#### Scenario: 环境变量支持
- **WHEN** 配置 LLM API 密钥
- **THEN** 支持使用环境变量 ${LLM_API_KEY}
- **AND** 保护敏感信息

### Requirement: 自定义模板支持
系统 MUST 支持使用自定义 Markdown 模板生成报告。

#### Scenario: 使用默认模板
- **WHEN** 用户未指定自定义模板
- **THEN** 系统使用内置的默认模板生成报告
- **AND** 默认模板包含所有标准数据源章节

#### Scenario: 使用自定义模板
- **WHEN** 用户在配置文件中指定了自定义模板路径
- **THEN** 系统加载并使用该模板生成报告
- **AND** 模板可以自定义报告的结构和样式

### Requirement: 输出方式
系统 MUST 支持多种输出方式，包括终端输出和文件输出。

#### Scenario: 输出到终端
- **WHEN** 用户未指定输出文件
- **THEN** 系统将生成的 Markdown 报告输出到终端
- **AND** 支持终端的 Markdown 渲染（如使用代码块显示）

#### Scenario: 输出到文件
- **WHEN** 用户通过 --output 参数指定输出文件路径
- **THEN** 系统将生成的 Markdown 报告写入指定文件
- **AND** 如果文件已存在，则覆盖原有内容

### Requirement: Web 界面支持
系统 MUST 提供 Web 界面，允许用户通过浏览器查看和生成日报。

#### Scenario: 访问 Web 界面
- **WHEN** 用户启动 Web 服务并访问浏览器
- **THEN** 系统显示日报生成界面
- **AND** 界面包含日期选择器和生成按钮

#### Scenario: 通过 Web 生成日报
- **WHEN** 用户在 Web 界面选择日期并点击生成
- **THEN** 系统调用数据收集和报告生成逻辑
- **AND** 在页面中渲染生成的 Markdown 报告

#### Scenario: Web API 端点
- **WHEN** 客户端调用 POST /api/report 端点
- **THEN** 系统接收日期范围参数并返回生成的 Markdown 报告
- **AND** 支持 JSON 和纯文本两种响应格式

### Requirement: 数据汇总统计
系统 MUST 在报告中提供汇总统计信息。

#### Scenario: 统计各数据源条目数量
- **WHEN** 生成报告
- **THEN** 报告顶部显示每个数据源的条目数量
- **AND** 包括 Git 提交数、会议数、任务变更数、文档数

#### Scenario: 空数据处理
- **WHEN** 某个数据源在指定时间范围内无数据
- **THEN** 报告中该数据源章节显示"无记录"
- **AND** 不在统计中计算该数据源

### Requirement: 报告元数据
系统 MUST 在报告中包含生成时间和数据源状态信息。

#### Scenario: 包含生成时间
- **WHEN** 生成报告
- **THEN** 报告标题下显示报告生成时间
- **AND** 格式为 "生成于: YYYY-MM-DD HH:MM:SS"

#### Scenario: 包含数据源状态
- **WHEN** 生成报告
- **THEN** 报告末尾显示各数据源的收集状态
- **AND** 标注成功或失败的数据源
- **AND** 如果失败，显示错误信息