# Web Penetration Testing Tool

一个专业的Web渗透测试工具集，包含多种常见的Web安全测试功能。

## 🖥️ 桌面版 (Wails)

本项目现在包含一个**黑客风格的桌面GUI版本**，使用Wails v2 + Tailwind CSS构建。

### 快速启动桌面版

```bash
# 开发模式（最简单，推荐）
./dev.sh

# 或手动运行
export PATH=$PATH:$(go env GOPATH)/bin
wails dev

# 编译应用
./build.sh

# 运行编译后的应用
./build/bin/web-penetration.app/Contents/MacOS/web-penetration
```

**⚠️ 注意**: 不要使用 `go run app.go`！必须使用 `wails dev` 或 `./dev.sh`

详细使用说明请查看 [USAGE.md](./USAGE.md)

### 桌面版特性
- 🎨 黑客风格UI（绿色主题+发光效果）
- 📱 现代化响应式界面
- ⌨️ 快捷键导航
- 📊 实时终端输出
- 🎯 图形化操作所有功能

## 📦 CLI版本

同时保留命令行版本，可通过以下方式使用：

```bash
# 编译CLI版本
go build -o web-pen cmd/*.go main.go
```

## ⚠️ 免责声明

本工具仅用于合法的安全测试和教育目的。未经授权使用本工具进行任何攻击行为是违法的。使用者需自行承担使用本工具的所有法律责任。

## 功能列表

### 1. 端口扫描 (scan)
- 扫描指定主机的开放端口
- 支持端口范围指定和自定义端口列表
- 检测服务类型和版本信息

**用法示例:**
```bash
# 扫描默认端口范围 (1-1000)
web-pen scan example.com

# 扫描指定端口范围
web-pen scan example.com --ports=1-65535

# 扫描指定端口
web-pen scan example.com --ports=80,443,8080,8443
```

### 2. 模糊测试 (fuzz)
- 检测隐藏文件和目录
- 检测备份文件
- 检测配置文件
- 支持自定义字典

**用法示例:**
```bash
# 使用默认字典
web-pen fuzz https://example.com

# 使用自定义字典
web-pen fuzz https://example.com --wordlist=custom-wordlist.txt
```

### 3. 暴力破解 (brute)
- HTTP基础认证暴力破解
- 表单登录暴力破解
- 自定义密码字典
- 多线程支持

**用法示例:**
```bash
# HTTP基础认证暴力破解
web-pen brute https://example.com/admin --username=admin --password-list=passwords.txt

# 自定义密码字典
web-pen brute https://example.com/login --username=admin --password-list=top-passwords.txt
```

### 4. SQL注入检测 (sqli)
- 自动检测SQL注入漏洞
- 支持GET和POST参数测试
- 多种注入类型检测 (Union, Boolean-based, Time-based)
- Payload生成

**用法示例:**
```bash
# 检测默认参数
web-pen sqli https://example.com/page.php?id=1

# 检测指定参数
web-pen sqli https://example.com/page.php?id=1 --parameter=id

# 检测POST参数
web-pen sqli https://example.com/search.php --parameter=query
```

### 5. XSS漏洞检测 (xss)
- 反射型XSS检测
- 存储型XSS检测
- DOM型XSS检测
- 多种Payload测试

**用法示例:**
```bash
# 使用默认参数
web-pen xss https://example.com/search?q=test

# 指定测试参数
web-pen xss https://example.com/search?q=test --parameter=q

# 测试多个参数
web-pen xss https://example.com/comment --parameter=comment
```

### 6. 敏感文件扫描 (filescan)
- 自动扫描常见敏感文件
- 备份文件检测
- 配置文件检测
- .git/.svn等版本控制文件检测

**用法示例:**
```bash
# 扫描目标
web-pen filescan https://example.com
```

## 安装与使用

### 安装依赖
```bash
go mod download
```

### 编译
```bash
go build -o web-pen
```

### 运行
```bash
# 查看帮助
./web-pen --help

# 查看所有子命令
./web-pen

# 使用具体功能
./web-pen scan example.com
```

## 项目结构

```
web-penetration/
├── cmd/                    # 命令行命令
│   ├── root.go            # 根命令
│   ├── scanner.go         # 端口扫描
│   ├── fuzz.go            # 模糊测试
│   ├── bruteforce.go      # 暴力破解
│   ├── sqli.go            # SQL注入检测
│   ├── xss.go             # XSS检测
│   └── filescan.go        # 文件扫描
├── internal/              # 内部功能模块
│   ├── scanner/           # 端口扫描实现
│   ├── fuzzer/            # 模糊测试实现
│   ├── bruteforce/        # 暴力破解实现
│   ├── sqli/              # SQL注入检测实现
│   ├── xss/               # XSS检测实现
│   └── filescan/          # 文件扫描实现
├── main.go                # 程序入口
├── go.mod                 # Go模块定义
└── README.md              # 项目文档
```

## 开发计划

- [x] 实现端口扫描功能（TCP连接扫描）✅
- [ ] 实现模糊测试功能（目录和文件暴力破解）
- [ ] 实现暴力破解功能（HTTP认证暴力破解）
- [ ] 实现SQL注入检测（多种注入技术）
- [ ] 实现XSS检测（反射型、存储型、DOM型）
- [ ] 实现敏感文件扫描
- [ ] 添加代理支持
- [ ] 添加报告生成功能
- [ ] 添加多线程/并发支持
- [ ] 优化用户体验和输出格式

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License

