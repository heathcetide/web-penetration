# 故障排除指南

## ❌ 错误: "index.html: file does not exist"

**原因**: Wails 需要在 `frontend/dist/` 目录中找到文件，但文件在其他位置。

**解决方案**:
```bash
# 方法1: 使用同步脚本
./sync-frontend.sh

# 方法2: 手动复制
cp frontend/index.html frontend/dist/
cp frontend/app.js frontend/dist/
```

## ❌ 错误: "wails: command not found"

**原因**: Wails CLI 没有安装或不在 PATH 中。

**解决方案**:
```bash
# 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 添加到 PATH
export PATH=$PATH:$(go env GOPATH)/bin

# 验证安装
wails version
```

## ❌ 错误: "go run app.go" 不能运行

**原因**: Wails 应用需要特定的构建系统。

**解决方案**: 始终使用 `wails dev` 或 `./dev.sh` 运行。

## ✅ 正确的运行方式

```bash
# 开发模式（推荐）
./dev.sh

# 或手动运行
export PATH=$PATH:$(go env GOPATH)/bin
wails dev
```

## 🔄 修改前端文件后

1. 编辑 `frontend/index.html` 或 `frontend/app.js`
2. 运行 `./sync-frontend.sh` 同步文件
3. Wails 会自动检测并重新加载

## 📋 检查清单

启动前确保：
- [ ] Wails CLI 已安装
- [ ] PATH 包含 `$(go env GOPATH)/bin`
- [ ] `frontend/dist/index.html` 存在
- [ ] `frontend/dist/app.js` 存在
- [ ] `go mod tidy` 已运行

## 🆘 仍有问题？

1. 删除构建缓存：`rm -rf build/`
2. 清理 Go 缓存：`go clean -cache`
3. 重新运行：`./dev.sh`

