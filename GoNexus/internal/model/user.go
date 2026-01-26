package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
// 包含：基础信息、安全信息、个人资料、状态控制
type User struct {

	// 包含 ID(uint主键), CreatedAt, UpdatedAt, DeletedAt(软删除)
	gorm.Model

	// 核心身份信息 (Identity)
	// UUID: 分布式唯一ID。
	// 为什么有了 ID 还要 UUID？因为自增 ID (1, 2, 3) 容易被爬虫遍历猜测用户量。
	// 对外暴露接口时用 UUID，对内数据库关联用 ID (性能更高)。
	UUID string `gorm:"type:char(36);uniqueIndex;not null;comment:分布式唯一ID"`

	Username string `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`

	Password string `gorm:"type:varchar(100);not null;comment:加密后的密码"`

	Phone string `gorm:"type:varchar(20);index;comment:手机号"`
	Email string `gorm:"type:varchar(100);index;comment:邮箱"`

	// 个人资料 (Profile) - RAG 也会用到这些做语义分析
	Nickname string `gorm:"type:varchar(50);comment:昵称"`
	Avatar   string `gorm:"type:varchar(255);default:'https://example.com/default.png';comment:头像URL"`

	// 个性签名：RAG 搜索时，可以根据签名匹配兴趣相投的好友
	Signature string `gorm:"type:varchar(255);comment:个性签名"`

	// 性别：0-未知 1-男 2-女 (使用 int 也就是 tinyint 存储，性能优于 string)
	Gender string `gorm:"type:varchar(20);default:0;comment:性别"`

	// Tags: 用户标签 (JSON格式存储)。
	// 例如：["Java", "Go", "二次元"]。
	// GORM 支持 serializer 自动序列化，Go代码里是用 []string，库里存的是 JSON 字符串。
	Tags []string `gorm:"serializer:json;type:json;comment:用户标签(JSON)"`

	// 状态与风控 (Status & Audit)

	// 状态：1-正常 2-冻结/拉黑
	Status int `gorm:"type:tinyint(1);default:1;comment:账户状态"`

	// 记录最后一次登录信息，用于安全审计
	LastLoginIp   string    `gorm:"type:varchar(50);comment:最后登录IP"`
	LastLoginTime time.Time `gorm:"comment:最后登录时间戳"`
	Birthday      time.Time `json:"birthday"`
	Location      string    `json:"location"`
	IsOnline      bool      `json:"is_online" gorm:"-"` // 在线状态(仅内存)
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate 是 GORM 的 Hook (钩子函数)
// 作用：在数据插入数据库 **之前**，自动执行这段代码。
// 复试亮点：利用 Hook 机制自动生成 UUID，业务代码就不用管了，解耦！
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return
}

// UserProfileResponse 定义一个干净的结构体，只包含前端需要的字段
type UserProfileResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastLoginTime time.Time `json:"last_login_time"`
	Username      string    `json:"username"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Email         string    `json:"email"`
	Tags          []string  `json:"tags"`      // 自动序列化的标签
	Signature     string    `json:"signature"` // 个性签名
	Gender        string    `json:"gender"`
	Birthday      time.Time `json:"birthday"`
	Location      string    `json:"location"`
	IsOnline      bool      `json:"is_online"`    // 在线状态
	UnreadCount   int       `json:"unread_count"` // 未读消息数
	LastMsg       string    `json:"last_msg"`     // 最后一条消息内容
}
