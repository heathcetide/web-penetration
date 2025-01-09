package service

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"unicode"
	"web_penetration/internal/model"
)

type SecurityPolicyService struct {
	db *gorm.DB
}

func NewSecurityPolicyService(db *gorm.DB) *SecurityPolicyService {
	return &SecurityPolicyService{db: db}
}

// 验证密码强度
func (s *SecurityPolicyService) ValidatePassword(userID uint, password string) error {
	var policy model.SecurityPolicy
	if err := s.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", userID).
		First(&policy).Error; err != nil {
		// 使用默认策略
		if err := s.db.Where("is_default = ?", true).First(&policy).Error; err != nil {
			return err
		}
	}

	if len(password) < policy.MinPasswordLength {
		return fmt.Errorf("密码长度不能小于%d位", policy.MinPasswordLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return fmt.Errorf("密码必须包含大写字母")
	}
	if policy.RequireLowercase && !hasLower {
		return fmt.Errorf("密码必须包含小写字母")
	}
	if policy.RequireNumbers && !hasNumber {
		return fmt.Errorf("密码必须包含数字")
	}
	if policy.RequireSpecialChars && !hasSpecial {
		return fmt.Errorf("密码必须包含特殊字符")
	}

	// 检查密码历史
	if policy.PreventPasswordReuse > 0 {
		var histories []model.PasswordHistory
		if err := s.db.Where("user_id = ?", userID).
			Order("created_at DESC").
			Limit(policy.PreventPasswordReuse).
			Find(&histories).Error; err != nil {
			return err
		}

		for _, history := range histories {
			if history.HashedPassword == password {
				return fmt.Errorf("不能使用最近%d次使用过的密码", policy.PreventPasswordReuse)
			}
		}
	}

	return nil
}

// 检查登录尝试
func (s *SecurityPolicyService) CheckLoginAttempts(userID uint, ip string) error {
	var policy model.SecurityPolicy
	if err := s.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", userID).
		First(&policy).Error; err != nil {
		// 使用默认策略
		if err := s.db.Where("is_default = ?", true).First(&policy).Error; err != nil {
			return err
		}
	}

	// 检查IP限制
	if policy.RestrictedIPs != "" {
		restrictedIPs := strings.Split(policy.RestrictedIPs, ",")
		for _, restrictedIP := range restrictedIPs {
			if strings.TrimSpace(restrictedIP) == ip {
				return fmt.Errorf("IP地址被限制")
			}
		}
	}

	// 检查登录失败次数
	var failCount int64
	if err := s.db.Model(&model.LoginAttempt{}).
		Where("user_id = ? AND status = ? AND created_at > ?",
			userID, false, time.Now().Add(-policy.LockoutDuration)).
		Count(&failCount).Error; err != nil {
		return err
	}

	if int(failCount) >= policy.MaxLoginAttempts {
		return fmt.Errorf("登录失败次数过多，账户已被锁定%v", policy.LockoutDuration)
	}

	return nil
}

// 检查会话数量
func (s *SecurityPolicyService) CheckConcurrentSessions(userID uint) error {
	var policy model.SecurityPolicy
	if err := s.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", userID).
		First(&policy).Error; err != nil {
		// 使用默认策略
		if err := s.db.Where("is_default = ?", true).First(&policy).Error; err != nil {
			return err
		}
	}

	var activeSessions int64
	if err := s.db.Model(&model.UserSession{}).
		Where("user_id = ? AND is_active = ? AND expires_at > ?",
			userID, true, time.Now()).
		Count(&activeSessions).Error; err != nil {
		return err
	}

	if int(activeSessions) >= policy.ConcurrentSessions {
		return fmt.Errorf("已达到最大同时在线数量限制: %d", policy.ConcurrentSessions)
	}

	return nil
}

// 检查密码过期
func (s *SecurityPolicyService) CheckPasswordExpiration(userID uint) error {
	var policy model.SecurityPolicy
	if err := s.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", userID).
		First(&policy).Error; err != nil {
		// 使用默认策略
		if err := s.db.Where("is_default = ?", true).First(&policy).Error; err != nil {
			return err
		}
	}

	if policy.PasswordExpireDays > 0 {
		var lastPasswordChange model.PasswordHistory
		if err := s.db.Where("user_id = ?", userID).
			Order("created_at DESC").
			First(&lastPasswordChange).Error; err != nil {
			return err
		}

		expirationDate := lastPasswordChange.CreatedAt.AddDate(0, 0, policy.PasswordExpireDays)
		if time.Now().After(expirationDate) {
			return fmt.Errorf("密码已过期，请修改密码")
		}
	}

	return nil
}

// 应用安全策略
func (s *SecurityPolicyService) ApplySecurityPolicy(userID uint, policyID uint) error {
	userPolicy := &model.UserSecurityPolicy{
		UserID:           userID,
		SecurityPolicyID: policyID,
	}
	return s.db.Create(userPolicy).Error
}

// 创建或更新安全策略
func (s *SecurityPolicyService) SaveSecurityPolicy(policy *model.SecurityPolicy) error {
	if policy.ID == 0 {
		return s.db.Create(policy).Error
	}
	return s.db.Save(policy).Error
}

// 获取用户的安全策略
func (s *SecurityPolicyService) GetUserSecurityPolicy(userID uint) (*model.SecurityPolicy, error) {
	var policy model.SecurityPolicy
	err := s.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", userID).
		First(&policy).Error
	if err != nil {
		// 使用默认策略
		err = s.db.Where("is_default = ?", true).First(&policy).Error
	}
	return &policy, err
}
