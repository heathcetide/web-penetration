#!/bin/bash

# 同步前端文件到 dist 目录
echo "🔄 同步前端文件..."

cp frontend/index.html frontend/dist/index.html
cp frontend/app.js frontend/dist/app.js

echo "✅ 同步完成！"
echo "📍 已更新:"
echo "   - frontend/dist/index.html"
echo "   - frontend/dist/app.js"

