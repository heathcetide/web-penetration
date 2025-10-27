#!/bin/bash

echo "🔧 快速修复端口扫描功能"
echo "======================"
echo ""

# 停止所有Wails进程
echo "1️⃣ 停止现有进程..."
pkill -f "wails dev" 2>/dev/null || true
pkill -f "wails dev" 2>/dev/null || true
sleep 2

# 清理旧绑定
echo "2️⃣ 清理旧绑定..."
rm -rf frontend/wailsjs

# 同步前端文件
echo "3️⃣ 同步前端文件..."
./sync-frontend.sh

# 检查structs配置
echo "4️⃣ 检查配置..."
if grep -q "Structs:" app.go; then
    echo "✅ app.go 配置正确"
else
    echo "❌ app.go 配置错误"
    exit 1
fi

# 启动应用
echo ""
echo "5️⃣ 启动应用..."
echo "提示: 如果绑定文件未自动生成，请按 Ctrl+C，然后运行:"
echo "   ./create-binding.sh"
echo ""
echo ""

export PATH=$PATH:$(go env GOPATH)/bin

echo "启动Wails开发模式..."
wails dev

