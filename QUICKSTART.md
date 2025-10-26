# 快速开始指南

## 🎯 如何运行应用

**重要**: Wails 应用不能使用 `go run` 运行！必须使用 Wails 命令。

### 方式1: 开发模式（推荐）✨

```bash
# 首次运行前，确保文件在正确位置
./sync-frontend.sh

# 使用提供的脚本（最简单）
./dev.sh

# 或者手动运行
export PATH=$PATH:$(go env GOPATH)/bin
wails dev
```

**提示**: 修改 `frontend/index.html` 或 `frontend/app.js` 后，运行 `./sync-frontend.sh` 来同步文件

**开发模式特点**:
- ✅ 支持热重载（修改前端代码自动刷新）
- ✅ 显示开发者工具
- ✅ 实时查看控制台日志
- ✅ 按 `Ctrl+C` 退出

### 方式2: 编译后运行

```bash
# 编译应用
./build.sh

# 或手动编译
export PATH=$PATH:$(go env GOPATH)/bin
wails build

# 运行编译后的应用
./build/bin/web-penetration.app/Contents/MacOS/web-penetration
```

## 🔧 常见问题

### Q: 为什么不能用 `go run app.go`？

A: Wails 应用使用了自定义构建系统，需要特定的构建标签和资源嵌入。直接使用 `go run` 无法正确编译和运行。

### Q: 提示 "wails: command not found"

A: 需要先安装 Wails CLI：
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Q: 开发模式无法启动

A: 确保：
1. 已安装 Wails CLI
2. 已运行 `go mod download`
3. PATH 环境变量包含 `$(go env GOPATH)/bin`

### Q: 如何查看日志？

开发模式下，控制台会显示所有日志。也可以在应用内右键点击，选择"开发者工具"查看浏览器控制台。

## 📖 下一步

1. 运行 `./dev.sh` 启动应用
2. 查看 [USAGE.md](./USAGE.md) 了解功能使用
3. 查看 [README.md](./README.md) 了解项目详情

## 🎨 界面预览

```
┌─────────────────────────────────────────────────┐
│ Web Penetration Tool v1.0  ● Connected      │
├──────────┬──────────────────────────────────┤
│[1] 端口  │  [输入框] Target: _________      │
│[2] 模糊  │  [输入框] Ports: 1-1000          │
│[3] SQLI  │  [按钮]  执行扫描                │
│[4] XSS   │                                  │
│[5] 暴力  │                                  │
│[6] 文件  │                                  │
├──────────┴──────────────────────────────────┤
│ Terminal Output:                            │
│ [时间] [INFO] 扫描目标: example.com        │
│ [时间] [OK] 发现端口: 80, 443              │
└─────────────────────────────────────────────┘
```

现在运行 `./dev.sh` 来启动应用吧！🚀

