package dto

import "time"

type CreatePostRequest struct {
	Content    string   `json:"content" binding:"required"`
	Media      []string `json:"media"`
	Visibility int      `json:"visibility"` // 0-公开, 1-好友, 2-私密
	Location   string   `json:"location"`
	Mood       string   `json:"mood"`
}

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	ParentID *uint `json:"parent_id"`
}

type LikePostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}

type GetMomentsRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Type     string `form:"type"` // "all" (广场) or "friend" (朋友圈) or "user" (个人)
	UserID   uint   `form:"user_id"` // 查看特定用户的动态
}

type PostResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	UserNickname string    `json:"user_nickname"`
	UserAvatar   string    `json:"user_avatar"`
	Content      string    `json:"content"`
	Media        []string  `json:"media"`
	Visibility   int       `json:"visibility"`
	Location     string    `json:"location"`
	Mood         string    `json:"mood"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	IsLiked      bool      `json:"is_liked"` // 当前用户是否已点赞
	CreatedAt    time.Time `json:"created_at"`
	Comments     []CommentResponse `json:"comments,omitempty"` // 仅在详情页或预加载时返回
}

type CommentResponse struct {
	ID           uint      `json:"id"`
	PostID       uint      `json:"post_id"`
	UserID       uint      `json:"user_id"`
	UserNickname string    `json:"user_nickname"`
	UserAvatar   string    `json:"user_avatar"`
	Content      string    `json:"content"`
	ParentID     *uint     `json:"parent_id"`
	CreatedAt    time.Time `json:"created_at"`
}
