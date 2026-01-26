package service

import (
	"encoding/json"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"strings"
)

type MomentService struct{}

var MomentServiceApp = new(MomentService)

// CreatePost 发布动态
func (s *MomentService) CreatePost(userID uint, req dto.CreatePostRequest) error {
	// 简单的关键词分析模拟 (TODO: 后续对接真实 AI 服务)
	mood := "Neutral"
	if req.Mood != "" {
		mood = req.Mood
	} else {
		mood = analyzeMood(req.Content)
	}

	mediaJSON, _ := json.Marshal(req.Media)

	post := &model.Post{
		UserID:     userID,
		Content:    req.Content,
		Media:      string(mediaJSON),
		Visibility: req.Visibility,
		Location:   req.Location,
		Mood:       mood,
	}

	return repository.MomentRepo.CreatePost(post)
}

func analyzeMood(content string) string {
	content = strings.ToLower(content)
	if strings.Contains(content, "开心") || strings.Contains(content, "快乐") || strings.Contains(content, "哈哈") || strings.Contains(content, "棒") || strings.Contains(content, "happy") {
		return "Happy"
	}
	if strings.Contains(content, "难过") || strings.Contains(content, "伤心") || strings.Contains(content, "哭") || strings.Contains(content, "累") || strings.Contains(content, "sad") {
		return "Sad"
	}
	if strings.Contains(content, "生气") || strings.Contains(content, "愤怒") || strings.Contains(content, "烦") || strings.Contains(content, "angry") {
		return "Angry"
	}
	if strings.Contains(content, "惊讶") || strings.Contains(content, "哇") || strings.Contains(content, "wow") {
		return "Surprised"
	}
	return "Neutral"
}

// GetMoments 获取动态列表
func (s *MomentService) GetMoments(userID uint, req dto.GetMomentsRequest) ([]dto.PostResponse, int64, error) {
	// 确定查询的用户ID (如果是查看个人主页，则为 req.UserID; 否则为当前用户 userID)
	targetUserID := userID
	if req.Type == "user" && req.UserID != 0 {
		targetUserID = req.UserID
	}

	posts, total, err := repository.MomentRepo.GetPosts(req.Type, targetUserID, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		// 解析媒体 JSON
		var media []string
		_ = json.Unmarshal([]byte(post.Media), &media)

		// 检查当前用户是否点赞
		isLiked, _ := repository.MomentRepo.IsLiked(post.ID, userID)

		// 获取评论 (这里只获取前3条作为预览，或者不获取，取决于需求。暂时不获取，点击详情再拉取)
		// comments, _ := repository.MomentRepo.GetComments(post.ID)

		res := dto.PostResponse{
			ID:           post.ID,
			UserID:       post.UserID,
			UserNickname: post.User.Nickname,
			UserAvatar:   post.User.Avatar,
			Content:      post.Content,
			Media:        media,
			Visibility:   post.Visibility,
			Location:     post.Location,
			Mood:         post.Mood,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			IsLiked:      isLiked,
			CreatedAt:    post.CreatedAt,
		}

		// 如果用户昵称为空，使用用户名
		if res.UserNickname == "" {
			res.UserNickname = post.User.Username
		}

		responses = append(responses, res)
	}

	return responses, total, nil
}

// CreateComment 发表评论
func (s *MomentService) CreateComment(userID uint, req dto.CreateCommentRequest) (dto.CommentResponse, error) {
	comment := &model.Comment{
		PostID:   req.PostID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	err := repository.MomentRepo.CreateComment(comment)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	// 为了返回完整的用户信息，需要重新查一下或者手动填充
	// 这里简单起见，假设前端已经有当前用户的信息，或者我们再查一次 User
	user, _ := repository.GetUserByID(userID)

	return dto.CommentResponse{
		ID:           comment.ID,
		PostID:       comment.PostID,
		UserID:       comment.UserID,
		UserNickname: user.Nickname,
		UserAvatar:   user.Avatar,
		Content:      comment.Content,
		ParentID:     comment.ParentID,
		CreatedAt:    comment.CreatedAt,
	}, nil
}

// GetComments 获取评论
func (s *MomentService) GetComments(postID uint) ([]dto.CommentResponse, error) {
	comments, err := repository.MomentRepo.GetComments(postID)
	if err != nil {
		return nil, err
	}

	var responses []dto.CommentResponse
	for _, c := range comments {
		res := dto.CommentResponse{
			ID:           c.ID,
			PostID:       c.PostID,
			UserID:       c.UserID,
			UserNickname: c.User.Nickname,
			UserAvatar:   c.User.Avatar,
			Content:      c.Content,
			ParentID:     c.ParentID,
			CreatedAt:    c.CreatedAt,
		}
		if res.UserNickname == "" {
			res.UserNickname = c.User.Username
		}
		responses = append(responses, res)
	}
	return responses, nil
}

// ToggleLike 点赞
func (s *MomentService) ToggleLike(userID uint, req dto.LikePostRequest) (bool, int, error) {
	return repository.MomentRepo.ToggleLike(req.PostID, userID)
}
