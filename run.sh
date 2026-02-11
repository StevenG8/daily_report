#!/bin/bash

# Daily Report Generator - 运行脚本
# 用法: ./run.sh [选项]

set -e

# 默认配置
CONFIG_FILE="config.yaml"
DATE_RANGE="today"
OUTPUT_FILE=""
MODE="template"
TEMPLATE_FILE=""

# 显示帮助信息
show_help() {
    cat << EOF
Daily Report Generator - 运行脚本

用法: ./run.sh [选项]

选项:
    -c, --config FILE      配置文件路径 (默认: config.yaml)
    -d, --date RANGE       日期范围 (默认: today)
                           可选值: today, yesterday, YYYY-MM-DD, YYYY-MM-DD,YYYY-MM-DD
    -o, --output FILE      输出文件路径 (默认: 终端输出)
    -m, --mode MODE        报告模式 (默认: template)
                           可选值: template, llm
    -t, --template FILE    自定义模板文件路径
    -h, --help             显示此帮助信息

示例:
    # 生成今天的日报
    ./run.sh

    # 生成昨天的日报
    ./run.sh -d yesterday

    # 输出到文件
    ./run.sh -o report.md

    # 使用自定义模板
    ./run.sh -t custom_template.md

    # 指定配置文件
    ./run.sh -c my_config.yaml

环境变量:
    JIRA_API_TOKEN         Jira API 令牌
    OPENAI_API_KEY         OpenAI API 密钥
    FEISHU_APP_ID          飞书应用 ID
    FEISHU_APP_SECRET      飞书应用密钥

EOF
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -d|--date)
            DATE_RANGE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -m|--mode)
            MODE="$2"
            shift 2
            ;;
        -t|--template)
            TEMPLATE_FILE="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "错误: 配置文件不存在: $CONFIG_FILE"
    echo "提示: 可以复制示例配置文件: cp examples/config.example.yaml config.yaml"
    exit 1
fi

# 检查二进制文件是否存在，如果不存在则构建
if [ ! -f "daily_report" ]; then
    echo "正在构建 daily_report..."
    go build -o daily_report ./cmd/cli
    echo "构建完成！"
fi

# 构建命令
CMD="./daily_report --config $CONFIG_FILE --date $DATE_RANGE --mode $MODE"

# 添加可选参数
if [ -n "$OUTPUT_FILE" ]; then
    CMD="$CMD --output $OUTPUT_FILE"
fi

if [ -n "$TEMPLATE_FILE" ]; then
    CMD="$CMD --template $TEMPLATE_FILE"
fi

# 执行命令
echo "正在生成日报..."
echo "配置文件: $CONFIG_FILE"
echo "日期范围: $DATE_RANGE"
echo "模式: $MODE"
if [ -n "$OUTPUT_FILE" ]; then
    echo "输出文件: $OUTPUT_FILE"
fi
echo "---"

eval $CMD

echo "---"
echo "日报生成完成！"