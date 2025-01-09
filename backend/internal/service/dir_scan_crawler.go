package service

import (
    "net/url"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

// 爬虫配置
type CrawlerConfig struct {
    MaxDepth     int      `json:"max_depth"`      // 最大爬取深度
    IncludeJS    bool     `json:"include_js"`     // 是否包含JS文件
    IncludeCSS   bool     `json:"include_css"`    // 是否包含CSS文件
    ExcludePaths []string `json:"exclude_paths"`  // 排除路径
    AllowDomains []string `json:"allow_domains"`  // 允许的域名
}

// 爬虫服务
type DirCrawler struct {
    service *DirScanService
    config  *CrawlerConfig
    visited map[string]bool
    depth   map[string]int
}

// 创建爬虫
func NewDirCrawler(service *DirScanService, config *CrawlerConfig) *DirCrawler {
    return &DirCrawler{
        service: service,
        config:  config,
        visited: make(map[string]bool),
        depth:   make(map[string]int),
    }
}

// 爬取URL
func (c *DirCrawler) Crawl(baseURL string, results chan<- *dirScanResult) {
    if c.visited[baseURL] {
        return
    }
    c.visited[baseURL] = true

    // 检查深度
    depth := c.depth[baseURL]
    if depth >= c.config.MaxDepth {
        return
    }

    // 发送请求
    resp, err := c.service.client.Get(baseURL)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    // 解析HTML
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return
    }

    // 提取链接
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        if !exists {
            return
        }

        // 解析URL
        u, err := url.Parse(href)
        if err != nil {
            return
        }

        // 处理相对路径
        if !u.IsAbs() {
            base, _ := url.Parse(baseURL)
            u = base.ResolveReference(u)
        }

        // 检查域名
        if !c.isAllowedDomain(u.Host) {
            return
        }

        // 检查排除路径
        if c.isExcludedPath(u.Path) {
            return
        }

        // 递归爬取
        nextURL := u.String()
        c.depth[nextURL] = depth + 1
        c.Crawl(nextURL, results)
    })

    // 提取资源链接
    if c.config.IncludeJS {
        c.extractResources(doc, "script[src]", "src", results)
    }
    if c.config.IncludeCSS {
        c.extractResources(doc, "link[rel=stylesheet]", "href", results)
    }
}

// 提取资源链接
func (c *DirCrawler) extractResources(doc *goquery.Document, selector, attr string, results chan<- *dirScanResult) {
    doc.Find(selector).Each(func(i int, s *goquery.Selection) {
        if href, exists := s.Attr(attr); exists {
            if u, err := url.Parse(href); err == nil {
                results <- &dirScanResult{
                    URL:  u.String(),
                    Type: "resource",
                }
            }
        }
    })
}

// 检查域名是否允许
func (c *DirCrawler) isAllowedDomain(domain string) bool {
    if len(c.config.AllowDomains) == 0 {
        return true
    }
    for _, d := range c.config.AllowDomains {
        if strings.HasSuffix(domain, d) {
            return true
        }
    }
    return false
}

// 检查路径是否排除
func (c *DirCrawler) isExcludedPath(path string) bool {
    for _, p := range c.config.ExcludePaths {
        if strings.HasPrefix(path, p) {
            return true
        }
    }
    return false
}
 