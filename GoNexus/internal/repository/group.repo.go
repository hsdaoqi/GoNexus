package repository

import (
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/pkg/global"

	"gorm.io/gorm"
)

// CreateGroup 创建群组 (事务：建群 + 加群主)
func CreateGroup(group *model.Group) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建群组基本信息
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		// 2. 自动把创建者加入群成员表 (身份: 3=群主)
		member := model.GroupMember{
			GroupID: group.ID, // 这里会自动拿到刚才插入生成的 ID
			UserID:  group.OwnerID,
			Role:    3, // 群主
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetJoinedGroups 获取我加入的群列表 (用于 SidePanel 展示)
func GetJoinedGroups(userID uint) ([]model.Group, error) {
	var groups []model.Group
	// 联表查询：从 group_members 表反查 groups 表
	// 使用表别名 g / gm，避免 MySQL 对保留字 `groups` 解析出错
	err := global.DB.Table("`groups` AS g").
		Select("g.*, gm.unread_count").
		Joins("JOIN group_members AS gm ON gm.group_id = g.id").
		Where("gm.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// GetGroupByID 查询群
func GetGroupByID(id uint) (*model.Group, error) {
	var g model.Group
	if err := global.DB.First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

// SaveGroup 保存群信息
func SaveGroup(g *model.Group) error {
	return global.DB.Save(g).Error
}

// GetGroupMembers 获取群成员列表
func GetGroupMembers(groupID uint) ([]dto.GroupMemberResponse, error) {
	var members []dto.GroupMemberResponse
	// 联表查询获取用户头像和昵称
	err := global.DB.Table("group_members").
		Select("group_members.*, users.avatar as user_avatar, users.nickname as user_name").
		Joins("JOIN users ON users.id = group_members.user_id").
		Where("group_members.group_id = ?", groupID).
		Find(&members).Error
	return members, err
}

// AddGroupMember 添加群成员
func AddGroupMember(member *model.GroupMember) error {
	return global.DB.Create(member).Error
}

// CheckGroupMember 检查是否已经是成员
func CheckGroupMember(groupID, userID uint) bool {
	var count int64
	global.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	return count > 0
}

// GetGroupMemberIDs 获取群成员的用户ID列表
func GetGroupMemberIDs(groupID uint) ([]uint, error) {
	var ids []uint
	err := global.DB.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &ids).Error
	return ids, err
}

// RemoveGroupMember 移除成员
func RemoveGroupMember(groupID, userID uint) error {
	// Unscoped 硬删除，彻底移除
	return global.DB.Unscoped().Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&model.GroupMember{}).Error
}

// UpdateMemberMuteStatus 更新禁言状态
func UpdateMemberMuteStatus(groupID, userID uint, muted int) error {
	return global.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("muted", muted).Error
}

// GetMemberRole 获取成员角色
func GetMemberRole(groupID, userID uint) (int, error) {
	var member model.GroupMember
	err := global.DB.Select("role").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error
	return member.Role, err
}

// IsMemberMuted 检查成员是否被禁言
func IsMemberMuted(groupID, userID uint) (bool, error) {
	var member model.GroupMember
	err := global.DB.Select("muted").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error
	if err != nil {
		return false, err
	}
	return member.Muted == 1, nil
}

// UpdateMemberRole 更新成员角色
func UpdateMemberRole(groupID, userID uint, role int) error {
	return global.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", role).Error
}

// TransferGroupOwner 转让群主 (事务)
func TransferGroupOwner(groupID, oldOwnerID, newOwnerID uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新群表 owner_id
		if err := tx.Model(&model.Group{}).Where("id = ?", groupID).Update("owner_id", newOwnerID).Error; err != nil {
			return err
		}

		// 2. 更新旧群主角色 -> 普通成员 (1)
		if err := tx.Model(&model.GroupMember{}).
			Where("group_id = ? AND user_id = ?", groupID, oldOwnerID).
			Update("role", 1).Error; err != nil {
			return err
		}

		// 3. 更新新群主角色 -> 群主 (3)
		if err := tx.Model(&model.GroupMember{}).
			Where("group_id = ? AND user_id = ?", groupID, newOwnerID).
			Update("role", 3).Error; err != nil {
			return err
		}

		return nil
	})
}

// IncrementGroupUnread 增加群未读数
func IncrementGroupUnread(groupID, senderID uint) error {
	return global.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id != ?", groupID, senderID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
}

// ClearGroupUnread 清除群未读数
func ClearGroupUnread(userID, groupID uint) error {
	return global.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("unread_count", 0).Error
}
