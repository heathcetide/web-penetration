package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"web_penetration/internal/model"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

// 创建用户组
func (h *GroupHandler) CreateGroup(c *gin.Context) {
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

// 更新用户组
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	var group model.UserGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户组失败"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// 删除用户���
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.UserGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 分配权限
func (h *GroupHandler) AssignPermissions(c *gin.Context) {
	var req struct {
		GroupID       uint   `json:"group_id" binding:"required"`
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 删除旧的权限
	if err := h.db.Where("group_id = ?", req.GroupID).Delete(&model.GroupPermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除旧权限失败"})
		return
	}

	// 添加新的权限
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

func (h *GroupHandler) ListGroups(c *gin.Context) {
	var groups []model.UserGroup
	if err := h.db.Find(&groups).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, groups)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	id := c.Param("id")
	var group model.UserGroup
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "组不存在"})
		return
	}
	c.JSON(200, group)
}

func (h *GroupHandler) AddUsers(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的组ID"})
		return
	}

	var req struct {
		UserIDs []uint `json:"user_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var members []model.UserGroupMember
	for _, userID := range req.UserIDs {
		members = append(members, model.UserGroupMember{
			UserID:  userID,
			GroupID: uint(groupID),
		})
	}

	if err := h.db.Create(&members).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "添加成功"})
}

// 从用户组中移除用户
func (h *GroupHandler) RemoveUsers(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的组ID"})
		return
	}

	var req struct {
		UserIDs []uint `json:"user_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 删除用户组成员关系
	if err := h.db.Where("group_id = ? AND user_id IN ?", groupID, req.UserIDs).
		Delete(&model.UserGroupMember{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "移除用户失败"})
		return
	}

	c.JSON(200, gin.H{"message": "移除成功"})
}
