# Web Penetration Tool - 桌面版使用指南

## 🚀 快速开始

### 运行应用

#### 方式1: 使用编译后的应用（推荐）
```bash
./build/bin/web-penetration.app/Contents/MacOS/web-penetration
```

#### 方式2: 开发模式
```bash
export PATH=$PATH:$(go env GOPATH)/bin
wails dev
```

开发模式支持热重载，修改前端代码后会自动刷新。

## 🎨 界面介绍

### 黑客风格设计
- **主题色**: 黑客绿色 (#00ff00)
- **背景**: 纯黑
- **字体**: 等宽字体 (Share Tech Mono)
- **特效**: 文字发光、扫描线效果

### 功能模块

#### 1. 端口扫描 (Port Scanner)
- **位置**: 左侧导航 [1]
- **功能**: 扫描目标主机的开放端口
- **输入**: 
  - 目标主机 (如: example.com)
  - 端口范围 (如: 1-1000)

#### 2. 模糊测试 (Fuzzing)
- **位置**: 左侧导航 [2]
- **功能**: 发现隐藏文件和目录
- **输入**:
  - 目标URL (如: https://example.com)
  - 字典文件路径

#### 3. SQL注入检测 (SQL Injection)
- **位置**: 左侧导航 [3]
- **功能**: 检测SQL注入漏洞
- **输入**:
  - 目标URL
  - 参数名称

#### 4. XSS漏洞检测 (Cross-Site Scripting)
- **位置**: 左侧导航 [4]
- **功能**: 检测XSS漏洞
- **输入**:
  - 目标URL
  - 参数名称

#### 5. 暴力破解 (Brute Force)
- **位置**: 左侧导航 [5]
- **功能**: HTTP认证暴力破解
- **输入**:
  - 目标URL
  - 用户名
  - 密码字典

#### 6. 文件扫描 (File Scanner)
- **位置**: 左侧导航 [6]
- **功能**: 扫描敏感文件
- **输入**:
  - 目标URL

## 📊 输出终端

底部终端显示所有操作的实时输出，支持：
- 时间戳
- 颜色分类（信息/成功/警告/危险）
- 滚动查看历史

## 🔧 技术栈

- **后端**: Go 1.21
- **前端框架**: Wails v2
- **前端库**: Tailwind CSS
- **样式**: 纯CSS实现黑客风格

## 🏗️ 项目结构

```
web-penetration/
├── app.go                # Wails应用入口
├── internal/
│   └── app.go           # 业务逻辑（与前端交互）
├── frontend/
│   ├── index.html       # 主界面HTML
│   └── app.js           # 前端逻辑
├── cmd/                  # CLI命令（保留）
└── build/               # 编译输出
```

## 💻 开发

### 修改前端
直接编辑 `frontend/index.html` 和 `frontend/app.js`

### 修改后端逻辑
编辑 `internal/app.go`，函数会自动暴露给前端

### 添加新功能
1. 在 `internal/app.go` 添加新的Go函数
2. 运行 `wails dev` 重新生成绑定
3. 在 `frontend/app.js` 中调用新的函数

## ⚠️ 重要提示

本工具仅用于合法的安全测试目的。使用前请确保：
1. 获得目标系统的明确授权
2. 遵守当地法律法规
3. 仅用于教育和个人学习

## 📝 TODO

当前版本为演示版，核心功能需要实现：
- [ ] 端口扫描算法
- [ ] 模糊测试实现
- [ ] SQL注入检测逻辑
- [ ] XSS检测算法
- [ ] 暴力破解实现
- [ ] 文件扫描功能
- [ ] 添加进度条
- [ ] 支持多线程
- [ ] 导出报告功能

## 🐛 故障排除

### 问题: 应用无法启动
**解决**: 检查是否已安装Wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 问题: 前端没有显示
**解决**: 确保 `frontend/dist` 目录存在（为空也可以）

### 问题: 功能没有反应
**解决**: 打开浏览器开发者工具查看JavaScript错误

## 📖 相关文档

- [Wails官方文档](https://wails.io/docs/gettingstarted/introduction)
- [Tailwind CSS](https://tailwindcss.com/)
- [Go语言](https://golang.org/)

