# 调试端口扫描功能

## 问题分析

目前的问题是：前端点击扫描按钮没有反应。可能的原因：

1. **Wails绑定文件缺失** - `frontend/wailsjs/go/main/App.js` 不存在
2. **前端代码无法调用后端方法** - 导入路径错误
3. **应用未正确启动** - 需要确保 Structs 配置正确

## 已完成的修复

✅ 已在 `app.go` 中添加 `Structs: []interface{}{app}` 配置
✅ 端口扫描后端逻辑完整实现
✅ 前端调用代码已更新

## 下一步

需要用户手动运行：

```bash
# 1. 停止所有Wails进程
pkill -f wails

# 2. 清理并重新启动
cd /Users/cetide/Desktop/web-penetration
rm -rf frontend/wailsjs
./sync-frontend.sh

# 3. 启动应用
./dev.sh
```

然后等待绑定文件生成。

## 验证方法

启动后检查：

```bash
# 应该能看到这些文件
ls -la frontend/wailsjs/go/main/*.js
```

如果有 `App.js`，说明绑定成功。

## 手动创建绑定（如果自动生成失败）

如果Wails没有自动生成绑定，可以手动创建：

```bash
# 创建绑定目录
mkdir -p frontend/wailsjs/go/main

# 创建App.js（简化版本）
cat > frontend/wailsjs/go/main/App.js << 'EOF'
export function ScanPorts(target, portRange) {
    return window.Go.main.App.ScanPorts(target, portRange);
}

export function FuzzURL(url, wordlist) {
    return window.Go.main.App.FuzzURL(url, wordlist);
}

export function TestSQLi(url, parameter) {
    return window.Go.main.App.TestSQLi(url, parameter);
}

export function TestXSS(url, parameter) {
    return window.Go.main.App.TestXSS(url, parameter);
}

export function BruteForce(url, username, passwordList) {
    return window.Go.main.App.BruteForce(url, username, passwordList);
}

export function ScanFiles(url) {
    return window.Go.main.App.ScanFiles(url);
}
EOF

# 然后重启应用
pkill -f wails
./dev.sh
```

