package crawler

import (
	"net/url"
	"path"
	"strings"
)

// URLFilter URL过滤器接口
type URLFilter interface {
	// Filter 过滤URL
	Filter(url string) bool
}

// CompositeURLFilter 组合过滤器
type CompositeURLFilter struct {
	filters []URLFilter
}

// FileExtensionFilter 文件扩展名过滤器
type FileExtensionFilter struct {
	allowedExts map[string]bool
}

// DomainFilter 域名过滤器
type DomainFilter struct {
	allowedDomains map[string]bool
	allowSubdomains bool
}

// PathFilter 路径过滤器
type PathFilter struct {
	allowedPaths []string
	excludedPaths []string
}

func NewFileExtensionFilter(exts []string) *FileExtensionFilter {
	allowed := make(map[string]bool)
	for _, ext := range exts {
		allowed[strings.ToLower(ext)] = true
	}
	return &FileExtensionFilter{allowedExts: allowed}
}

func (f *FileExtensionFilter) Filter(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	ext := strings.ToLower(path.Ext(u.Path))
	if ext == "" {
		return true
	}
	return f.allowedExts[ext]
}

func NewDomainFilter(domains []string, allowSubdomains bool) *DomainFilter {
	allowed := make(map[string]bool)
	for _, domain := range domains {
		allowed[strings.ToLower(domain)] = true
	}
	return &DomainFilter{
		allowedDomains: allowed,
		allowSubdomains: allowSubdomains,
	}
}

func (f *DomainFilter) Filter(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	domain := strings.ToLower(u.Hostname())
	
	if f.allowSubdomains {
		for allowedDomain := range f.allowedDomains {
			if strings.HasSuffix(domain, allowedDomain) {
				return true
			}
		}
		return false
	}
	
	return f.allowedDomains[domain]
} 