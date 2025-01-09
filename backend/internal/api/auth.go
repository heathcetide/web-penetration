package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
	"web_penetration/internal/model"
	"web_penetration/internal/service"
)

type AuthHandler struct {
	authService    *service.AuthService
	db             *gorm.DB
	sessionService *service.SessionService
	redis          *redis.Client
	mfaService     *service.MFAService
}

func NewAuthHandler(authService *service.AuthService, db *gorm.DB, sessionService *service.SessionService, redis *redis.Client, mfaService *service.MFAService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		db:             db,
		sessionService: sessionService,
		redis:          redis,
		mfaService:     mfaService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "user",
	}

	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := h.db.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 生成 token
	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	// 更新用户登录信息
	user.LastLoginAt = time.Now()
	user.LastLoginIP = c.ClientIP()
	user.LoginCount++
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 创建会话
	if err := h.sessionService.CreateSession(user.ID, token, c.ClientIP(), c.Request.UserAgent()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 记录登录日志
	userLog := &model.UserLog{
		UserID:    user.ID,
		Action:    "login",
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    true,
		Message:   "登录成功",
	}
	h.db.Create(userLog)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查是否需要MFA
	var policy model.SecurityPolicy
	if err := h.db.Joins("JOIN user_security_policies usp ON usp.security_policy_id = security_policies.id").
		Where("usp.user_id = ?", user.ID).
		First(&policy).Error; err == nil && policy.RequireMFA {

		// 检查是否已经设置MFA
		var mfaMethod model.MFAMethod
		if err := h.db.Where("user_id = ? AND is_verified = ?", user.ID, true).
			First(&mfaMethod).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":              "需要设置多因素认证",
				"mfa_required":       true,
				"mfa_setup_required": true,
			})
			return
		}

		// 生成临时token用于MFA验证
		tempToken := uuid.New().String()
		ctx := context.Background()
		if err := h.redis.Set(ctx, fmt.Sprintf("mfa_temp_token:%s", tempToken), user.ID, 5*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"mfa_required": true,
			"temp_token":   tempToken,
			"mfa_type":     mfaMethod.Type,
		})
		return
	}

	// 生成 token
	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	// 更新用户登录信息
	user.LastLoginAt = time.Now()
	user.LastLoginIP = c.ClientIP()
	user.LoginCount++
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 创建会话
	if err := h.sessionService.CreateSession(user.ID, token, c.ClientIP(), c.Request.UserAgent()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 记录登录日志
	userLog := &model.UserLog{
		UserID:    user.ID,
		Action:    "login",
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    true,
		Message:   "登录成功",
	}
	h.db.Create(userLog)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// 验证MFA
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	var req struct {
		TempToken string `json:"temp_token" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从Redis获取用户ID
	ctx := context.Background()
	userID, err := h.redis.Get(ctx, fmt.Sprintf("mfa_temp_token:%s", req.TempToken)).Uint64()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的临���token"})
		return
	}

	// 验证MFA代码
	verified, err := h.mfaService.VerifyTOTP(uint(userID), req.Code)
	if err != nil || !verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "验证码无效"})
		return
	}

	// 删除临时token
	h.redis.Del(ctx, fmt.Sprintf("mfa_temp_token:%s", req.TempToken))

	// 生成正式token
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	// 更新用户登录信息
	user.LastLoginAt = time.Now()
	user.LastLoginIP = c.ClientIP()
	user.LoginCount++
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 创建会话
	if err := h.sessionService.CreateSession(user.ID, token, c.ClientIP(), c.Request.UserAgent()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 记录登录日志
	userLog := &model.UserLog{
		UserID:    user.ID,
		Action:    "login",
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    true,
		Message:   "登录成功",
	}
	h.db.Create(userLog)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 验证刷新令牌
	claims, err := h.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "无效的刷新令牌"})
		return
	}

	// 生成新的访问令牌
	token, err := h.authService.GenerateToken(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.sessionService.InvalidateSession(token); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "已登出"})
}

func (h *AuthHandler) SetupMFA(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 生成 MFA 密钥
	secret := h.mfaService.GenerateSecret()
	qrCode, err := h.mfaService.GenerateQRCode(secret)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 保存 MFA 配置
	mfa := &model.MFAMethod{
		UserID: userID,
		Type:   "totp",
		Secret: secret,
	}
	if err := h.db.Create(mfa).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"secret":  secret,
		"qr_code": qrCode,
	})
}
