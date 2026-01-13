package dto

import "encoding/json"

// 消息类型常量 (移到这里，方便全局引用)
const (
	TypeHeartbeat = 0 // 心跳
	TypeText      = 1 // 文本
	TypeImage     = 2 // 图片
	TypeAudio     = 3 // 语音
	TypeFriendReq = 4 // 好友申请通知
	TypeFriendAns = 5 //对方同意了我的申请
)

// 聊天类型常量
const (
	ChatTypePrivate = 1 // 单聊
	ChatTypeGroup   = 2 // 群聊
	ChatTypeSystem  = 3 // 系统广播 (新增)
)

// ProtocolMsg WebSocket 传输协议对象
type ProtocolMsg struct {
	// 1. 基础信息
	Code int    `json:"code"` // 状态码 (用于服务端返回发送成功/失败)
	Msg  string `json:"msg"`  // 错误提示

	// 2. 路由信息
	Type       int  `json:"type"` // 消息类型
	FromUserID uint `json:"from_user_id"`
	ToUserID   uint `json:"to_user_id"` // 接收人或群ID
	ChatType   int  `json:"chat_type"`  // 1单聊 2群聊

	// 3. 消息载体
	Content  string `json:"content"` // 文本
	Url      string `json:"url"`     // 图片地址
	FileSize int64  `json:"file_size"`

	// 4. 用户信息 (前端展示用，数据库不存，但转发时要带上)
	SenderAvatar   string `json:"sender_avatar"`
	SenderNickname string `json:"sender_nickname"`

	// 5. 好友申请相关字段 (当type=4时使用)
	RequestID     uint   `json:"request_id,omitempty"`     // 申请记录ID
	RequesterName string `json:"requester_name,omitempty"` // 申请者姓名
	VerifyMsg     string `json:"verify_msg,omitempty"`     // 验证消息

	// 6. 元数据
	SendTime string `json:"send_time"` // 2023-12-12 10:00:00
}

// ToBytes 序列化
func (m *ProtocolMsg) ToBytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
