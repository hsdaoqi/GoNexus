package dto

import "go-nexus/internal/model"

// GroupMemberResponse 群成员详细信息
type GroupMemberResponse struct {
	model.GroupMember
	UserAvatar string `json:"user_avatar"`
	UserName   string `json:"user_name"`
}
