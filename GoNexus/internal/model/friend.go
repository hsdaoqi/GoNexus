package model

import "gorm.io/gorm"

// --- 常量定义 ---
const (
	// 好友申请状态
	RequestStatusPending  = 0 // 等待处理
	RequestStatusAccepted = 1 // 已同意
	RequestStatusRefused  = 2 // 已拒绝

	// 来源渠道 (复试亮点：数据埋点)
	SourceSearch = 1 // 搜索发现
	SourceGroup  = 2 // 群聊
	SourceQrCode = 3 // 二维码
)

// FriendRequest 好友申请表 (记录“加好友”这个事件)
type FriendRequest struct {
	gorm.Model
	// 谁(Requester) 向 谁(Receiver) 发起了申请
	RequesterID uint `json:"requester_id" gorm:"not null;index;comment:申请人ID"`
	ReceiverID  uint `json:"receiver_id" gorm:"not null;index;comment:接收人ID"`

	// 验证信息
	VerifyMsg string `json:"verify_msg" gorm:"type:varchar(100);comment:验证消息"`
	// 处理状态
	Status int `json:"status" gorm:"type:tinyint;default:0;comment:状态(0待处理/1已同意/2已拒绝)"`
}

// Friend 好友关系表 (只存已确立的关系)
type Friend struct {
	gorm.Model
	// 复合唯一索引：确保 A 和 B 之间只能有一条关系记录，防止重复
	UserID   uint `json:"user_id" gorm:"not null;index;uniqueIndex:idx_user_friend;comment:我的ID"`
	FriendID uint `json:"friend_id" gorm:"not null;index;uniqueIndex:idx_user_friend;comment:对方ID"`

	// 备注名 (只是我看对方的名字)
	Remark string `json:"remark" gorm:"type:varchar(50);comment:备注名"`

	// 关系设置
	IsDisturb bool `json:"is_disturb" gorm:"default:false;comment:免打扰"`
	IsTop     bool `json:"is_top" gorm:"default:false;comment:置顶"`

	// 来源 (记录你们是怎么认识的)
	Source int `json:"source" gorm:"type:tinyint;default:1;comment:来源"`

	// 未读消息数
	UnreadCount int `json:"unread_count" gorm:"default:0;comment:未读消息数"`
}
