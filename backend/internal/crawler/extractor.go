package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// ContentExtractor 内容提取器接口
type ContentExtractor interface {
	// Extract 从HTML中提取内容
	Extract(html string) (*ExtractedContent, error)
}

// ExtractedContent 提取的内容
type ExtractedContent struct {
	Title       string
	Description string
	Keywords    []string
	MainContent string
	Images      []ImageInfo
	Links       []LinkInfo
	Scripts     []string
	Styles      []string
}

// ImageInfo 图片信息
type ImageInfo struct {
	URL         string
	Alt         string
	Dimensions  string
	Size        int64
	ContentType string
}

// LinkInfo 链接信息
type LinkInfo struct {
	URL         string
	Text        string
	IsExternal  bool
	ContentType string
	NoFollow    bool
}

// DefaultExtractor 默认内容提取器
type DefaultExtractor struct {
	opts *ExtractorOptions
}

// ExtractorOptions 提取器配置
type ExtractorOptions struct {
	ExtractImages    bool
	ExtractLinks     bool
	ExtractScripts   bool
	ExtractStyles    bool
	MinContentLength int
	ExcludeSelectors []string
}

func NewDefaultExtractor(opts *ExtractorOptions) ContentExtractor {
	if opts == nil {
		opts = &ExtractorOptions{
			ExtractImages:    true,
			ExtractLinks:     true,
			ExtractScripts:   true,
			ExtractStyles:    true,
			MinContentLength: 100,
		}
	}
	return &DefaultExtractor{opts: opts}
}

func (e *DefaultExtractor) Extract(html string) (*ExtractedContent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	content := &ExtractedContent{
		Images:  make([]ImageInfo, 0),
		Links:   make([]LinkInfo, 0),
		Scripts: make([]string, 0),
		Styles:  make([]string, 0),
	}

	// 提取标题和元数据
	content.Title = doc.Find("title").Text()
	content.Description = doc.Find("meta[name=description]").AttrOr("content", "")
	keywords := doc.Find("meta[name=keywords]").AttrOr("content", "")
	if keywords != "" {
		content.Keywords = strings.Split(keywords, ",")
	}

	// 提取主要内容
	mainContent := doc.Find("article, .content, .main, #content, #main").First()
	if mainContent.Length() > 0 {
		// 移除不需要的元素
		mainContent.Find("script, style, iframe, .ad").Remove()
		content.MainContent = strings.TrimSpace(mainContent.Text())
	}

	// 提取图片
	if e.opts.ExtractImages {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists {
				img := ImageInfo{
					URL: src,
					Alt: s.AttrOr("alt", ""),
				}
				content.Images = append(content.Images, img)
			}
		})
	}

	// 提取链接
	if e.opts.ExtractLinks {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				link := LinkInfo{
					URL:      href,
					Text:     s.Text(),
					NoFollow: s.AttrOr("rel", "") == "nofollow",
				}
				content.Links = append(content.Links, link)
			}
		})
	}

	return content, nil
} 