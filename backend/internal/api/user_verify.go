package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"web_penetration/internal/model"
	"web_penetration/internal/service"
)

type UserVerifyHandler struct {
	db           *gorm.DB
	emailService *service.EmailService
}

func NewUserVerifyHandler(db *gorm.DB, emailService *service.EmailService) *UserVerifyHandler {
	return &UserVerifyHandler{
		db:           db,
		emailService: emailService,
	}
}

// 发送邮箱验证码
func (h *UserVerifyHandler) SendEmailVerifyCode(c *gin.Context) {
	userID := c.GetUint("userID")
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	code := service.GenerateVerifyCode()
	if err := h.emailService.SendVerificationEmail(user.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证码失败"})
		return
	}

	// 将验证码保存到缓存中
	if err := service.SaveVerifyCode(user.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存验证码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

// 验证邮箱
func (h *UserVerifyHandler) VerifyEmail(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证验证码
	if !service.VerifyCode(user.Email, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期"})
		return
	}

	// 更新用户验证状态
	user.IsVerified = true
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新验证状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证成功"})
}
