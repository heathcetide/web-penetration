package scan

import (
	"fmt"
	"strconv"
	"strings"
)

// PortRange 端口范围
type PortRange struct {
	Start int
	End   int
}

// ParsePortRange 解析端口范围
func ParsePortRange(rangeStr string) ([]int, error) {
	var ports []int
	
	// 处理多个端口范围，用逗号分隔
	ranges := strings.Split(rangeStr, ",")
	for _, r := range ranges {
		// 处理单个范围
		if strings.Contains(r, "-") {
			// 处理范围格式 (例如: 80-100)
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range format: %s", r)
			}
			
			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, err
			}
			
			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, err
			}
			
			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			// 处理单个端口
			port, err := strconv.Atoi(strings.TrimSpace(r))
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
	}
	
	return ports, nil
} 