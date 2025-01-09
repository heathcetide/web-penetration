package service

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"web_penetration/internal/model"
)

type MFAService struct {
	db *gorm.DB
}

func NewMFAService(db *gorm.DB) *MFAService {
	return &MFAService{db: db}
}

// 生成TOTP密钥
func (s *MFAService) GenerateTOTPSecret() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}

// 启用TOTP
func (s *MFAService) EnableTOTP(userID uint) (*model.MFAMethod, error) {
	secret, err := s.GenerateTOTPSecret()
	if err != nil {
		return nil, err
	}

	method := &model.MFAMethod{
		UserID: userID,
		Type:   "totp",
		Secret: secret,
	}

	if err := s.db.Create(method).Error; err != nil {
		return nil, err
	}

	return method, nil
}

// 验证TOTP
func (s *MFAService) VerifyTOTP(userID uint, code string) (bool, error) {
	var method model.MFAMethod
	if err := s.db.Where("user_id = ? AND type = ?", userID, "totp").
		First(&method).Error; err != nil {
		return false, err
	}

	// TODO: 实现TOTP验证逻辑
	// 这里需要使用如 github.com/pquerna/otp 这样的库来实现

	return true, nil
}

// 生成备用恢复码
func (s *MFAService) GenerateBackupCodes(userID uint) ([]string, error) {
	codes := make([]string, 8)
	for i := range codes {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = fmt.Sprintf("%08x", bytes)
	}

	var method model.MFAMethod
	if err := s.db.Where("user_id = ?", userID).First(&method).Error; err != nil {
		return nil, err
	}

	method.BackupCodes = strings.Join(codes, ",")
	if err := s.db.Save(&method).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

// 发送验证码
func (s *MFAService) SendVerificationCode(userID uint, mfaType string) error {
	var method model.MFAMethod
	if err := s.db.Where("user_id = ? AND type = ?", userID, mfaType).
		First(&method).Error; err != nil {
		return err
	}

	// 使用 crypto/rand 生成安全的随机数
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return err
	}
	// 将随机字节转换为6位数字
	randomInt := uint32(randomBytes[0]) | uint32(randomBytes[1])<<8 |
		uint32(randomBytes[2])<<16 | uint32(randomBytes[3])<<24
	code := fmt.Sprintf("%06d", randomInt%1000000)

	verification := &model.MFAVerification{
		UserID:    userID,
		Type:      mfaType,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.db.Create(verification).Error; err != nil {
		return err
	}

	switch mfaType {
	case "sms":
		// TODO: 实现短信发送
	case "email":
		// TODO: 实现邮件发送
	}

	return nil
}

// 生成密钥
func (s *MFAService) GenerateSecret() string {
	// 实现 TOTP 密钥生成
	return "test_secret" // TODO: 实现实际的密钥生成
}

// 生成二维码
func (s *MFAService) GenerateQRCode(secret string) (string, error) {
	// 实现二维码生成
	return "test_qr_code", nil // TODO: 实现实际的二维码生成
}

// 验证 MFA 代码
func (s *MFAService) VerifyCode(secret, code string) bool {
	// 实现 TOTP 验证
	return true // TODO: 实现实际的代码验证
}
