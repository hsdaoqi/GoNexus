package model

import "gorm.io/gorm"

// Group 群组基本信息
type Group struct {
	gorm.Model
	Name    string `gorm:"type:varchar(50);comment:群名"`
	OwnerID uint   `gorm:"index;comment:群主ID"`
	Avatar  string `gorm:"type:varchar(255);comment:群头像"`
	Notice  string `gorm:"type:varchar(255);comment:群公告"`
}

// GroupMember 群成员关系表
type GroupMember struct {
	gorm.Model
	GroupID  uint   `gorm:"index;comment:群ID"`
	UserID   uint   `gorm:"index;comment:成员ID"`
	Nickname string `gorm:"type:varchar(50);comment:群内昵称"`
	Role     int    `gorm:"type:tinyint;default:1;comment:1普通/2管理员/3群主"`
	Muted    int    `gorm:"type:tinyint;default:0;comment:0正常/1禁言"`
}
