package dto

import "time"

// ChatHistoryResponse 聊天历史响应
type ChatHistoryResponse struct {
	ID         uint      `json:"id"`
	FromUserID uint      `json:"from_user_id"`
	ToUserID   uint      `json:"to_user_id"`
	ChatType   int       `json:"chat_type"`
	MsgType    int       `json:"msg_type"`
	Content    string    `json:"content"`
	Url        string    `json:"url"`
	FileSize   int64     `json:"file_size"`
	CreatedAt  time.Time `json:"created_at"`

	// 发送者信息
	SenderAvatar   string `json:"sender_avatar"`
	SenderNickname string `json:"sender_nickname"`
}
