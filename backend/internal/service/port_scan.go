package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"web_penetration/internal/model"
)

type PortScanService struct {
	db *gorm.DB
}

func NewPortScanService(db *gorm.DB) *PortScanService {
	return &PortScanService{db: db}
}

// 创建扫描任务
func (s *PortScanService) CreateScanTask(task *model.ScanTask) error {
	task.Status = "pending"
	return s.db.Create(task).Error
}

// 执行扫描任务
func (s *PortScanService) ExecuteScanTask(task *model.ScanTask) error {
	// 更新任务状态
	task.Status = "running"
	task.StartTime = time.Now()
	if err := s.db.Save(task).Error; err != nil {
		return err
	}

	// 解析目标
	var targets []string
	if task.Targets != "" {
		if err := json.Unmarshal([]byte(task.Targets), &targets); err != nil {
			return s.handleTaskError(task, err)
		}
	} else {
		targets = []string{task.Target}
	}

	// 解析端口范围
	ports, err := s.parsePortRange(task.PortRange)
	if err != nil {
		return s.handleTaskError(task, err)
	}

	// 创建工作池
	jobs := make(chan *ScanJob, len(targets)*len(ports))
	results := make(chan *model.ScanResult, len(targets)*len(ports))
	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < task.Concurrency; i++ {
		wg.Add(1)
		go s.scanWorker(jobs, results, &wg, task.Timeout)
	}

	// 提交扫描任务
	totalJobs := 0
	for _, target := range targets {
		for _, port := range ports {
			jobs <- &ScanJob{
				TaskID:   task.ID,
				Target:   target,
				IP:       target, // 这里可以添加DNS解析
				Port:     port,
				Protocol: strings.ToLower(task.ScanType),
			}
			totalJobs++
		}
	}
	close(jobs)

	// 处理结果
	go func() {
		wg.Wait()
		close(results)
	}()

	// 保存结果
	completedJobs := 0
	for result := range results {
		if err := s.db.Create(result).Error; err != nil {
			return err
		}
		completedJobs++
		task.Progress = float64(completedJobs) / float64(totalJobs) * 100
		s.db.Save(task)
	}

	// 完成任务
	task.Status = "completed"
	task.EndTime = time.Now()
	task.Progress = 100
	return s.db.Save(task).Error
}

// 扫描任务作业
type ScanJob struct {
	TaskID   uint   `json:"task_id"`
	Target   string `json:"target"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// 扫描工作协程
func (s *PortScanService) scanWorker(jobs <-chan *ScanJob, results chan<- *model.ScanResult, wg *sync.WaitGroup, timeout int) {
	defer wg.Done()

	for job := range jobs {
		result := &model.ScanResult{
			TaskID:   job.TaskID,
			IP:       job.IP,
			Port:     job.Port,
			Protocol: job.Protocol,
		}

		// 执行端口扫描
		state, banner := s.scanPort(job.IP, job.Port, job.Protocol, timeout)
		result.Status = state
		result.Banner = banner

		// 如果端口开放，进行服务识别
		if state == "open" {
			service, version := s.identifyService(banner, job.Port, job.Protocol)
			result.Service = service
			result.Version = version

			// 获取指纹信息
			fingerprint := s.getFingerprint(banner, service)
			result.Fingerprint = fingerprint

			// 检查相关漏洞
			vulns := s.checkVulnerabilities(service, version)
			vulnsJSON, _ := json.Marshal(vulns)
			result.RawData = string(vulnsJSON)
		}

		results <- result
	}
}

// 执行端口扫描
func (s *PortScanService) scanPort(ip string, port int, protocol string, timeout int) (state string, banner string) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(protocol, address, time.Duration(timeout)*time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			return "closed", ""
		}
		return "filtered", ""
	}
	defer conn.Close()

	// 尝试获取banner
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	banner = string(buffer[:n])

	return "open", banner
}

// 获取指纹信息
func (s *PortScanService) getFingerprint(banner string, service string) string {
	var fingerprint struct {
		Banner     string   `json:"banner"`
		Keywords   []string `json:"keywords"`
		Signatures []string `json:"signatures"`
	}

	// TODO: 实现指纹识别逻辑
	// 1. 提取banner特征
	// 2. 匹配指纹库
	// 3. 识别技术栈

	fingerprintJSON, _ := json.Marshal(fingerprint)
	return string(fingerprintJSON)
}

// 解析目标
func (s *PortScanService) parseTargets(targets string, targetType string) ([]string, error) {
	var result []string
	switch targetType {
	case "ip":
		result = strings.Split(targets, ",")
	case "domain":
		// TODO: 实现域名解析
	case "cidr":
		// TODO: 实现CIDR解析
	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
	return result, nil
}

// 解析端口范围
func (s *PortScanService) parsePortRange(portRange string) ([]int, error) {
	var ports []int
	ranges := strings.Split(portRange, ",")

	for _, r := range ranges {
		if strings.Contains(r, "-") {
			// 处理端口范围，如：80-100
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}

			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}

			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}

			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			// 处理单个端口
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

// 处理任务错误
func (s *PortScanService) handleTaskError(task *model.ScanTask, err error) error {
	task.Status = "failed"
	task.Error = err.Error()
	task.EndTime = time.Now()
	return s.db.Save(task).Error
}

// 获取扫描统计信息
func (s *PortScanService) GetScanStats(taskID uint) (map[string]interface{}, error) {
	var stats struct {
		TotalPorts     int64
		OpenPorts      int64
		ClosedPorts    int64
		FilteredPorts  int64
		CommonServices map[string]int64
		TopVulns       []string
	}

	stats.CommonServices = make(map[string]int64)

	// 统计端口状态
	if err := s.db.Model(&model.ScanResult{}).
		Where("task_id = ?", taskID).
		Count(&stats.TotalPorts).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.ScanResult{}).
		Where("task_id = ? AND state = ?", taskID, "open").
		Count(&stats.OpenPorts).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.ScanResult{}).
		Where("task_id = ? AND state = ?", taskID, "closed").
		Count(&stats.ClosedPorts).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.ScanResult{}).
		Where("task_id = ? AND state = ?", taskID, "filtered").
		Count(&stats.FilteredPorts).Error; err != nil {
		return nil, err
	}

	// 统计常见服务
	rows, err := s.db.Model(&model.ScanResult{}).
		Select("service, count(*) as count").
		Where("task_id = ? AND state = ?", taskID, "open").
		Group("service").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var service string
		var count int64
		if err := rows.Scan(&service, &count); err != nil {
			return nil, err
		}
		stats.CommonServices[service] = count
	}

	// TODO: 分析Top漏洞

	return map[string]interface{}{
		"total_ports":     stats.TotalPorts,
		"open_ports":      stats.OpenPorts,
		"closed_ports":    stats.ClosedPorts,
		"filtered_ports":  stats.FilteredPorts,
		"common_services": stats.CommonServices,
		"top_vulns":       stats.TopVulns,
		"open_rate":       float64(stats.OpenPorts) / float64(stats.TotalPorts),
		"filtered_rate":   float64(stats.FilteredPorts) / float64(stats.TotalPorts),
	}, nil
}

// 使用配置执行扫描
func (s *PortScanService) ScanWithConfig(target string, configStr string, results chan<- *model.ScanResult, cancelChan <-chan struct{}) error {
	var config ScanConfig
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return err
	}

	// 创建工作池
	jobs := make(chan *ScanJob, 1)
	var wg sync.WaitGroup

	// 启动工作协程
	wg.Add(1)
	go s.scanWorker(jobs, results, &wg, int(config.Timeout.Seconds()))

	// 提交扫描任务
	select {
	case <-cancelChan:
		close(jobs)
		wg.Wait()
		return fmt.Errorf("scan canceled")
	case jobs <- &ScanJob{
		Target:   target,
		Port:     config.Port,
		Protocol: config.ScanType,
	}:
	}

	close(jobs)
	wg.Wait()
	return nil
}

// TCP连接扫描
func (s *PortScanService) scanPortConnect(ip string, port int, timeout int) (state string, banner string) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			return "closed", ""
		}
		return "filtered", ""
	}
	defer conn.Close()

	// 尝试获取banner
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	banner = string(buffer[:n])

	return "open", banner
}
