package model

import "gorm.io/gorm"

// Group 群组基本信息
type Group struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50);comment:群名"`
	OwnerID     uint   `json:"owner_id" gorm:"index;comment:群主ID"`
	Avatar      string `json:"avatar" gorm:"type:varchar(255);comment:群头像"`
	Notice      string `json:"notice" gorm:"type:varchar(255);comment:群公告"`
	Type        int    `json:"type" gorm:"type:int;commit:群类型"` //1.普通群，2公共大厅
	UnreadCount int    `json:"unread_count" gorm:"-"`           // 未读数 (仅内存)
}

// GroupMember 群成员关系表
type GroupMember struct {
	gorm.Model
	GroupID  uint   `json:"group_id" gorm:"index;comment:群ID"`
	UserID   uint   `json:"user_id" gorm:"index;comment:成员ID"`
	Nickname string `json:"nickname" gorm:"type:varchar(50);comment:群内昵称"`
	Role     int    `json:"role" gorm:"type:tinyint;default:1;comment:1普通/2管理员/3群主"`
	Muted    int    `json:"muted" gorm:"type:tinyint;default:0;comment:0正常/1禁言"`

	// 未读消息数
	UnreadCount int `json:"unread_count" gorm:"default:0;comment:未读消息数"`
}
