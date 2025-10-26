package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}

// parsePortRange 解析端口范围
// 支持格式: "1-1000" 或 "80,443,8080"
func parsePortRange(portRange string) ([]int, error) {
	var ports []int

	// 检查是否是逗号分隔的端口列表
	if strings.Contains(portRange, ",") {
		portList := strings.Split(portRange, ",")
		for _, portStr := range portList {
			port, err := strconv.Atoi(strings.TrimSpace(portStr))
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", portStr)
			}
			ports = append(ports, port)
		}
		return ports, nil
	}

	// 检查是否是范围格式 "start-end"
	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port range format: %s", portRange)
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start port: %s", parts[0])
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end port: %s", parts[1])
		}

		// 限制扫描范围，避免过大
		if end-start > 10000 {
			return nil, fmt.Errorf("port range too large (max 10000)")
		}

		for port := start; port <= end; port++ {
			ports = append(ports, port)
		}
		return ports, nil
	}

	// 单个端口
	port, err := strconv.Atoi(strings.TrimSpace(portRange))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %s", portRange)
	}
	return []int{port}, nil
}

// scanSinglePort 扫描单个端口
func scanSinglePort(target string, port int, timeout time.Duration) ScanResult {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		return ScanResult{Port: port, Status: "关闭"}
	}
	defer conn.Close()

	// 端口开放
	return ScanResult{Port: port, Status: "开放"}
}

// ScanPorts 执行端口扫描
func ScanPorts(target, portRange string) ([]ScanResult, error) {
	// 解析端口范围
	ports, err := parsePortRange(portRange)
	if err != nil {
		return nil, fmt.Errorf("解析端口范围失败: %v", err)
	}

	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用通道限制并发数（避免过多连接）
	semaphore := make(chan struct{}, 100) // 最多100个并发连接
	timeout := 2 * time.Second

	fmt.Printf("开始扫描 %s，端口数: %d\n", target, len(ports))

	for _, port := range ports {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			result := scanSinglePort(target, p, timeout)

			// 只记录开放的端口
			if result.Status == "开放" {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()

	if len(results) == 0 {
		return []ScanResult{{Port: 0, Status: "未发现开放端口"}}, nil
	}

	return results, nil
}
