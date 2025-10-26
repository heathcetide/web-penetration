#!/bin/bash

# 测试端口扫描功能的CLI脚本

echo "🧪 测试端口扫描功能"
echo "==================="
echo ""

# 编译CLI版本
echo "📦 编译CLI版本..."
go build -o web-pen-cli ./cmd/*.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"
echo ""

# 测试本地扫描
echo "🔍 测试本地端口扫描 (扫描 127.0.0.1:80,443)..."
./web-pen-cli scan 127.0.0.1 --ports=80,443

echo ""
echo "✨ 测试完成！"
echo ""
echo "💡 提示: 可以使用 ./dev.sh 启动GUI版本进行更详细的测试"

# 清理
rm web-pen-cli

