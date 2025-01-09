package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
	"web_penetration/internal/model"
)

// 导出服务
type DirScanExporter struct {
	db      *gorm.DB
	service *DirScanService
}

// 导出任务结果
func (e *DirScanExporter) Export(taskID uint, opts *ExportOptions) (string, error) {
	// 获取结果
	var results []*model.DirScanResult
	query := e.db.Where("task_id = ?", taskID)

	// 应用时间范围
	if opts.TimeRange != "" {
		duration, err := time.ParseDuration(opts.TimeRange)
		if err == nil {
			query = query.Where("created_at >= ?", time.Now().Add(-duration))
		}
	}

	if err := query.Find(&results).Error; err != nil {
		return "", err
	}

	// 应用过滤器
	if opts.Filter != "" {
		var filter ResultFilter
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return "", fmt.Errorf("invalid filter: %v", err)
		}
		results = FilterResults(results, &filter)
	}

	// 根据格式导出
	switch opts.Format {
	case ExportFormatJSON:
		return e.exportJSON(results)
	case ExportFormatCSV:
		return e.exportCSV(results)
	case ExportFormatHTML:
		return e.exportHTML(results)
	default:
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

// 导出为JSON
func (e *DirScanExporter) exportJSON(results []*model.DirScanResult) (string, error) {
	filename := fmt.Sprintf("scan_results_%d.json", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return "", err
	}

	return filename, nil
}

// 导出为CSV
func (e *DirScanExporter) exportCSV(results []*model.DirScanResult) (string, error) {
	filename := fmt.Sprintf("scan_results_%d.csv", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{"URL", "Status", "Type", "Title", "Length", "Found", "Error"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// 写入数据
	for _, r := range results {
		row := []string{
			r.URL,
			strconv.Itoa(r.StatusCode),
			r.Type,
			r.Title,
			strconv.FormatInt(r.Length, 10),
			r.Found.Format(time.RFC3339),
			r.Error,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return filename, nil
}

// 导出HTML报告
func (e *DirScanExporter) exportHTML(results []*model.DirScanResult) (string, error) {
	// TODO: 实现HTML报告生成
	return "", nil
}
