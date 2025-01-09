package service

import (
	"gorm.io/gorm"
	"time"
)

// VariableManager 处理工作流变量的管理
type VariableManager struct {
	variables map[string]interface{}
	db        *gorm.DB
	cache     *CacheService
	logger    *LoggerService
}

// NewVariableManager 创建一个新的变量管理器实例
func NewVariableManager(db *gorm.DB, cache *CacheService, logger *LoggerService) *VariableManager {
	return &VariableManager{
		variables: make(map[string]interface{}),
		db:        db,
		cache:     cache,
		logger:    logger,
	}
}

// SetVariable 设置变量值
func (vm *VariableManager) SetVariable(name string, value interface{}) {
	vm.variables[name] = value
	// 添加日志记录
	vm.logger.LogSystem("INFO", "VariableManager", "SetVariable", 
		"设置变量: "+name, map[string]interface{}{
			"name": name,
			"value": value,
		})
	
	// 添加缓存操作，设置1小时过期时间
	vm.cache.Set(name, value, time.Hour)
}

// GetVariable 获取变量值
func (vm *VariableManager) GetVariable(name string) (interface{}, bool) {
	// 首先尝试从缓存获取
	if value, ok := vm.cache.Get(name); ok == nil {
		return value, true
	}
	
	value, ok := vm.variables[name]
	return value, ok
}

// DeleteVariable 删除变量
func (vm *VariableManager) DeleteVariable(name string) {
	delete(vm.variables, name)
	vm.cache.Delete(name)
	vm.logger.LogSystem("INFO", "VariableManager", "DeleteVariable", 
		"删除变量: "+name, map[string]interface{}{
			"name": name,
		})
}

// ClearVariables 清除所有变量
func (vm *VariableManager) ClearVariables() {
	vm.variables = make(map[string]interface{})
	// 如果 CacheService 没有 Clear 方法，我们可以删除这行
	// vm.cache.Clear()
	vm.logger.LogSystem("INFO", "VariableManager", "ClearVariables", 
		"清除所有变量", nil)
}
