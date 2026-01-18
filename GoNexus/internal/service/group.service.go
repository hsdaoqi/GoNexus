package service

import (
	"errors"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
)

// CreateGroup 业务逻辑
func CreateGroup(ownerID uint, name, avatar, notice string) (*model.Group, error) {
	group := &model.Group{
		Name:    name,
		OwnerID: ownerID,
		Avatar:  avatar,
		Notice:  notice,
		Type:    1, // 默认为普通群 (Type=2 是公共大厅)
	}

	// 如果没传头像，给个随机默认头像 (DiceBear API)
	if group.Avatar == "" {
		group.Avatar = "https://api.dicebear.com/7.x/identicon/svg?seed=" + name
	}

	if err := repository.CreateGroup(group); err != nil {
		return nil, err
	}
	return group, nil
}

// GetJoinedGroups 业务逻辑
func GetJoinedGroups(userID uint) ([]model.Group, error) {
	return repository.GetJoinedGroups(userID)
}

// UpdateGroup 仅群主可更新群名称/头像/公告
func UpdateGroup(ownerID, groupID uint, name, avatar, notice string) error {
	g, err := repository.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if g.OwnerID != ownerID {
		return errors.New("仅群主可编辑群资料")
	}
	// 按需更新
	if name != "" {
		g.Name = name
	}
	if avatar != "" {
		g.Avatar = avatar
	}
	if notice != "" {
		g.Notice = notice
	}
	return repository.SaveGroup(g)
}

// GetGroupMembers 获取群成员
func GetGroupMembers(groupID uint) ([]dto.GroupMemberResponse, error) {
	return repository.GetGroupMembers(groupID)
}

// InviteFriendToGroup 邀请好友入群
func InviteFriendToGroup(groupID, userID, friendID uint) error {
	// 1. 验证是否是群成员 (只有群成员能邀请)
	if !repository.CheckGroupMember(groupID, userID) {
		return errors.New("你不是该群成员，无法邀请")
	}

	// 2. 验证目标是否已经在群里
	if repository.CheckGroupMember(groupID, friendID) {
		return errors.New("该用户已经在群里了")
	}

	// 3. 入库
	member := &model.GroupMember{
		GroupID: groupID,
		UserID:  friendID,
		Role:    1, // 普通成员
	}
	return repository.AddGroupMember(member)
}

// KickMember 踢出群成员
func KickMember(operatorID, groupID, targetID uint) error {
	// 1. 检查操作者权限
	operatorRole, err := repository.GetMemberRole(groupID, operatorID)
	if err != nil {
		return err
	}
	// 只有群主(3)和管理员(2)可以踢人
	if operatorRole < 2 {
		return errors.New("权限不足")
	}

	// 2. 检查目标身份
	targetRole, err := repository.GetMemberRole(groupID, targetID)
	if err != nil {
		return errors.New("目标成员不存在")
	}
	// 不能踢比自己大或平级的人 (群主不能被踢，管理员不能踢管理员)
	if targetRole >= operatorRole {
		return errors.New("无法移除该成员")
	}

	return repository.RemoveGroupMember(groupID, targetID)
}

// MuteMember 禁言/解禁
func MuteMember(operatorID, groupID, targetID uint, muteState int) error {
	// 检查权限
	operatorRole, err := repository.GetMemberRole(groupID, operatorID)
	if err != nil {
		return err
	}
	// 仅管理员和群主可禁言
	if operatorRole < 2 {
		return errors.New("权限不足")
	}
	
	// 不能禁言比自己大或平级的人
	targetRole, err := repository.GetMemberRole(groupID, targetID)
	if err == nil && targetRole >= operatorRole {
		return errors.New("无法操作该成员")
	}

	return repository.UpdateMemberMuteStatus(groupID, targetID, muteState)
}

// SetGroupAdmin 设置/取消管理员
func SetGroupAdmin(ownerID, groupID, targetID uint, isAdmin bool) error {
	// 1. 只有群主可以设置管理员
	g, err := repository.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if g.OwnerID != ownerID {
		return errors.New("仅群主可设置管理员")
	}

	// 2. 目标角色
	role := 1
	if isAdmin {
		role = 2
	}

	return repository.UpdateMemberRole(groupID, targetID, role)
}

// TransferGroupOwner 转让群主
func TransferGroupOwner(ownerID, groupID, targetID uint) error {
	// 1. 只有群主可以转让
	g, err := repository.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if g.OwnerID != ownerID {
		return errors.New("仅群主可转让群组")
	}

	// 2. 目标必须是群成员
	if !repository.CheckGroupMember(groupID, targetID) {
		return errors.New("目标用户不在群内")
	}

	return repository.TransferGroupOwner(groupID, ownerID, targetID)
}
