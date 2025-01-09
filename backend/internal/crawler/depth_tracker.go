package crawler

import (
	"sync"
)

// DepthTracker URL深度追踪器
type DepthTracker struct {
	depths map[string]int
	mutex  sync.RWMutex
}

func NewDepthTracker() *DepthTracker {
	return &DepthTracker{
		depths: make(map[string]int),
	}
}

func (t *DepthTracker) GetDepth(url string) int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.depths[url]
}

func (t *DepthTracker) SetDepth(url string, depth int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.depths[url] = depth
}

func (t *DepthTracker) IsExceeded(url string, maxDepth int) bool {
	if maxDepth <= 0 {
		return false
	}
	return t.GetDepth(url) >= maxDepth
} 