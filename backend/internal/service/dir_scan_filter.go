package service

import (
	"strings"
	"web_penetration/internal/model"
)

// 过滤结果
func FilterResults(results []*model.DirScanResult, filter *ResultFilter) []*model.DirScanResult {
	var filtered []*model.DirScanResult

	for _, result := range results {
		if !matchesFilter(result, filter) {
			continue
		}
		filtered = append(filtered, result)
	}

	return filtered
}

// 检查是否匹配过滤条件
func matchesFilter(result *model.DirScanResult, filter *ResultFilter) bool {
	// 检查状态码
	if len(filter.StatusCodes) > 0 {
		matched := false
		for _, code := range filter.StatusCodes {
			if result.StatusCode == code {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查类型
	if len(filter.Types) > 0 {
		matched := false
		for _, t := range filter.Types {
			if result.Type == t {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查深度
	if filter.MinDepth > 0 && result.Depth < filter.MinDepth {
		return false
	}
	if filter.MaxDepth > 0 && result.Depth > filter.MaxDepth {
		return false
	}

	// 检查关键词
	if len(filter.Keywords) > 0 {
		matched := false
		for _, kw := range filter.Keywords {
			if strings.Contains(result.URL, kw) || strings.Contains(result.Title, kw) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查排除目录
	for _, dir := range filter.ExcludeDirs {
		if strings.HasPrefix(result.URL, dir) {
			return false
		}
	}

	return true
}
