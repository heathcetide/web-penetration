package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

// Renderer 定义页面渲染器接口
type Renderer interface {
	// Render 渲染页面
	Render(url string) (string, error)
	// Close 关闭渲染器
	Close() error
}

// ChromeRenderer 实现基于Chrome的渲染器
type ChromeRenderer struct {
	ctx    context.Context
	cancel context.CancelFunc
	opts   *RendererOptions
}

// RendererOptions 渲染器配置选项
type RendererOptions struct {
	Timeout       time.Duration
	WaitForLoad   bool
	ScrollToBottom bool
	UserAgent     string
	Proxy         string
}

// NewChromeRenderer 创建Chrome渲染器
func NewChromeRenderer(opts *RendererOptions) (*ChromeRenderer, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.NoSandbox,
		chromedp.Headless,
	}

	if opts.UserAgent != "" {
		options = append(options, chromedp.UserAgent(opts.UserAgent))
	}
	if opts.Proxy != "" {
		options = append(options, chromedp.ProxyServer(opts.Proxy))
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel = chromedp.NewContext(ctx)

	return &ChromeRenderer{
		ctx:    ctx,
		cancel: cancel,
		opts:   opts,
	}, nil
}

// Render 实现页面渲染
func (r *ChromeRenderer) Render(url string) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.opts.Timeout)
	defer cancel()

	var html string
	tasks := []chromedp.Action{
		chromedp.Navigate(url),
	}

	if r.opts.WaitForLoad {
		tasks = append(tasks, chromedp.WaitReady("body", chromedp.ByQuery))
	}

	if r.opts.ScrollToBottom {
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := chromedp.Evaluate(`
				window.scrollTo(0, document.body.scrollHeight);
				new Promise(resolve => setTimeout(resolve, 1000));
			`, nil).Do(ctx)
			return err
		}))
	}

	tasks = append(tasks, chromedp.OuterHTML("html", &html))

	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		return "", err
	}

	return html, nil
}

func (r *ChromeRenderer) Close() error {
	r.cancel()
	return nil
} 