package model

import (
	"gorm.io/gorm"
	"time"
)

// API日志
type APILog struct {
	gorm.Model
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	Query        string        `json:"query"`
	IP           string        `json:"ip"`
	UserAgent    string        `json:"user_agent"`
	UserID       uint         `json:"user_id"`
	Status       int          `json:"status"`
	Error        string       `json:"error"`
	Latency      time.Duration `json:"latency"`
	RequestAt    time.Time    `json:"request_at"`
	RequestBody  string       `json:"request_body"`
	ResponseBody string       `json:"response_body"`
} 