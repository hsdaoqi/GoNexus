package repository

import (
	"go-nexus/internal/model"
	"go-nexus/pkg/global"

	"gorm.io/gorm"
)

type MomentRepository struct{}

var MomentRepo = new(MomentRepository)

// CreatePost 发布动态
func (r *MomentRepository) CreatePost(post *model.Post) error {
	return global.DB.Create(post).Error
}

// GetPostByID 获取单条动态
func (r *MomentRepository) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	err := global.DB.Preload("User").First(&post, id).Error
	return &post, err
}

// GetPosts 获取动态列表
// scope: "all" (广场), "friend" (朋友圈), "user" (个人)
func (r *MomentRepository) GetPosts(scope string, userID uint, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	db := global.DB.Model(&model.Post{}).Preload("User")

	if scope == "user" {
		// 查看特定用户的动态 (userID 此时是目标用户ID)
		db = db.Where("user_id = ?", userID)
		// 如果查看的是别人的主页，可能需要过滤私密动态 (Visibility=2)
		// 这里暂不处理复杂的权限，假设查看自己或公开的
	} else if scope == "friend" {
		// 朋友圈：自己 + 好友 的动态
		// 1. 获取好友ID列表 (子查询)
		subQuery := global.DB.Table("friends").Select("friend_id").Where("user_id = ?", userID)

		// 2. 查询: (作者是自己 OR 作者在好友列表中) AND (可见性 != 私密)
		// 注意：自己的私密动态应该可见，好友的私密动态不可见
		// 简化逻辑：查询 (UserID = me) OR (UserID IN friends AND Visibility IN (0, 1))

		// 这是一个组合条件
		db = db.Where(
			global.DB.Where("user_id = ?", userID).
				Or(
					global.DB.Where("user_id IN (?)", subQuery).Where("visibility IN ?", []int{0, 1}),
				),
		)

	} else { // scope == "all" (广场)
		// 广场只看公开的 (Visibility=0)
		db = db.Where("visibility = ?", 0)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// DeletePost 删除动态
func (r *MomentRepository) DeletePost(id uint) error {
	return global.DB.Delete(&model.Post{}, id).Error
}

// CreateComment 发表评论
func (r *MomentRepository) CreateComment(comment *model.Comment) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		// 更新评论计数
		return tx.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	})
}

// GetComments 获取评论列表
func (r *MomentRepository) GetComments(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := global.DB.Where("post_id = ?", postID).Preload("User").Order("created_at asc").Find(&comments).Error
	return comments, err
}

// ToggleLike 点赞/取消点赞
func (r *MomentRepository) ToggleLike(postID, userID uint) (bool, int, error) {
	var like model.Like
	var isLiked bool
	var likeCount int

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否已点赞
		result := tx.Where("post_id = ? AND user_id = ?", postID, userID).First(&like)

		if result.Error == nil {
			// 已点赞 -> 取消点赞
			if err := tx.Delete(&like).Error; err != nil {
				return err
			}
			// 减少计数
			if err := tx.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
				return err
			}
			isLiked = false
		} else if result.Error == gorm.ErrRecordNotFound {
			// 未点赞 -> 点赞
			newLike := model.Like{PostID: postID, UserID: userID}
			if err := tx.Create(&newLike).Error; err != nil {
				return err
			}
			// 增加计数
			if err := tx.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
				return err
			}
			isLiked = true
		} else {
			return result.Error
		}

		// 获取最新点赞数
		var post model.Post
		if err := tx.Select("like_count").First(&post, postID).Error; err != nil {
			return err
		}
		likeCount = post.LikeCount

		return nil
	})

	return isLiked, likeCount, err
}

// IsLiked 检查用户是否点赞
func (r *MomentRepository) IsLiked(postID, userID uint) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Like{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	return count > 0, err
}
