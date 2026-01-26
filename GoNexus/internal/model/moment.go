package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 动态/朋友圈内容
type Post struct {
	gorm.Model

	UserID uint `gorm:"index;not null;comment:发布者ID"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Content string `gorm:"type:text;comment:动态内容"`
	
	// Media 存储图片/视频URL列表，JSON格式
	Media string `gorm:"type:json;comment:媒体资源URL列表(JSON)"`

	// Visibility: 0-公开(广场可见), 1-好友可见, 2-仅自己可见
	Visibility int `gorm:"type:tinyint;default:0;comment:可见性:0-公开,1-好友,2-私密"`

	LikeCount    int `gorm:"default:0;comment:点赞数"`
	CommentCount int `gorm:"default:0;comment:评论数"`

	// 扩展字段：位置、心情等
	Location string `gorm:"type:varchar(255);comment:位置"`
	Mood     string `gorm:"type:varchar(50);comment:心情(AI分析或手动选择)"`
}

func (Post) TableName() string {
	return "posts"
}

// Comment 评论
type Comment struct {
	gorm.Model

	PostID uint `gorm:"index;not null;comment:关联动态ID"`
	Post   Post `gorm:"foreignKey:PostID"`

	UserID uint `gorm:"index;not null;comment:评论者ID"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Content string `gorm:"type:varchar(500);not null;comment:评论内容"`

	// ParentID 支持回复评论
	ParentID *uint `gorm:"index;comment:父评论ID"`
}

func (Comment) TableName() string {
	return "comments"
}

// Like 点赞
type Like struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`

	PostID uint `gorm:"index;not null;uniqueIndex:idx_post_user;comment:动态ID"`
	UserID uint `gorm:"index;not null;uniqueIndex:idx_post_user;comment:点赞者ID"`
}

func (Like) TableName() string {
	return "likes"
}
