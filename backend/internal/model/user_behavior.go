package model

import (
	"time"
)

// 用户行为趋势
type UserBehaviorTrend struct {
	Date   time.Time `json:"date"`
	Count  int64     `json:"count"`
	Action string    `json:"action"`
}
