package service

import (
	"encoding/json"
	"sync"
)

// WorkflowVariable 工作流变量
type WorkflowVariable struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
}

// WorkflowVariableService 工作流变量服务
type WorkflowVariableService struct {
	variables map[string]*WorkflowVariable
	mu        sync.RWMutex
}

// NewWorkflowVariableService 创建工作流变量服务
func NewWorkflowVariableService() *WorkflowVariableService {
	return &WorkflowVariableService{
		variables: make(map[string]*WorkflowVariable),
	}
}

// SetVariable 设置变量
func (s *WorkflowVariableService) SetVariable(name string, value interface{}, varType string, desc string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.variables[name] = &WorkflowVariable{
		Name:        name,
		Value:       value,
		Type:        varType,
		Description: desc,
	}
}

// GetVariable 获取变量
func (s *WorkflowVariableService) GetVariable(name string) (*WorkflowVariable, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.variables[name]
	return v, ok
}

// DeleteVariable 删除变量
func (s *WorkflowVariableService) DeleteVariable(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.variables, name)
}

// ExportVariables 导出所有变量
func (s *WorkflowVariableService) ExportVariables() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return json.Marshal(s.variables)
}
