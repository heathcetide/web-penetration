#!/bin/bash

echo "🔨 手动创建Wails绑定文件"
echo "========================"
echo ""

# 创建目录
mkdir -p frontend/wailsjs/go/main

# 创建App.js
cat > frontend/wailsjs/go/main/App.js << 'EOF'
export function ScanPorts(target, portRange) {
    return window.go?.main?.App?.ScanPorts(target, portRange) || 
           Promise.resolve(`扫描: ${target} (${portRange})`);
}

export function FuzzURL(url, wordlist) {
    return window.go?.main?.App?.FuzzURL(url, wordlist) || 
           Promise.resolve(`模糊测试: ${url}`);
}

export function TestSQLi(url, parameter) {
    return window.go?.main?.App?.TestSQLi(url, parameter) || 
           Promise.resolve(`SQL注入测试: ${url}`);
}

export function TestXSS(url, parameter) {
    return window.go?.main?.App?.TestXSS(url, parameter) || 
           Promise.resolve(`XSS测试: ${url}`);
}

export function BruteForce(url, username, passwordList) {
    return window.go?.main?.App?.BruteForce(url, username, passwordList) || 
           Promise.resolve(`暴力破解: ${url}`);
}

export function ScanFiles(url) {
    return window.go?.main?.App?.ScanFiles(url) || 
           Promise.resolve(`文件扫描: ${url}`);
}
EOF

# 创建App.d.ts
cat > frontend/wailsjs/go/main/App.d.ts << 'EOF'
export function ScanPorts(target: string, portRange: string): Promise<string>;
export function FuzzURL(url: string, wordlist: string): Promise<string>;
export function TestSQLi(url: string, parameter: string): Promise<string>;
export function TestXSS(url: string, parameter: string): Promise<string>;
export function BruteForce(url: string, username: string, passwordList: string): Promise<string>;
export function ScanFiles(url: string): Promise<string>;
EOF

echo "✅ 绑定文件已创建！"
echo ""
echo "文件位置:"
echo "  - frontend/wailsjs/go/main/App.js"
echo "  - frontend/wailsjs/go/main/App.d.ts"
echo ""
echo "现在运行: ./quick-fix.sh"

