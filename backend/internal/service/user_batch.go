package service

import (
	"gorm.io/gorm"
	"web_penetration/internal/model"
)

type UserBatchService struct {
	db *gorm.DB
}

func NewUserBatchService(db *gorm.DB) *UserBatchService {
	return &UserBatchService{db: db}
}

// 批量更新用户状态
func (s *UserBatchService) BatchUpdateStatus(userIDs []uint, status int) error {
	return s.db.Model(&model.User{}).
		Where("id IN ?", userIDs).
		Update("status", status).Error
}

// 批量分配用户组
func (s *UserBatchService) BatchAssignGroups(userIDs []uint, groupIDs []uint) error {
	members := make([]model.UserGroupMember, 0)
	for _, userID := range userIDs {
		for _, groupID := range groupIDs {
			members = append(members, model.UserGroupMember{
				UserID:  userID,
				GroupID: groupID,
			})
		}
	}
	return s.db.Create(&members).Error
}

// 批量移除用户组
func (s *UserBatchService) BatchRemoveFromGroups(userIDs []uint, groupIDs []uint) error {
	return s.db.Where("user_id IN ? AND group_id IN ?", userIDs, groupIDs).
		Delete(&model.UserGroupMember{}).Error
}

// 批量删除用户
func (s *UserBatchService) BatchDeleteUsers(userIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户组关联
		if err := tx.Where("user_id IN ?", userIDs).
			Delete(&model.UserGroupMember{}).Error; err != nil {
			return err
		}

		// 删除权限关联
		if err := tx.Where("user_id IN ?", userIDs).
			Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		// 删除用户
		return tx.Delete(&model.User{}, userIDs).Error
	})
}
