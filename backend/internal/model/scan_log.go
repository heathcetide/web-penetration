package model

import "gorm.io/gorm"

// ScanLog 表示扫描日志
type ScanLog struct {
	gorm.Model
	TaskID   uint   `json:"task_id" gorm:"index"`
	Type     string `json:"type" gorm:"size:50"`
	Level    string `json:"level" gorm:"size:20"`
	Message  string `json:"message" gorm:"type:text"`
	Metadata string `json:"metadata" gorm:"type:text"`
	URL      string `json:"url" gorm:"size:1024"`
}
