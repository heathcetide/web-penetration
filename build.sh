#!/bin/bash

# 添加 Wails 到 PATH
export PATH=$PATH:$(go env GOPATH)/bin

# 编译应用
echo "🔨 编译 Web Penetration Tool..."
echo ""

wails build

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 编译成功！"
    echo "📍 应用位置: ./build/bin/web-penetration.app/Contents/MacOS/web-penetration"
    echo ""
    echo "运行应用:"
    echo "  ./build/bin/web-penetration.app/Contents/MacOS/web-penetration"
fi

