package crawler

import (
	"regexp"
	"sync"
)

// ResultFilter 结果过滤器接口
type ResultFilter interface {
	// Filter 过滤结果
	Filter(result *ParseResult) bool
	
	// Reset 重置过滤器
	Reset()
}

// CompositeFilter 组合过滤器
type CompositeFilter struct {
	filters []ResultFilter
}

// RegexFilter 正则过滤器
type RegexFilter struct {
	patterns []*regexp.Regexp
}

// DomainFilter 域名过滤器
type DomainFilter struct {
	allowedDomains map[string]bool
	mutex          sync.RWMutex
}

func NewCompositeFilter(filters ...ResultFilter) *CompositeFilter {
	return &CompositeFilter{filters: filters}
}

func (f *CompositeFilter) Filter(result *ParseResult) bool {
	for _, filter := range f.filters {
		if !filter.Filter(result) {
			return false
		}
	}
	return true
}

func NewRegexFilter(patterns []string) (*RegexFilter, error) {
	var regexps []*regexp.Regexp
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexps = append(regexps, re)
	}
	return &RegexFilter{patterns: regexps}, nil
}

func (f *RegexFilter) Filter(result *ParseResult) bool {
	for _, pattern := range f.patterns {
		if pattern.MatchString(result.URL) {
			return true
		}
	}
	return false
} 