package service

import (
	"errors"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
)

type GroupService struct{}

var GroupSrv = &GroupService{}

// CreateGroup 业务逻辑
func (g *GroupService) Create(ownerID uint, name, avatar, notice string) (*model.Group, error) {
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

	if err := repository.GroupRepo.Create(group); err != nil {
		return nil, err
	}
	return group, nil
}

// GetJoinedGroups 业务逻辑
func (g *GroupService) GetJoinedGroups(userID uint) ([]model.Group, error) {
	return repository.GroupRepo.GetJoinedGroups(userID)
}

// UpdateGroup 仅群主可更新群名称/头像/公告
func (g *GroupService) Update(ownerID, groupID uint, name, avatar, notice string) error {
	group, err := repository.GroupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if group.OwnerID != ownerID {
		return errors.New("仅群主可编辑群资料")
	}
	// 按需更新
	if name != "" {
		group.Name = name
	}
	if avatar != "" {
		group.Avatar = avatar
	}
	if notice != "" {
		group.Notice = notice
	}
	return repository.GroupRepo.Save(group)
}

// GetGroupMembers 获取群成员
func (g *GroupService) GetMembers(groupID uint) ([]dto.GroupMemberResponse, error) {
	return repository.GroupRepo.GetMembers(groupID)
}

// InviteFriendToGroup 邀请好友入群
func (g *GroupService) InviteFriend(groupID, userID, friendID uint) error {
	// 1. 验证是否是群成员 (只有群成员能邀请)
	if !repository.GroupRepo.CheckMember(groupID, userID) {
		return errors.New("你不是该群成员，无法邀请")
	}

	// 2. 验证目标是否已经在群里
	if repository.GroupRepo.CheckMember(groupID, friendID) {
		return errors.New("该用户已经在群里了")
	}

	// 3. 入库
	member := &model.GroupMember{
		GroupID: groupID,
		UserID:  friendID,
		Role:    1, // 普通成员
	}
	return repository.GroupRepo.AddMember(member)
}

// KickMember 踢出群成员
func (g *GroupService) KickMember(operatorID, groupID, targetID uint) error {
	// 1. 检查操作者权限
	operatorRole, err := repository.GroupRepo.GetMemberRole(groupID, operatorID)
	if err != nil {
		return err
	}
	// 只有群主(3)和管理员(2)可以踢人
	if operatorRole < 2 {
		return errors.New("权限不足")
	}

	// 2. 检查目标身份
	targetRole, err := repository.GroupRepo.GetMemberRole(groupID, targetID)
	if err != nil {
		return errors.New("目标成员不存在")
	}
	// 不能踢比自己大或平级的人 (群主不能被踢，管理员不能踢管理员)
	if targetRole >= operatorRole {
		return errors.New("无法移除该成员")
	}

	return repository.GroupRepo.RemoveMember(groupID, targetID)
}

// MuteMember 禁言/解禁
func (g *GroupService) MuteMember(operatorID, groupID, targetID uint, muteState int) error {
	// 检查权限
	operatorRole, err := repository.GroupRepo.GetMemberRole(groupID, operatorID)
	if err != nil {
		return err
	}
	// 仅管理员和群主可禁言
	if operatorRole < 2 {
		return errors.New("权限不足")
	}

	// 不能禁言比自己大或平级的人
	targetRole, err := repository.GroupRepo.GetMemberRole(groupID, targetID)
	if err == nil && targetRole >= operatorRole {
		return errors.New("无法操作该成员")
	}

	return repository.GroupRepo.UpdateMemberMuteStatus(groupID, targetID, muteState)
}

// SetGroupAdmin 设置/取消管理员
func (g *GroupService) SetAdmin(operatorID, groupID, memberID uint, isAdmin bool) error {
	// 1. 检查操作者是否是群主
	group, err := repository.GroupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if group.OwnerID != operatorID {
		return errors.New("仅群主可设置管理员")
	}

	// 2. 检查目标是否是群成员
	if !repository.GroupRepo.CheckMember(groupID, memberID) {
		return errors.New("目标不是群成员")
	}

	// 3. 设置角色
	// 3 = 群主, 2 = 管理员, 1 = 普通成员
	role := 1
	if isAdmin {
		role = 2
	}
	return repository.GroupRepo.UpdateMemberRole(groupID, memberID, role)
}

// TransferGroupOwner 转让群主
func (g *GroupService) TransferOwner(operatorID, groupID, newOwnerID uint) error {
	// 1. 检查操作者是否是群主
	group, err := repository.GroupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if group.OwnerID != operatorID {
		return errors.New("仅群主可转让群")
	}

	// 2. 检查新群主是否是群成员
	if !repository.GroupRepo.CheckMember(groupID, newOwnerID) {
		return errors.New("目标不是群成员")
	}

	if operatorID == newOwnerID {
		return errors.New("不能转让自己")
	}

	// 3. 执行转让事务
	return repository.GroupRepo.TransferGroupOwner(groupID, operatorID, newOwnerID)
}
