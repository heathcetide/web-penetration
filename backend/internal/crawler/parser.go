package crawler

import (
	"bytes"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

// ParseResult 表示页面解析结果
type ParseResult struct {
	// 当前页面URL
	URL string
	
	// 发现的链接
	Links []string
	
	// 页面标题
	Title string
	
	// 提取的文本内容
	Text string
	
	// 找到的表单
	Forms []Form
	
	// 发现的资源（图片、脚本、样式表等）
	Resources map[string][]string
	
	// 页面元数据
	Metadata map[string]string
}

// Form 表示HTML表单
type Form struct {
	Action string
	Method string
	Inputs []FormInput
}

// FormInput 表示表单输入字段
type FormInput struct {
	Name     string
	Type     string
	Value    string
	Required bool
}

// Parser 定义页面解析器接口
type Parser interface {
	// Parse 解析页面内容
	Parse(result *FetchResult) (*ParseResult, error)
	
	// SetOptions 设置解析器选项
	SetOptions(opts *ParserOptions)
}

// ParserOptions 定义解析器配置选项
type ParserOptions struct {
	// 是否提取链接
	ExtractLinks bool
	
	// 是否提取表单
	ExtractForms bool
	
	// 是否提取文本
	ExtractText bool
	
	// 是否提取元数据
	ExtractMetadata bool
	
	// 自定义选择器
	CustomSelectors map[string]string
}

// defaultParser 实现默认的解析器
type defaultParser struct {
	opts *ParserOptions
}

// NewDefaultParser 创建默认解析器
func NewDefaultParser() Parser {
	return &defaultParser{
		opts: &ParserOptions{
			ExtractLinks:    true,
			ExtractForms:    true,
			ExtractText:     true,
			ExtractMetadata: true,
			CustomSelectors: make(map[string]string),
		},
	}
}

// Parse 实现解析逻辑
func (p *defaultParser) Parse(result *FetchResult) (*ParseResult, error) {
	doc, err := html.Parse(bytes.NewReader(result.Body))
	if err != nil {
		return nil, err
	}

	parseResult := &ParseResult{
		URL:       result.URL,
		Links:     make([]string, 0),
		Resources: make(map[string][]string),
		Metadata:  make(map[string]string),
	}

	// 解析页面内容
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 提取链接
			if p.opts.ExtractLinks && n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						if link := p.normalizeURL(result.URL, attr.Val); link != "" {
							parseResult.Links = append(parseResult.Links, link)
						}
					}
				}
			}

			// 提取表单
			if p.opts.ExtractForms && n.Data == "form" {
				if form := p.parseForm(n); form != nil {
					parseResult.Forms = append(parseResult.Forms, *form)
				}
			}

			// 提取资源
			switch n.Data {
			case "img":
				p.extractResource(n, "src", "images", parseResult)
			case "script":
				p.extractResource(n, "src", "scripts", parseResult)
			case "link":
				p.extractResource(n, "href", "stylesheets", parseResult)
			}
		}

		// 提取文本
		if p.opts.ExtractText && n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				parseResult.Text += text + " "
			}
		}

		// 递归处理子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return parseResult, nil
}

// SetOptions 设置解析器选项
func (p *defaultParser) SetOptions(opts *ParserOptions) {
	p.opts = opts
}

// 辅助方法：规范化URL
func (p *defaultParser) normalizeURL(base, ref string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	resolvedURL := baseURL.ResolveReference(refURL)
	return resolvedURL.String()
}

// 辅助方法：解析表单
func (p *defaultParser) parseForm(n *html.Node) *Form {
	form := &Form{
		Inputs: make([]FormInput, 0),
	}

	// 获取表单属性
	for _, attr := range n.Attr {
		switch attr.Key {
		case "action":
			form.Action = attr.Val
		case "method":
			form.Method = strings.ToUpper(attr.Val)
		}
	}

	// 解析表单输入字段
	var parseInputs func(*html.Node)
	parseInputs = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "input" {
			input := FormInput{}
			for _, attr := range node.Attr {
				switch attr.Key {
				case "name":
					input.Name = attr.Val
				case "type":
					input.Type = attr.Val
				case "value":
					input.Value = attr.Val
				case "required":
					input.Required = true
				}
			}
			if input.Name != "" {
				form.Inputs = append(form.Inputs, input)
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			parseInputs(c)
		}
	}
	parseInputs(n)

	return form
}

// 辅助方法：提取资源
func (p *defaultParser) extractResource(n *html.Node, attrKey, resourceType string, result *ParseResult) {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			if url := p.normalizeURL(result.URL, attr.Val); url != "" {
				result.Resources[resourceType] = append(result.Resources[resourceType], url)
			}
			break
		}
	}
} 