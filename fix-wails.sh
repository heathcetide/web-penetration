#!/bin/bash

echo "🔧 修复Wails绑定问题..."
echo ""

# 停止现有进程
pkill -f "wails dev" 2>/dev/null || true
sleep 1

# 清理旧的绑定
echo "1️⃣ 清理旧绑定..."
rm -rf frontend/wailsjs
mkdir -p frontend/wailsjs/go/main
mkdir -p frontend/wailsjs/runtime

# 同步前端文件
echo "2️⃣ 同步前端文件..."
./sync-frontend.sh

# 重新生成绑定
echo "3️⃣ 启动Wails（会自动生成绑定）..."
export PATH=$PATH:$(go env GOPATH)/bin

# 后台启动
nohup wails dev > /tmp/wails-dev.log 2>&1 &
WAILS_PID=$!

echo "✅ Wails已启动 (PID: $WAILS_PID)"
echo "📋 查看日志: tail -f /tmp/wails-dev.log"
echo ""
echo "等待5秒后检查绑定文件..."

sleep 5

if [ -f "frontend/wailsjs/go/main/App.js" ]; then
    echo "✅ 绑定文件已生成！"
    cat frontend/wailsjs/go/main/App.js | head -20
else
    echo "❌ 绑定文件仍未生成"
    echo "正在查看日志..."
    tail -20 /tmp/wails-dev.log
fi

echo ""
echo "🌐 应用地址: http://localhost:34115"
echo "按 Ctrl+C 停止并查看日志"

