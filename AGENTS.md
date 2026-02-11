
# AGENTS.md - 项目上下文文档

## 项目概述

**项目名称:** daily_report

**项目类型:** Go 代码项目

**项目状态:** 初始化阶段

**目的:** 这是一个使用 Go 语言开发的项目，模块名为 `daily_report`。目前项目仅包含基础的 Go 模块配置（`go.mod`），尚未实现任何具体功能。

## 技术栈

- **编程语言:** Go (Golang)
- **模块名称:** daily_report
- **Go 版本:** 未指定（使用 Go 默认版本）

## 项目结构

```
daily_report/
├── .idea/              # JetBrains IDE 配置目录
│   ├── .gitignore
│   ├── copilot.data.migration.ask2agent.xml
│   ├── daily_report.iml
│   ├── modules.xml
│   └── workspace.xml
├── AGENTS.md           # 本文档 - AI 代理上下文
└── go.mod              # Go 模块定义文件
```

## 构建和运行

### 初始化依赖

```bash
go mod tidy
```

### 构建项目

```bash
go build
```

### 运行项目

```bash
go run .
```

### 运行测试

```bash
go test ./...
```

**注意:** 由于项目尚未实现任何功能，上述命令可能需要根据实际开发需求进行调整。

## 开发约定

### 代码结构建议

根据 Go 项目标准实践，建议遵循以下目录结构：

```
daily_report/
├── cmd/                # 主应用程序入口
│   └── daily_report/
│       └── main.go
├── internal/           # 私有应用代码
├── pkg/                # 可被外部使用的库代码
├── api/                # API 定义
├── configs/            # 配置文件
├── scripts/            # 构建、部署脚本
├── test/               # 测试辅助文件
└── docs/               # 项目文档
```

### 编码规范

- 遵循 Go 官方代码风格指南：https://golang.org/doc/effective_go
- 使用 `gofmt` 格式化代码
- 使用 `go vet` 进行静态检查
- 为导出的函数、类型、常量编写文档注释

### 版本控制

- 项目当前未初始化为 Git 仓库
- 建议添加 `.gitignore` 文件以排除 IDE 配置和编译产物

## 待实现功能

项目目前处于初始状态，需要根据实际需求确定功能方向并实现。可能的开发方向包括：

- 命令行工具
- Web 服务
- 库/SDK
- 数据处理工具

## 依赖管理

项目使用 Go Modules 进行依赖管理。依赖声明在 `go.mod` 文件中。

### 添加依赖

```bash
go get <package-name>
```

### 移除未使用的依赖

```bash
go mod tidy
```

## 备注

- 本文档由 iFlow CLI 自动生成，用于为后续的 AI 交互提供项目上下文
- 随着项目发展，建议定期更新本文档以反映最新的项目状态和约定
