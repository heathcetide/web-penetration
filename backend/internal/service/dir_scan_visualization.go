package service

import (
	"gorm.io/gorm"
	"sort"
	"strings"
	"web_penetration/internal/model"
)

// 可视化服务
type DirScanVisualizer struct {
	db      *gorm.DB
	service *DirScanService
}

// 时间序列数据点
type TimeSeriesPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	Label     string  `json:"label,omitempty"`
}

// 生成目录树
func (v *DirScanVisualizer) GenerateDirectoryTree(taskID uint) (*DirTreeNode, error) {
	var results []*model.DirScanResult
	if err := v.db.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	root := &DirTreeNode{
		Name: "/",
		Path: "/",
		Type: "directory",
	}

	// 构建树结构
	for _, r := range results {
		v.addToTree(root, r)
	}

	return root, nil
}

// 添加节点到树
func (v *DirScanVisualizer) addToTree(root *DirTreeNode, result *model.DirScanResult) {
	parts := strings.Split(strings.Trim(result.URL, "/"), "/")
	current := root

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
				Name: part,
				Path: strings.Join(parts[:i+1], "/"),
				Type: "file",
			}
			if i < len(parts)-1 || result.IsDir {
				newNode.Type = "directory"
			}
			current.Children = append(current.Children, newNode)
			current = newNode
		}
	}

	// 更新节点信息
	current.Size = result.Length
	current.Count++
	current.Metadata = map[string]interface{}{
		"status_code":  result.StatusCode,
		"content_type": result.ContentType,
		"found_time":   result.Found,
	}
}

// 生成时间序列数据
func (v *DirScanVisualizer) GenerateTimeSeries(taskID uint, metric string) ([]*TimeSeriesPoint, error) {
	var points []*TimeSeriesPoint

	switch metric {
	case "requests":
		// 统计请求数量
		var results []*model.DirScanResult
		if err := v.db.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
			return nil, err
		}

		// 按时间分组
		timeMap := make(map[int64]int)
		for _, r := range results {
			t := r.Found.Unix() - (r.Found.Unix() % 60) // 按分钟分组
			timeMap[t]++
		}

		// 转换为时间序列
		for t, count := range timeMap {
			points = append(points, &TimeSeriesPoint{
				Timestamp: t,
				Value:     float64(count),
				Label:     "requests",
			})
		}

	case "response_time":
		// 统计响应时间
		var metrics []*model.DirScanMetric
		if err := v.db.Where("task_id = ? AND metric_name = ?", taskID, "scan_time").
			Find(&metrics).Error; err != nil {
			return nil, err
		}

		for _, m := range metrics {
			points = append(points, &TimeSeriesPoint{
				Timestamp: m.Timestamp.Unix(),
				Value:     m.MetricValue,
				Label:     "response_time",
			})
		}
	}

	// 排序
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp < points[j].Timestamp
	})

	return points, nil
}
