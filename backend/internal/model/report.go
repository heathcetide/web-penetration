package model

import (
	"gorm.io/gorm"
	"time"
)

// 报告
type Report struct {
	gorm.Model
	Name        string          `json:"name" gorm:"size:100"`
	Type        string          `json:"type" gorm:"size:50"` // scan, audit, analysis
	Status      string          `json:"status" gorm:"size:50"`
	Template    ReportTemplate  `json:"-" gorm:"foreignKey:TemplateID"`
	TemplateID  uint            `json:"template_id" gorm:"index"`
	Sections    []ReportSection `json:"sections" gorm:"foreignKey:ReportID"`
	Content     string          `json:"content" gorm:"type:text"`
	GeneratedAt time.Time       `json:"generated_at"`
}

// 报告章节
type ReportSection struct {
	gorm.Model
	ReportID uint   `json:"report_id" gorm:"index"`
	Title    string `json:"title" gorm:"size:100"`
	Content  string `json:"content" gorm:"type:text"`
	Order    int    `json:"order"`
	Type     string `json:"type" gorm:"size:50"` // text, chart, table
}
