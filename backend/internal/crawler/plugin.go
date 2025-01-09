package crawler

import (
	"context"
	"plugin"
)

// Plugin 定义爬虫插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string
	
	// OnRequest 请求前处理
	OnRequest(ctx context.Context, req *http.Request) error
	
	// OnResponse 响应后处理
	OnResponse(ctx context.Context, resp *http.Response) error
	
	// OnError 错误处理
	OnError(ctx context.Context, err error)
}

// PluginManager 插件管理器
type PluginManager struct {
	plugins []Plugin
	mutex   sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
	}
}

// LoadPlugin 加载插件
func (pm *PluginManager) LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return err
	}

	plugin, ok := symPlugin.(Plugin)
	if !ok {
		return errors.New("invalid plugin type")
	}

	pm.mutex.Lock()
	pm.plugins = append(pm.plugins, plugin)
	pm.mutex.Unlock()

	return nil
}

// ExecuteOnRequest 执行所有插件的OnRequest
func (pm *PluginManager) ExecuteOnRequest(ctx context.Context, req *http.Request) error {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	for _, p := range pm.plugins {
		if err := p.OnRequest(ctx, req); err != nil {
			return err
		}
	}
	return nil
} 