package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/oschwald/maxminddb-golang"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
	"web_penetration/internal/model"
)

type GeoIPService struct {
	db *maxminddb.Reader
}

type DeviceAnalysisService struct {
	db           *gorm.DB
	geoIPService *GeoIPService // 第三方地理位置服务
}

func NewDeviceAnalysisService(db *gorm.DB, geoIPService *GeoIPService) *DeviceAnalysisService {
	return &DeviceAnalysisService{
		db:           db,
		geoIPService: geoIPService,
	}
}

// 生成设备指纹
func (s *DeviceAnalysisService) GenerateFingerprint(data map[string]interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonStr)
	return hex.EncodeToString(hash[:]), nil
}

// 分析设备可信度
func (s *DeviceAnalysisService) AnalyzeDeviceTrust(fingerprint *model.DeviceFingerprint, ip string) (float64, error) {
	score := 100.0

	// 检查设备历史
	var loginCount int64
	if err := s.db.Model(&model.DeviceLoginHistory{}).
		Where("fingerprint_id = ? AND login_status = ?", fingerprint.FingerprintID, "success").
		Count(&loginCount).Error; err != nil {
		return 0, err
	}

	if loginCount == 0 {
		score -= 30 // 新设备扣分
	}

	// 检查地理位置异常
	geoLocation, err := s.geoIPService.GetGeoLocation(ip)
	if err != nil {
		return 0, err
	}

	if geoLocation.IsProxy || geoLocation.IsVPN {
		score -= 20
	}

	if geoLocation.IsDataCenter {
		score -= 10
	}

	// 检查用户常用位置
	var userLocations []model.UserLocation
	if err := s.db.Where("user_id = ?", fingerprint.UserID).
		Order("frequency DESC").
		Limit(5).
		Find(&userLocations).Error; err != nil {
		return 0, err
	}

	locationTrusted := false
	for _, loc := range userLocations {
		if s.isNearLocation(geoLocation.Latitude, geoLocation.Longitude, loc.Latitude, loc.Longitude) {
			locationTrusted = true
			break
		}
	}

	if !locationTrusted {
		score -= 15
	}

	// 检查设备环境
	if s.hasAnomalousEnvironment(fingerprint) {
		score -= 25
	}

	return math.Max(0, score), nil
}

// 检查位置是否接近
func (s *DeviceAnalysisService) isNearLocation(lat1, lon1, lat2, lon2 float64) bool {
	// 使用Haversine公式计算两点间距离
	const earthRadius = 6371.0 // 地球半径，单位km
	const maxDistance = 100.0  // 最大允许距离，单位km

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance <= maxDistance
}

// 检查设备环境异常
func (s *DeviceAnalysisService) hasAnomalousEnvironment(fingerprint *model.DeviceFingerprint) bool {
	// 检查WebGL渲染器是否为虚拟设备
	if s.isVirtualGPU(fingerprint.WebGLRenderer) {
		return true
	}

	// 检查字体列表是否异常
	if s.hasAnomalousFonts(fingerprint.Fonts) {
		return true
	}

	// 检查插件列表是否异常
	if s.hasAnomalousPlugins(fingerprint.Plugins) {
		return true
	}

	return false
}

// 更新设备登录历史
func (s *DeviceAnalysisService) RecordDeviceLogin(userID uint, fingerprintID string, ip string, status string) error {
	geoLocation, err := s.geoIPService.GetGeoLocation(ip)
	if err != nil {
		return err
	}

	// 保存地理位置信息
	if err := s.db.Create(geoLocation).Error; err != nil {
		return err
	}

	// 记录登录历史
	history := &model.DeviceLoginHistory{
		UserID:        userID,
		FingerprintID: fingerprintID,
		IP:            ip,
		LoginTime:     time.Now(),
		LoginStatus:   status,
		GeoLocationID: geoLocation.ID,
	}

	if err := s.db.Create(history).Error; err != nil {
		return err
	}

	// 更新用户常用位置
	return s.updateUserLocation(userID, geoLocation)
}

// 更新用户常用位置
func (s *DeviceAnalysisService) updateUserLocation(userID uint, geoLocation *model.GeoLocation) error {
	var location model.UserLocation
	err := s.db.Where("user_id = ? AND country = ? AND region = ? AND city = ?",
		userID, geoLocation.Country, geoLocation.Region, geoLocation.City).
		First(&location).Error

	if err == gorm.ErrRecordNotFound {
		location = model.UserLocation{
			UserID:    userID,
			Country:   geoLocation.Country,
			Region:    geoLocation.Region,
			City:      geoLocation.City,
			Latitude:  geoLocation.Latitude,
			Longitude: geoLocation.Longitude,
			Frequency: 1,
		}
		return s.db.Create(&location).Error
	}

	location.Frequency++
	location.LastSeenAt = time.Now()
	return s.db.Save(&location).Error
}

// 获取设备访问统计
func (s *DeviceAnalysisService) GetDeviceStats(userID uint) (map[string]interface{}, error) {
	var stats struct {
		TotalDevices     int64
		TrustedDevices   int64
		UnusualLocations int64
		RiskyLogins      int64
	}

	// 统计设备数量
	if err := s.db.Model(&model.DeviceFingerprint{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalDevices).Error; err != nil {
		return nil, err
	}

	// 统计可信设备
	if err := s.db.Model(&model.DeviceFingerprint{}).
		Where("user_id = ? AND is_trusted = ?", userID, true).
		Count(&stats.TrustedDevices).Error; err != nil {
		return nil, err
	}

	// 统计异常位置登录
	if err := s.db.Model(&model.DeviceLoginHistory{}).
		Joins("JOIN geo_locations ON device_login_histories.geo_location_id = geo_locations.id").
		Where("device_login_histories.user_id = ? AND (geo_locations.is_proxy = ? OR geo_locations.is_vpn = ?)",
			userID, true, true).
		Count(&stats.UnusualLocations).Error; err != nil {
		return nil, err
	}

	// 统计风险登录
	if err := s.db.Model(&model.DeviceLoginHistory{}).
		Where("user_id = ? AND risk_score < ?", userID, 60).
		Count(&stats.RiskyLogins).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_devices":     stats.TotalDevices,
		"trusted_devices":   stats.TrustedDevices,
		"unusual_locations": stats.UnusualLocations,
		"risky_logins":      stats.RiskyLogins,
		"trust_rate":        float64(stats.TrustedDevices) / float64(stats.TotalDevices),
		"risk_login_rate":   float64(stats.RiskyLogins) / float64(stats.TotalDevices),
	}, nil
}

// 检查是否为虚拟GPU
func (s *DeviceAnalysisService) isVirtualGPU(renderer string) bool {
	virtualGPUs := []string{
		"VMware", "VirtualBox", "QEMU", "Microsoft Basic Display",
		"llvmpipe", "SwiftShader", "VirtIO",
	}

	for _, gpu := range virtualGPUs {
		if strings.Contains(strings.ToLower(renderer), strings.ToLower(gpu)) {
			return true
		}
	}
	return false
}

// 检查字体列表是否异常
func (s *DeviceAnalysisService) hasAnomalousFonts(fonts string) bool {
	var fontList []string
	json.Unmarshal([]byte(fonts), &fontList)

	// 检查字体数量是否异常
	if len(fontList) < 10 || len(fontList) > 500 {
		return true
	}

	// 检查是否缺少常见字体
	commonFonts := []string{"Arial", "Times New Roman", "Courier New"}
	commonFontCount := 0
	for _, font := range fontList {
		for _, common := range commonFonts {
			if strings.Contains(font, common) {
				commonFontCount++
				break
			}
		}
	}

	return commonFontCount < 2
}

// 检查插件列表是否异常
func (s *DeviceAnalysisService) hasAnomalousPlugins(plugins string) bool {
	var pluginList []string
	json.Unmarshal([]byte(plugins), &pluginList)

	// 检查插件数量是否异常
	if len(pluginList) == 0 || len(pluginList) > 50 {
		return true
	}

	// 检查是否有可疑插件
	suspiciousPlugins := []string{
		"proxy", "vpn", "anonymizer", "debugger",
		"automation", "recorder", "bot", "scraper",
	}

	for _, plugin := range pluginList {
		for _, suspicious := range suspiciousPlugins {
			if strings.Contains(strings.ToLower(plugin), suspicious) {
				return true
			}
		}
	}

	return false
}

// GeoIP服务方法
func (s *GeoIPService) GetGeoLocation(ip string) (*model.GeoLocation, error) {
	// TODO: 实现IP地理位置查询
	return &model.GeoLocation{
		IP:          ip,
		Country:     "Unknown",
		City:        "Unknown",
		ISP:         "Unknown",
		IsProxy:     false,
		IsVPN:       false,
		ThreatLevel: 0,
	}, nil
}
