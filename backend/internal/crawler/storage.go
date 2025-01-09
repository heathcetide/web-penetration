package crawler

import (
	"sync"
)

// Storage 定义存储接口
type Storage interface {
	// Store 存储解析结果
	Store(result *ParseResult) error
	
	// Get 获取指定URL的解析结果
	Get(url string) (*ParseResult, error)
	
	// GetAll 获取所有解析结果
	GetAll() ([]*ParseResult, error)
	
	// Clear 清空存储
	Clear() error
}

// MemoryStorage 实现基于内存的存储
type MemoryStorage struct {
	data sync.Map
}

// NewMemoryStorage 创建内存存储实例
func NewMemoryStorage() Storage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Store(result *ParseResult) error {
	s.data.Store(result.URL, result)
	return nil
}

func (s *MemoryStorage) Get(url string) (*ParseResult, error) {
	if value, ok := s.data.Load(url); ok {
		return value.(*ParseResult), nil
	}
	return nil, ErrNotFound
}

func (s *MemoryStorage) GetAll() ([]*ParseResult, error) {
	var results []*ParseResult
	s.data.Range(func(key, value interface{}) bool {
		results = append(results, value.(*ParseResult))
		return true
	})
	return results, nil
}

func (s *MemoryStorage) Clear() error {
	s.data = sync.Map{}
	return nil
} 