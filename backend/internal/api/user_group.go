package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"web_penetration/internal/model"
)

type UserGroupHandler struct {
	db *gorm.DB
}

func NewUserGroupHandler(db *gorm.DB) *UserGroupHandler {
	return &UserGroupHandler{db: db}
}

// 创建用户组
func (h *UserGroupHandler) CreateGroup(c *gin.Context) {
	var group model.UserGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户组失败"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// 分配权限
func (h *UserGroupHandler) AssignPermissions(c *gin.Context) {
	var req struct {
		GroupID       uint   `json:"group_id" binding:"required"`
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissions := make([]model.GroupPermission, 0)
	for _, permID := range req.PermissionIDs {
		permissions = append(permissions, model.GroupPermission{
			GroupID:      req.GroupID,
			PermissionID: permID,
		})
	}

	if err := h.db.Create(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分配权限失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限分配成功"})
}
