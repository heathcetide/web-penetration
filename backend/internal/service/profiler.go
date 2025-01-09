package service

import (
	"context"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

type ProfileType string

const (
	CPUProfile       ProfileType = "cpu"
	MemoryProfile    ProfileType = "memory"
	GoroutineProfile ProfileType = "goroutine"
	ThreadProfile    ProfileType = "thread"
	BlockProfile     ProfileType = "block"
)

type ProfilerService struct {
	logger     *LoggerService
	enabled    bool
	sampleRate int
	profiles   map[ProfileType]*pprof.Profile
	metrics    map[string]float64
	mu         sync.RWMutex
}

func NewProfilerService(logger *LoggerService) *ProfilerService {
	p := &ProfilerService{
		logger:     logger,
		enabled:    true,
		sampleRate: 100,
		profiles:   make(map[ProfileType]*pprof.Profile),
		metrics:    make(map[string]float64),
	}
	go p.collectMetrics()
	return p
}

// 开始性能分析
func (p *ProfilerService) StartProfiling(ctx context.Context, profileType ProfileType) error {
	if !p.enabled {
		return nil
	}

	switch profileType {
	case CPUProfile:
		return pprof.StartCPUProfile(nil)
	case MemoryProfile:
		runtime.GC()
		return nil
	}
	return nil
}

// 停止性能分析
func (p *ProfilerService) StopProfiling(profileType ProfileType) {
	if !p.enabled {
		return
	}

	switch profileType {
	case CPUProfile:
		pprof.StopCPUProfile()
	case MemoryProfile:
		p.profiles[MemoryProfile] = pprof.Lookup("heap")
	case GoroutineProfile:
		p.profiles[GoroutineProfile] = pprof.Lookup("goroutine")
	}
}

// 收集性能指标
func (p *ProfilerService) collectMetrics() {
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		p.mu.Lock()
		p.metrics["goroutines"] = float64(runtime.NumGoroutine())
		p.metrics["memory_alloc"] = float64(m.Alloc)
		p.metrics["memory_sys"] = float64(m.Sys)
		p.metrics["gc_cycles"] = float64(m.NumGC)
		p.mu.Unlock()

		// 记录性能指标
		p.logger.LogSystem(
			"info",
			"profiler",
			"metrics",
			"Performance metrics collected",
			p.metrics,
		)
	}
}

// 获取性能报告
func (p *ProfilerService) GetProfileReport() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]uint64{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"heap_alloc":  m.HeapAlloc,
			"heap_sys":    m.HeapSys,
			"heap_idle":   m.HeapIdle,
			"heap_inuse":  m.HeapInuse,
		},
		"gc": map[string]uint64{
			"num_gc":      uint64(m.NumGC),
			"pause_total": m.PauseTotalNs,
			"pause_ns":    m.PauseNs[(m.NumGC+255)%256],
		},
		"metrics": p.metrics,
	}
}
