package model

import (
	"gorm.io/gorm"
	"time"
)

// 设备指纹
type DeviceFingerprint struct {
	gorm.Model
	UserID         uint   `gorm:"index" json:"user_id"`
	FingerprintID  string `gorm:"size:64;uniqueIndex" json:"fingerprint_id"`
	UserAgent      string `gorm:"size:255" json:"user_agent"`
	OS             string `gorm:"size:50" json:"os"`
	Browser        string `gorm:"size:50" json:"browser"`
	ScreenRes      string `gorm:"size:20" json:"screen_res"`
	ColorDepth     string `gorm:"size:10" json:"color_depth"`
	Language       string `gorm:"size:20" json:"language"`
	Timezone       string `gorm:"size:50" json:"timezone"`
	WebGLRenderer  string `gorm:"size:255" json:"webgl_renderer"`
	Canvas         string `gorm:"size:32" json:"canvas"`         // Canvas指纹
	Fonts          string `gorm:"type:text" json:"fonts"`       // 已安装字体列表
	Plugins        string `gorm:"type:text" json:"plugins"`     // 浏览器插件列表
	IsTrusted     bool   `gorm:"default:false" json:"is_trusted"`
	LastSeenAt    time.Time `json:"last_seen_at"`
	TrustScore    float64   `json:"trust_score"`
}

// 地理位置信息
type GeoLocation struct {
	gorm.Model
	IP            string  `gorm:"size:50;index" json:"ip"`
	Country       string  `gorm:"size:50" json:"country"`
	Region        string  `gorm:"size:50" json:"region"`
	City          string  `gorm:"size:50" json:"city"`
	ISP           string  `gorm:"size:100" json:"isp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Timezone      string  `gorm:"size:50" json:"timezone"`
	ASN           string  `gorm:"size:20" json:"asn"`           // 自治系统编号
	Organization  string  `gorm:"size:100" json:"organization"` // ISP组织名称
	IsDataCenter  bool    `json:"is_data_center"`              // 是否为数据中心IP
	IsProxy       bool    `json:"is_proxy"`                    // 是否为代理IP
	IsVPN         bool    `json:"is_vpn"`                      // 是否为VPN
	ThreatLevel   int     `json:"threat_level"`                // 威胁等级
}

// 用户常用地理位置
type UserLocation struct {
	gorm.Model
	UserID      uint    `gorm:"index" json:"user_id"`
	Country     string  `gorm:"size:50" json:"country"`
	Region      string  `gorm:"size:50" json:"region"`
	City        string  `gorm:"size:50" json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Frequency   int     `json:"frequency"`    // 访问频率
	LastSeenAt  time.Time `json:"last_seen_at"`
	IsTrusted   bool    `json:"is_trusted"`
}

// 设备登录历史
type DeviceLoginHistory struct {
	gorm.Model
	UserID          uint      `gorm:"index" json:"user_id"`
	FingerprintID   string    `gorm:"size:64;index" json:"fingerprint_id"`
	IP              string    `gorm:"size:50" json:"ip"`
	LoginTime       time.Time `json:"login_time"`
	LoginStatus     string    `gorm:"size:20" json:"login_status"` // success, failed
	RiskScore       float64   `json:"risk_score"`
	GeoLocationID   uint      `gorm:"index" json:"geo_location_id"`
	AuthMethod      string    `gorm:"size:20" json:"auth_method"` // password, mfa, sso
	SessionDuration int       `json:"session_duration"`          // 会话时长(秒)
} 