package scan

import (
    "context"
    "sync"
    "time"
)

// VulnEngine 漏洞检测引擎
type VulnEngine struct {
    config      *ScanConfig
    rules       map[string][]*VulnRule
    detector    *VulnDetector
    analyzer    *VulnAnalyzer
    processor   ResultProcessor
    rateLimiter RateLimiter
    
    tasks       chan *VulnTask
    results     chan *VulnResult
    
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc
}

// VulnTask 漏洞检测任务
type VulnTask struct {
    ID          string       `json:"id"`
    Target      string       `json:"target"`
    Service     *ServiceInfo `json:"service"`
    Rules       []*VulnRule  `json:"rules"`
    StartTime   time.Time    `json:"start_time"`
    EndTime     time.Time    `json:"end_time"`
    Status      string       `json:"status"`
    Error       error        `json:"error,omitempty"`
}

// NewVulnEngine 创建漏洞检测引擎
func NewVulnEngine(config *ScanConfig) *VulnEngine {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &VulnEngine{
        config:      config,
        rules:       make(map[string][]*VulnRule),
        detector:    NewVulnDetector(),
        analyzer:    NewVulnAnalyzer(),
        rateLimiter: NewRateLimiter(ctx, config.RateLimit),
        tasks:       make(chan *VulnTask, config.BatchSize),
        results:     make(chan *VulnResult, config.BatchSize),
        ctx:         ctx,
        cancel:      cancel,
    }
}

// Start 启动漏洞检测引擎
func (e *VulnEngine) Start() error {
    // 启动工作协程
    for i := 0; i < e.config.Concurrency; i++ {
        e.wg.Add(1)
        go e.worker()
    }
    
    // 启动结果处理协程
    e.wg.Add(1)
    go e.processResults()
    
    return nil
}

// Stop 停止漏洞检测引擎
func (e *VulnEngine) Stop() {
    e.cancel()
    e.wg.Wait()
}

// worker 工作协程
func (e *VulnEngine) worker() {
    defer e.wg.Done()
    
    for {
        select {
        case <-e.ctx.Done():
            return
        case task := <-e.tasks:
            // 等待速率限制
            if err := e.rateLimiter.Wait(); err != nil {
                continue
            }
            
            // 执行漏洞检测
            results := e.detectVulns(task)
            
            // 发送结果
            for _, result := range results {
                e.results <- result
            }
        }
    }
}

// detectVulns 执行漏洞检测
func (e *VulnEngine) detectVulns(task *VulnTask) []*VulnResult {
    var results []*VulnResult
    
    // 获取适用的规则
    rules := e.getApplicableRules(task.Service)
    
    // 对每个规则执行检测
    for _, rule := range rules {
        result := e.detector.DetectVuln(task.Target, task.Service, rule)
        if result != nil {
            results = append(results, result)
        }
    }
    
    return results
}

// processResults 处理检测结果
func (e *VulnEngine) processResults() {
    defer e.wg.Done()
    
    for {
        select {
        case <-e.ctx.Done():
            return
        case result := <-e.results:
            // 分析结果
            e.analyzer.AddVuln(result)
            
            // 处理结果
            if err := e.processor.ProcessVuln(result); err != nil {
                // TODO: 处理错误
                continue
            }
            
            // 触发自动响应
            if e.shouldTriggerResponse(result) {
                e.triggerAutoResponse(result)
            }
        }
    }
}

// getApplicableRules 获取适用的规则
func (e *VulnEngine) getApplicableRules(service *ServiceInfo) []*VulnRule {
    var rules []*VulnRule
    
    // 获取服务特定规则
    if serviceRules, ok := e.rules[service.Name]; ok {
        rules = append(rules, serviceRules...)
    }
    
    // 获取通用规则
    if commonRules, ok := e.rules["common"]; ok {
        rules = append(rules, commonRules...)
    }
    
    return rules
}

// shouldTriggerResponse 判断是否需要触发自动响应
func (e *VulnEngine) shouldTriggerResponse(result *VulnResult) bool {
    // 根据漏洞严重程度判断
    switch result.Severity {
    case "critical", "high":
        return true
    default:
        return false
    }
}

// triggerAutoResponse 触发自动响应
func (e *VulnEngine) triggerAutoResponse(result *VulnResult) {
    // TODO: 实现自动响应逻辑
    // 1. 发送告警通知
    // 2. 执行应急处置
    // 3. 更新安全策略
}

// AddTask 添加检测任务
func (e *VulnEngine) AddTask(task *VulnTask) error {
    select {
    case e.tasks <- task:
        return nil
    default:
        return ErrTaskQueueFull
    }
} 