import * as wails from './wailsjs/go/main/App';
import { LogPrint } from './wailsjs/runtime';

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
            if (target) {
                addOutput(`[SCAN] 开始扫描: ${target}`, 'info');
                addOutput(`端口范围: ${ports}`, 'info');
                
                try {
                    const result = await wails.ScanPorts(target, ports);
                    const results = JSON.parse(result);
                    
                    if (results.length > 0 && results[0].port !== 0) {
                        addOutput(`✅ 发现 ${results.length} 个开放端口:`, 'success');
                        results.forEach(r => {
                            addOutput(`  端口 ${r.port}: ${r.status}`, 'success');
                        });
                    } else {
                        addOutput('❌ 未发现开放端口', 'warning');
                    }
                } catch (err) {
                    addOutput(`❌ 错误: ${err.message}`, 'danger');
                }
                
                addOutput('', 'info');
            }
        });
    }

    // Fuzz
    if (elements.fuzzBtn) {
        elements.fuzzBtn.addEventListener('click', async () => {
            const url = elements.fuzzInput.value;
            const wordlist = document.getElementById('fuzz-wordlist').value || 'wordlist.txt';
            if (url) {
                const result = await wails.FuzzURL(url, wordlist);
                addOutput(`[FUZZ] 目标: ${url}`, 'info');
                addOutput(`字典: ${wordlist}`, 'info');
                addOutput(`结果: ${result}`, 'success');
                addOutput('', 'info');
            }
        });
    }

    // SQLi
    if (elements.sqliBtn) {
        elements.sqliBtn.addEventListener('click', async () => {
            const url = elements.sqliInput.value;
            const param = document.getElementById('sqli-param').value || 'id';
            if (url) {
                const result = await wails.TestSQLi(url, param);
                addOutput(`[SQLI] 目标: ${url}`, 'info');
                addOutput(`参数: ${param}`, 'info');
                addOutput(`结果: ${result}`, 'warning');
                addOutput('', 'info');
            }
        });
    }

    // XSS
    if (elements.xssBtn) {
        elements.xssBtn.addEventListener('click', async () => {
            const url = elements.xssInput.value;
            const param = document.getElementById('xss-param').value || 'search';
            if (url) {
                const result = await wails.TestXSS(url, param);
                addOutput(`[XSS] 目标: ${url}`, 'info');
                addOutput(`参数: ${param}`, 'info');
                addOutput(`结果: ${result}`, 'warning');
                addOutput('', 'info');
            }
        });
    }

    // Brute
    if (elements.bruteBtn) {
        elements.bruteBtn.addEventListener('click', async () => {
            const url = elements.bruteInput.value;
            const username = document.getElementById('brute-username').value;
            const wordlist = document.getElementById('brute-wordlist').value || 'passwords.txt';
            if (url && username) {
                const result = await wails.BruteForce(url, username, wordlist);
                addOutput(`[BRUTE] 目标: ${url}`, 'info');
                addOutput(`用户: ${username}`, 'info');
                addOutput(`字典: ${wordlist}`, 'info');
                addOutput(`结果: ${result}`, 'danger');
                addOutput('', 'info');
            }
        });
    }

    // File Scan
    if (elements.fileScanBtn) {
        elements.fileScanBtn.addEventListener('click', async () => {
            const url = elements.fileScanInput.value;
            if (url) {
                const result = await wails.ScanFiles(url);
                addOutput(`[FILESCAN] 目标: ${url}`, 'info');
                addOutput(`结果: ${result}`, 'success');
                addOutput('', 'info');
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
