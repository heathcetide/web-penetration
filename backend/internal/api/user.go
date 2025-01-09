package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"web_penetration/internal/model"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var users []model.User
	var total int64
	query := h.db.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// 获取单个用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password,omitempty"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新基本信息
	user.Username = updateData.Username
	user.Email = updateData.Email
	user.Role = updateData.Role

	// 如果提供了新密码，则更新密码
	if updateData.Password != "" {
		if err := user.HashPassword(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
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

	if !user.CheckPassword(req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码错误"})
		return
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// 批量更新状态
func (h *UserHandler) BatchUpdateStatus(c *gin.Context) {
	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
		Status  int    `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&model.User{}).Where("id IN ?", req.UserIDs).
		Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 批量分配用户组
func (h *UserHandler) BatchAssignGroups(c *gin.Context) {
	var req struct {
		UserIDs  []uint `json:"user_ids" binding:"required"`
		GroupIDs []uint `json:"group_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members := make([]model.UserGroupMember, 0)
	for _, userID := range req.UserIDs {
		for _, groupID := range req.GroupIDs {
			members = append(members, model.UserGroupMember{
				UserID:  userID,
				GroupID: groupID,
			})
		}
	}

	if err := h.db.Create(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量分配用户组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分配成功"})
}

// 批量移除用户组
func (h *UserHandler) BatchRemoveFromGroups(c *gin.Context) {
	var req struct {
		UserIDs  []uint `json:"user_ids" binding:"required"`
		GroupIDs []uint `json:"group_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Where("user_id IN ? AND group_id IN ?", req.UserIDs, req.GroupIDs).
		Delete(&model.UserGroupMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量移除用户组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除成功"})
}

// 批量删除用户
func (h *UserHandler) BatchDeleteUsers(c *gin.Context) {
	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户组关联
		if err := tx.Where("user_id IN ?", req.UserIDs).
			Delete(&model.UserGroupMember{}).Error; err != nil {
			return err
		}

		// 删除权限关联
		if err := tx.Where("user_id IN ?", req.UserIDs).
			Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		// 删除用户
		return tx.Delete(&model.User{}, req.UserIDs).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 导出用户
func (h *UserHandler) ExportUsers(c *gin.Context) {
	var users []model.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户数据失败"})
		return
	}

	// TODO: 实现导出逻辑
	c.JSON(http.StatusOK, users)
}

// 导入用户
func (h *UserHandler) ImportUsers(c *gin.Context) {
	// TODO: 实现导入逻辑
	c.JSON(http.StatusOK, gin.H{"message": "导入成功"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []model.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var update struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", userID).Updates(update).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "更新成功"})
}
