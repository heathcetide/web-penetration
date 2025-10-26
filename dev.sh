#!/bin/bash

# 添加 Wails 到 PATH
export PATH=$PATH:$(go env GOPATH)/bin

# 启动开发模式
echo "🚀 启动 Web Penetration Tool 开发模式..."
echo "📝 提示: 按 Ctrl+C 退出"
echo ""

wails dev

