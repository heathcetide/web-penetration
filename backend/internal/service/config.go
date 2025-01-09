package service

import (
	"errors"
	"gorm.io/gorm"
	"web_penetration/internal/model"
)

type ConfigService struct {
	db     *gorm.DB
	cache  *CacheService
	logger *LoggerService
}

func NewConfigService(db *gorm.DB, cache *CacheService, logger *LoggerService) *ConfigService {
	return &ConfigService{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// 获取配置
func (s *ConfigService) GetConfig(module, key string) (*model.SystemConfig, error) {
	// 先从缓存获取
	cacheKey := "config:" + module + ":" + key
	if value, err := s.cache.Get(cacheKey); err == nil {
		return value.(*model.SystemConfig), nil
	}

	// 从数据库获取
	var config model.SystemConfig
	if err := s.db.Where("module = ? AND key = ?", module, key).First(&config).Error; err != nil {
		return nil, err
	}

	// 写入缓存
	s.cache.Set(cacheKey, &config, 0)
	return &config, nil
}

// 更新配置
func (s *ConfigService) UpdateConfig(userID uint, config *model.SystemConfig, reason string) error {
	// 检查只读配置
	var oldConfig model.SystemConfig
	if err := s.db.First(&oldConfig, config.ID).Error; err != nil {
		return err
	}
	if oldConfig.IsReadOnly {
		return errors.New("cannot modify read-only config")
	}

	// 记录变更
	change := &model.ConfigChange{
		ConfigID:  config.ID,
		OldValue:  oldConfig.Value,
		NewValue:  config.Value,
		ChangedBy: userID,
		Reason:    reason,
	}

	// 使用事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(config).Error; err != nil {
			return err
		}
		if err := tx.Create(change).Error; err != nil {
			return err
		}

		// 清除缓存
		cacheKey := "config:" + config.Module + ":" + config.Key
		s.cache.Delete(cacheKey)

		// 记录审计日志
		s.logger.LogAudit(
			userID,
			"update_config",
			"system_config",
			oldConfig,
			config,
			"",
			"",
		)

		return nil
	})
}

// 批量获取模块配置
func (s *ConfigService) GetModuleConfigs(module string) ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("module = ?", module).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// 初始化系统配置
func (s *ConfigService) InitSystemConfigs(configs []model.SystemConfig) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, config := range configs {
			var exists model.SystemConfig
			if err := tx.Where("module = ? AND key = ?", config.Module, config.Key).
				First(&exists).Error; err == gorm.ErrRecordNotFound {
				config.IsSystem = true
				if err := tx.Create(&config).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
