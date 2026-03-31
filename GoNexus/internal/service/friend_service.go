// internal/service/friend_service.go
package service

import (
	"errors"
	"fmt"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"sort"
	"strings"
	"time"
)

var (
	ErrCannotAddSelf      = errors.New("不能添加自己")
	ErrAlreadyFriend      = errors.New("你们已经是好友了")
	ErrRequestPending     = errors.New("已发送申请，请耐心等待")
	ErrRequestNotFound    = errors.New("申请记录不存在")
	ErrNoPermission       = errors.New("无权处理此申请")
	ErrRequestAlreadyDone = errors.New("该申请已被处理")
)

type OnlineChecker interface {
	IsUserOnline(userID uint) bool
}
type FriendService struct {
	onlineChecker OnlineChecker
}

// NewFriendService 构造时注入
func NewFriendService(checker OnlineChecker) *FriendService {
	return &FriendService{onlineChecker: checker}
}

var FriendSvc = NewFriendService(nil)

func (s *FriendService) SendFriendRequest(requesterID, receiverID uint, verifyMsg string) error {
	if requesterID == receiverID {
		return ErrCannotAddSelf
	}
	if _, err := repository.UserRepo.GetByID(receiverID); err != nil {
		return ErrUserNotFound
	}
	if repository.FriendRepo.IsFriend(requesterID, receiverID) {
		return ErrAlreadyFriend
	}
	if repository.FriendRepo.CheckRequestExist(requesterID, receiverID) {
		return ErrRequestPending
	}
	return repository.FriendRepo.CreateRequest(&model.FriendRequest{
		RequesterID: requesterID,
		ReceiverID:  receiverID,
		VerifyMsg:   verifyMsg,
		Status:      model.RequestStatusPending,
	})
}

// HandleFriendRequest 返回 requesterID 供 Controller 层发 WebSocket 通知
func (s *FriendService) HandleFriendRequest(userID, requestID uint, action int) (uint, error) {
	req, err := repository.FriendRepo.GetPendingRequest(requestID)
	if err != nil {
		return 0, ErrRequestNotFound
	}
	if req.ReceiverID != userID {
		return 0, ErrNoPermission
	}
	if req.Status != model.RequestStatusPending {
		return 0, ErrRequestAlreadyDone
	}

	if action == model.RequestStatusRefused {
		// 拒绝：直接删除记录
		return req.RequesterID, repository.FriendRepo.DeleteRecord(req.RequesterID, req.ReceiverID)
	}

	// 同意：事务由 Repository 负责
	return req.RequesterID, repository.FriendRepo.AcceptRequest(req)
}

func (s *FriendService) GetFriendList(userID uint) ([]model.UserProfileResponse, error) {
	friends, err := repository.FriendRepo.GetList(userID)
	if err != nil {
		return nil, err
	}
	for i := range friends {
		if s.onlineChecker != nil {
			friends[i].IsOnline = s.onlineChecker.IsUserOnline(friends[i].ID)
		}
	}
	return friends, nil
}
func (s *FriendService) DeleteFriend(userID, friendID uint) error {
	if err := repository.FriendRepo.DeleteRecord(userID, friendID); err != nil {
		return err
	}
	return repository.FriendRepo.DeleteRecord(friendID, userID)
}

func (s *FriendService) GetPendingRequests(userID uint) ([]dto.FriendRequestResponse, error) {
	return repository.FriendRepo.GetPendingRequests(userID)
}

func (s *FriendService) RecommendFriends(userID uint) ([]dto.RecommendResponse, error) {
	myFriendIDs, err := repository.FriendRepo.GetFriendIDs(userID)
	if err != nil {
		return nil, err
	}
	me, err := repository.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	candidateMap, err := repository.FriendRepo.GetSecondDegreeFriends(myFriendIDs, userID)
	if err != nil {
		return nil, err
	}

	friendSet := make(map[uint]bool)
	for _, id := range myFriendIDs {
		friendSet[id] = true
	}
	friendSet[userID] = true

	var candidateIDs []uint
	for id := range candidateMap {
		if !friendSet[id] {
			candidateIDs = append(candidateIDs, id)
		}
	}

	users, err := repository.FriendRepo.GetUsersByIDs(candidateIDs)
	if err != nil {
		return nil, err
	}

	var recommendations []dto.RecommendResponse
	for _, u := range users {
		commonFriendsCount := candidateMap[u.ID]
		commonTags := intersectTags(me.Tags, u.Tags)
		score := commonFriendsCount*5 + len(commonTags)*3

		var reasons []string
		if commonFriendsCount > 0 {
			reasons = append(reasons, fmt.Sprintf("%d位共同好友", commonFriendsCount))
		}
		if len(commonTags) > 0 {
			reasons = append(reasons, fmt.Sprintf("共同兴趣: %s", strings.Join(commonTags, ",")))
		}
		reason := strings.Join(reasons, " | ")
		if reason == "" {
			reason = "缘分推荐"
		}
		recommendations = append(recommendations, dto.RecommendResponse{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Tags:     u.Tags,
			Reason:   reason,
			Score:    score,
		})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}
	return recommendations, nil
}

func intersectTags(a, b []string) []string {
	m := make(map[string]bool)
	var res []string
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if m[item] {
			res = append(res, item)
		}
	}
	return res
}

// internal/service/friend_service.go
func (s *FriendService) calcRecommendScore(me *model.User, candidate model.User, commonFriends int) (int, string) {
	score := 0
	var reasons []string

	// 1. 共同好友（社交关系最强信号）
	if commonFriends > 0 {
		score += commonFriends * 5
		reasons = append(reasons, fmt.Sprintf("%d位共同好友", commonFriends))
	}

	// 2. 共同标签（兴趣匹配）
	commonTags := intersectTags(me.Tags, candidate.Tags)
	if len(commonTags) > 0 {
		score += len(commonTags) * 3
		reasons = append(reasons, fmt.Sprintf("共同兴趣: %s", strings.Join(commonTags, ",")))
	}

	// 3. 同城加分
	if me.Location != "" && me.Location == candidate.Location {
		score += 4
		reasons = append(reasons, "同城")
	}

	// 4. 在线用户优先（提升活跃用户曝光）
	if s.onlineChecker != nil && s.onlineChecker.IsUserOnline(candidate.ID) {
		score += 2
		reasons = append(reasons, "当前在线")
	}

	// 5. 新用户扶持（注册7天内的新用户加分，冷启动问题）
	if time.Since(candidate.CreatedAt) < 7*24*time.Hour {
		score += 3
		reasons = append(reasons, "新用户")
	}

	reason := strings.Join(reasons, " | ")
	if reason == "" {
		reason = "缘分推荐"
	}
	return score, reason
}
