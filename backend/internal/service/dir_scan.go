package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 目录扫描服务
type DirScanService struct {
	DB     *gorm.DB
	client *http.Client
}

// 扫描配置
type DirScanConfig struct {
	Extensions     []string          `json:"extensions"`      // 文件扩展名
	Methods        []string          `json:"methods"`         // HTTP方法
	Headers        map[string]string `json:"headers"`         // 自定义头
	Cookies        string            `json:"cookies"`         // Cookie
	UserAgent      string            `json:"user_agent"`      // UA
	Proxy          string            `json:"proxy"`           // 代理
	SkipSSL        bool              `json:"skip_ssl"`        // 跳过SSL验证
	FollowRedirect bool              `json:"follow_redirect"` // 跟随重定向
	Threads        int               `json:"threads"`         // 线程数
	Delay          int               `json:"delay"`           // 请求延迟(ms)
	Timeout        int               `json:"timeout"`         // 超时时间(s)
	Dictionaries   []string          `json:"dictionaries"`    // 字典列表
	Crawler        *CrawlerConfig    `json:"crawler"`         // 爬虫配置
	VulnScan       bool              `json:"vuln_scan"`       // 是否进行漏洞扫描
	RetryCount     int               `json:"retry_count"`     // 重试次数
	CreatedBy      uint              `json:"created_by"`      // 创建者ID
}

// 创建目录扫描服务
func NewDirScanService(db *gorm.DB) *DirScanService {
	return &DirScanService{
		DB: db,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second,
		},
	}
}

// 创建扫描任务
func (s *DirScanService) CreateScanTask(task *model.DirScanTask) error {
	// 初始化控制通道
	task.CancelChan = make(chan struct{})
	task.PauseChan = make(chan struct{})
	task.ResumeChan = make(chan struct{})

	task.Status = "pending"
	return s.DB.Create(task).Error
}

// 执行扫描任务
func (s *DirScanService) ExecuteScanTask(task *model.DirScanTask) error {
	// 更新任务状态
	task.Status = "running"
	task.StartTime = time.Now()
	if err := s.DB.Save(task).Error; err != nil {
		return err
	}

	// 解析配置
	var config DirScanConfig
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		return s.handleTaskError(task, err)
	}

	// 加载字典
	wordlist, err := s.loadDictionaries(config.Dictionaries)
	if err != nil {
		return s.handleTaskError(task, err)
	}

	// 创建工作池
	pool := &dirScanWorkerPool{
		jobs:      make(chan *dirScanJob, config.Threads*2),
		results:   make(chan *dirScanResult, config.Threads*2),
		workerNum: config.Threads,
		service:   s,
		task:      task,
		config:    &config,
	}

	// 启动工作池
	go pool.start()

	// 启动爬虫
	if config.Crawler != nil {
		crawler := NewDirCrawler(s, config.Crawler)
		go crawler.Crawl(task.Target, pool.results)
	}

	// 提交初始任务
	baseURL := task.Target
	for _, word := range wordlist {
		select {
		case <-task.CancelChan:
			return fmt.Errorf("task cancelled")
		case pool.jobs <- &dirScanJob{
			URL:    path.Join(baseURL, word),
			Method: "GET",
			Depth:  1,
			Parent: baseURL,
		}:
		}
	}

	// 等待扫描完成
	close(pool.jobs)
	pool.wg.Wait()
	close(pool.results)

	// 更新任务状态
	task.Status = "completed"
	task.EndTime = time.Now()
	task.Progress = 100
	return s.DB.Save(task).Error
}

// 启动工作池
func (p *dirScanWorkerPool) start() {
	// 启动工作协程
	for i := 0; i < p.workerNum; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	// 启动结果处理协程
	go p.handleResults()
}

// 工作协程
func (p *dirScanWorkerPool) worker() {
	defer p.wg.Done()

	for job := range p.jobs {
		// 检查任务是否取消
		select {
		case <-p.task.CancelChan:
			return
		default:
		}

		// 执行扫描
		result := p.scan(job)

		// 发送结果
		p.results <- result

		// 如果是目录且需要递归
		if result.IsDir && p.task.Recursive && job.Depth < p.task.MaxDepth {
			// 加载字典
			wordlist, _ := p.service.loadDictionaries(p.config.Dictionaries)

			// 提交子目录任务
			for _, word := range wordlist {
				select {
				case <-p.task.CancelChan:
					return
				case p.jobs <- &dirScanJob{
					URL:    path.Join(job.URL, word),
					Method: job.Method,
					Depth:  job.Depth + 1,
					Parent: job.URL,
				}:
				}
			}
		}

		// 延迟
		if p.config.Delay > 0 {
			time.Sleep(time.Duration(p.config.Delay) * time.Millisecond)
		}
	}
}

// 扫描作业
type dirScanJob struct {
	URL    string
	Method string
	Depth  int
	Parent string
}

// 扫描结果
type dirScanResult struct {
	URL         string          `json:"url"`
	Type        string          `json:"type"` // 结果类型(file/directory/resource)
	StatusCode  int             `json:"status_code"`
	ContentType string          `json:"content_type"`
	Length      int64           `json:"length"`
	Title       string          `json:"title"`
	Found       time.Time       `json:"found"`
	Error       string          `json:"error"`
	Parent      string          `json:"parent"`       // 父目录
	IsDir       bool            `json:"is_dir"`       // 是否是目录
	ScanTime    float64         `json:"scan_time"`    // 扫描耗时
	IsSensitive bool            `json:"is_sensitive"` // 是否敏感文件
	FileType    string          `json:"file_type"`    // 文件类型
	RiskLevel   string          `json:"risk_level"`   // 风险等级
	Fingerprint *WebFingerprint `json:"fingerprint"`  // 指纹信息
	VulnInfo    string          `json:"vuln_info"`    // 漏洞信息(JSON)
}

// 工作池
type dirScanWorkerPool struct {
	jobs      chan *dirScanJob
	results   chan *dirScanResult
	workerNum int
	wg        sync.WaitGroup
	service   *DirScanService
	task      *model.DirScanTask
	config    *DirScanConfig
}

// 行单个URL扫描
func (p *dirScanWorkerPool) scan(job *dirScanJob) *dirScanResult {
	start := time.Now()
	result := &dirScanResult{
		URL:    job.URL,
		Parent: job.Parent,
		Found:  time.Now(),
	}

	// 创建请求
	req, err := http.NewRequest(job.Method, job.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// 设置请求头
	for k, v := range p.config.Headers {
		req.Header.Set(k, v)
	}
	if p.config.UserAgent != "" {
		req.Header.Set("User-Agent", p.config.UserAgent)
	}

	// 发送请求
	resp, err := p.service.client.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// 设置基本信息
	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.Length = resp.ContentLength
	result.IsDir = strings.HasSuffix(job.URL, "/") || result.ContentType == "application/x-directory"
	result.ScanTime = time.Since(start).Seconds()

	// 提取标题
	if strings.Contains(result.ContentType, "text/html") {
		result.Title = p.service.extractTitle(bytes.NewReader(body))
	}

	// 检查敏感文件
	if sensitive, fileType, risk := p.service.checkSensitiveFile(job.URL); sensitive {
		result.IsSensitive = true
		result.FileType = fileType
		result.RiskLevel = risk
	}

	// ��别指纹
	result.Fingerprint = p.service.identifyFingerprint(resp, string(body))

	// 漏洞扫描
	if p.config.VulnScan {
		scanner, _ := NewVulnScanner()
		if vulns := scanner.Scan(result, string(body)); len(vulns) > 0 {
			vulnsJSON, _ := json.Marshal(vulns)
			result.VulnInfo = string(vulnsJSON)
		}
	}

	return result
}

// 处理扫描结果
func (p *dirScanWorkerPool) handleResults() {
	var stats model.DirScanStats
	stats.TaskID = p.task.ID
	stats.StartTime = p.task.StartTime

	for result := range p.results {
		// 保存结果到数据库
		scanResult := &model.DirScanResult{
			TaskID:      p.task.ID,
			URL:         result.URL,
			StatusCode:  result.StatusCode,
			ContentType: result.ContentType,
			Length:      result.Length,
			Title:       result.Title,
			IsDir:       result.IsDir,
			Parent:      result.Parent,
			Found:       time.Now(),
			ScanTime:    result.ScanTime,
		}
		if result.Error != "" {
			scanResult.Error = result.Error
		}

		if err := p.service.DB.Create(scanResult).Error; err != nil {
			continue
		}

		// 更新统计信息
		stats.TotalURLs++
		if result.Error != "" {
			stats.FailedURLs++
			stats.ErrorCount++
		} else {
			stats.SuccessURLs++
		}
		if result.IsDir {
			stats.Directories++
		} else {
			stats.Files++
		}
		stats.AvgResponseTime = (stats.AvgResponseTime*float64(stats.TotalURLs-1) + result.ScanTime) / float64(stats.TotalURLs)
	}

	// 保存统计信息
	stats.EndTime = time.Now()
	stats.Duration = stats.EndTime.Sub(stats.StartTime).Seconds()
	p.service.DB.Create(&stats)
}

// 处理任务错误
func (s *DirScanService) handleTaskError(task *model.DirScanTask, err error) error {
	task.Status = "failed"
	task.Error = err.Error()
	task.EndTime = time.Now()
	return s.DB.Save(task).Error
}

// 加载字典
func (s *DirScanService) loadDictionaries(dictNames []string) ([]string, error) {
	var wordlist []string
	var dicts []model.DirScanDict

	// 从数据库加载指定的字典
	if err := s.DB.Where("name IN ?", dictNames).Find(&dicts).Error; err != nil {
		return nil, err
	}

	// 合并字典内容
	for _, dict := range dicts {
		var words []string
		if err := json.Unmarshal([]byte(dict.Content), &words); err != nil {
			continue
		}
		wordlist = append(wordlist, words...)
	}

	// 去重
	return removeDuplicates(wordlist), nil
}

// 提取HTML标题
func (s *DirScanService) extractTitle(body io.Reader) string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return ""
	}
	return doc.Find("title").First().Text()
}

// 重
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// 敏感文件模式
var sensitivePatterns = []struct {
	Pattern string
	Type    string
	Risk    string
}{
	{`\.git/.*`, "Git Repository", "high"},
	{`\.svn/.*`, "SVN Repository", "high"},
	{`\.env`, "Environment File", "high"},
	{`wp-config\.php`, "WordPress Config", "high"},
	{`config\.php`, "PHP Config", "medium"},
	{`\.htaccess`, "Apache Config", "medium"},
	{`\.bak$`, "Backup File", "medium"},
	{`\.swp$`, "Vim Swap File", "low"},
	{`\.old$`, "Old File", "low"},
	{`\.txt$`, "Text File", "info"},
}

// 检查敏感文件
func (s *DirScanService) checkSensitiveFile(url string) (bool, string, string) {
	for _, pattern := range sensitivePatterns {
		matched, _ := regexp.MatchString(pattern.Pattern, url)
		if matched {
			return true, pattern.Type, pattern.Risk
		}
	}
	return false, "", ""
}

// 网站指纹
type WebFingerprint struct {
	Name    string            `json:"name"`
	Headers map[string]string `json:"headers"`
	Body    []string          `json:"body"`
	Files   []string          `json:"files"`
	Version string            `json:"version"`
}

// 加载指纹库
func (s *DirScanService) loadFingerprints() ([]WebFingerprint, error) {
	var fingerprints []WebFingerprint
	data, err := ioutil.ReadFile("configs/web_fingerprints.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &fingerprints); err != nil {
		return nil, err
	}
	return fingerprints, nil
}

// 识别网站指纹
func (s *DirScanService) identifyFingerprint(resp *http.Response, body string) *WebFingerprint {
	fingerprints, err := s.loadFingerprints()
	if err != nil {
		return nil
	}

	for _, fp := range fingerprints {
		matched := true

		// 检查响应头
		for k, v := range fp.Headers {
			if resp.Header.Get(k) != v {
				matched = false
				break
			}
		}

		// 检查响应体
		for _, pattern := range fp.Body {
			if !strings.Contains(body, pattern) {
				matched = false
				break
			}
		}

		if matched {
			return &fp
		}
	}

	return nil
}

// GetTask 获取任务信息
func (s *DirScanService) GetTask(taskID uint) (*model.DirScanTask, error) {
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// 获取任务状态
func (s *DirScanService) GetTaskStatus(taskID uint) (*model.DirScanTask, error) {
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// 导出结果
func (s *DirScanService) ExportResults(taskID uint, opts *ExportOptions) (*ExportResult, error) {
	// 获取任务结果
	var results []*model.DirScanResult
	if err := s.DB.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	// 根据选项导出
	return &ExportResult{
		TaskID: taskID,
		Data:   results,
		Format: opts.Format,
	}, nil
}

// 获取目录树
func (s *DirScanService) GetDirectoryTree(taskID uint) (*DirTreeNode, error) {
	var results []*model.DirScanResult
	if err := s.DB.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	return buildDirectoryTree(results), nil
}

// 生成报告
func (s *DirScanService) GenerateReport(taskID uint, format string, config map[string]interface{}) (*ReportResult, error) {
	// 获取任务信息和结果
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	return &ReportResult{
		TaskID: taskID,
		URL:    fmt.Sprintf("/reports/%d.%s", taskID, format),
	}, nil
}

// 列出任务
func (s *DirScanService) ListTasks() ([]*model.DirScanTask, error) {
	var tasks []*model.DirScanTask
	if err := s.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// 删除任务
func (s *DirScanService) DeleteTask(taskID uint) error {
	return s.DB.Delete(&model.DirScanTask{}, taskID).Error
}

// 停止任务
func (s *DirScanService) StopTask(taskID uint) error {
	return s.DB.Model(&model.DirScanTask{}).
		Where("id = ?", taskID).
		Update("status", "stopped").Error
}

// 恢复任务
func (s *DirScanService) ResumeTask(taskID uint) error {
	return s.DB.Model(&model.DirScanTask{}).
		Where("id = ?", taskID).
		Update("status", "running").Error
}

// 构建目录树
func buildDirectoryTree(results []*model.DirScanResult) *DirTreeNode {
	root := &DirTreeNode{
		Name:     "/",
		Path:     "/",
		Type:     "directory",
		Children: make([]*DirTreeNode, 0),
	}

	// 构建目录树
	for _, result := range results {
		addToTree(root, result)
	}

	return root
}

// 添加节点到目录树
func addToTree(root *DirTreeNode, result *model.DirScanResult) {
	parts := strings.Split(strings.Trim(result.URL, "/"), "/")
	current := root

	// 遍历路径部分
	for i, part := range parts {
		found := false
		for _, child := range current.Children {
			if child.Name == part {
				current = child
				found = true
				break
			}
		}

		if !found {
			newNode := &DirTreeNode{
				Name:     part,
				Path:     strings.Join(parts[:i+1], "/"),
				Type:     "file",
				Children: make([]*DirTreeNode, 0),
			}

			if i < len(parts)-1 || result.IsDir {
				newNode.Type = "directory"
			}

			current.Children = append(current.Children, newNode)
			current = newNode
		}
	}

	// 设置文件节点的元数据
	if current.Type == "file" {
		current.Size = result.Length
		current.Metadata = map[string]interface{}{
			"content_type": result.ContentType,
			"status_code":  result.StatusCode,
			"found_time":   result.Found,
		}
	}
}

// 创建批量任务
func (s *DirScanService) CreateBatchTask(name string, targets []string, config map[string]interface{}) (*model.DirScanTask, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	task := &model.DirScanTask{
		Name:       name,
		Target:     targets[0], // 使用第一个目标作为主目标
		Targets:    strings.Join(targets, ","),
		Config:     string(configJSON),
		Status:     "pending",
		Progress:   0,
		CreatedBy:  config["created_by"].(uint),
		CancelChan: make(chan struct{}),
		PauseChan:  make(chan struct{}),
		ResumeChan: make(chan struct{}),
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// 获取批量任务
func (s *DirScanService) GetBatchTask(taskID uint) (*model.DirScanTask, error) {
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// 获取批量任务进度
func (s *DirScanService) GetBatchProgress(taskID uint) (*BatchProgress, error) {
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	targets := strings.Split(task.Targets, ",")
	total := len(targets)
	completed := int(task.Progress * float64(total) / 100)

	return &BatchProgress{
		Total:     total,
		Completed: completed,
		Failed:    0, // TODO: 实现失败统计
		Progress:  task.Progress,
		Status:    task.Status,
	}, nil
}

// 获取扫描结果
func (s *DirScanService) GetResults(taskID uint) ([]*model.DirScanResult, error) {
	var results []*model.DirScanResult
	if err := s.DB.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// 获取统计信息
func (s *DirScanService) GetStats(taskID uint) (*model.DirScanStats, error) {
	var stats model.DirScanStats
	if err := s.DB.Where("task_id = ?", taskID).First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

// 获取漏洞信息
func (s *DirScanService) GetVulnerabilities(taskID uint) ([]*model.Vulnerability, error) {
	var vulns []*model.Vulnerability
	if err := s.DB.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return nil, err
	}
	return vulns, nil
}

// 获取性能指标
func (s *DirScanService) GetMetrics(taskID uint) (map[string]interface{}, error) {
	var task model.DirScanTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	// 获取统计信息
	var stats model.DirScanStats
	s.DB.Where("task_id = ?", taskID).First(&stats)

	return map[string]interface{}{
		"request_count":     stats.TotalURLs,
		"success_count":     stats.SuccessURLs,
		"error_count":       stats.ErrorCount,
		"avg_response_time": stats.AvgResponseTime,
		"scan_duration":     stats.Duration,
		"memory_usage":      0, // TODO: 实现内存使用统计
		"cpu_usage":         0, // TODO: 实现CPU使用统计
	}, nil
}
