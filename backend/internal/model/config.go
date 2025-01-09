package model

import (
	"gorm.io/gorm"
)

// 系统配置
type SystemConfig struct {
	gorm.Model
	Module      string `json:"module" gorm:"size:50"`              // 模块名称
	Key         string `json:"key" gorm:"size:100"`                // 配置键
	Value       string `json:"value" gorm:"type:text"`             // 配置值
	Type        string `json:"type" gorm:"size:20"`                // 值类型
	Description string `json:"description" gorm:"size:255"`         // 描述
	IsSystem    bool   `json:"is_system"`                          // 是否系统配置
	IsReadOnly  bool   `json:"is_read_only"`                       // 是否只读
}

// 配置变更记录
type ConfigChange struct {
	gorm.Model
	ConfigID  uint   `json:"config_id" gorm:"index"`              // 配置ID
	OldValue  string `json:"old_value" gorm:"type:text"`          // 旧值
	NewValue  string `json:"new_value" gorm:"type:text"`          // 新值
	ChangedBy uint   `json:"changed_by"`                          // 修改人
	Reason    string `json:"reason" gorm:"size:255"`              // 修改原因
}