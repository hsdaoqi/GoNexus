package model

import "gorm.io/gorm"

const (
	// 消息类型
	MsgTypeText  = 1 // 文本
	MsgTypeImage = 2 // 图片
	MsgTypeAudio = 3 // 语音
	MsgTypeFile  = 4 // 文件

	// 聊天类型 (复试考点: 区分单聊群聊)
	ChatTypePrivate = 1 // 单聊
	ChatTypeGroup   = 2 // 群聊
)

type Message struct {
	gorm.Model
	// 核心关系
	FromUserID uint `json:"from_user_id" gorm:"not null;index;comment:发送者ID"`
	ToUserID   uint `json:"to_user_id" gorm:"not null;index;comment:接收者ID(如果是群聊，这里存GroupID)"`

	// 聊天类型：1-单聊，2-群聊
	// 这个字段决定了 ToUserID 到底是 UserID 还是 GroupID
	ChatType int `json:"chat_type" gorm:"type:tinyint;not null;comment:聊天类型"`

	// 消息内容
	MsgType int    `json:"msg_type" gorm:"type:tinyint;not null;comment:消息类型(文本/图片/语音)"`
	Content string `json:"content" gorm:"type:text;comment:文本内容"` // 文本消息存这里

	// 多媒体扩展 (如果是图片/文件，存 URL)
	Url        string `json:"url" gorm:"type:varchar(255);comment:文件/图片地址"`
	PicPreview string `json:"pic_preview" gorm:"type:varchar(255);comment:图片缩略图(优化加载)"`
	FileSize   int64  `json:"file_size" gorm:"comment:文件大小(字节)"`
	FileName   string `json:"file_name" gorm:"type:varchar(255);comment:文件名"`

	// 撤回功能预留
	IsRevoked bool `json:"is_revoked" gorm:"default:false;comment:是否已撤回"`
}
