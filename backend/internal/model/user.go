package model

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"size:50;uniqueIndex" json:"username"`
	Password    string `gorm:"type:text" json:"-"`
	Email       string `gorm:"size:100;uniqueIndex" json:"email"`
	Role        string `gorm:"default:user" json:"role"`
	Avatar      string `gorm:"size:255" json:"avatar"`
	LastLoginAt time.Time `json:"last_login_at"`
	LastLoginIP string    `gorm:"size:50" json:"last_login_ip"`
	Status      int       `gorm:"default:1" json:"status"`
	LoginCount  int       `gorm:"default:0" json:"login_count"`
	Phone       string    `gorm:"size:20" json:"phone"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
}

// 密码加密
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
} 