package dto

// FriendRequestResponse 好友申请响应结构 (包含申请者信息)
type FriendRequestResponse struct {
	ID              uint   `json:"id"`               // 申请记录ID
	RequesterID     uint   `json:"requester_id"`     // 申请者ID
	RequesterName   string `json:"requester_name"`   // 申请者昵称
	RequesterAvatar string `json:"requester_avatar"` // 申请者头像
	VerifyMsg       string `json:"verify_msg"`       // 验证消息
	CreatedAt       string `json:"created_at"`       // 申请时间
}

// RecommendResponse 推荐好友响应
type RecommendResponse struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Tags      []string `json:"tags"`
	Reason    string   `json:"reason"` // 推荐理由
	Score     int      `json:"score"`  // 推荐分值
	IsOnline  bool     `json:"is_online"`
}
