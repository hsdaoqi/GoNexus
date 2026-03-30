package service

import (
	"errors"
	"fmt"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"go-nexus/pkg/global"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// OnlineStatusChecker 在线状态检查接口
type OnlineStatusChecker interface {
	IsUserOnline(userID uint) bool
}

// SendFriendRequest 发起好友申请
func SendFriendRequest(requesterID, receiverID uint, verifyMsg string) error {
	// 1. 不能加自己
	if requesterID == receiverID {
		return errors.New("不能添加自己")
	}

	// 2. 检查对方存不存在 (严谨)
	if _, err := repository.GetUserByID(receiverID); err != nil {
		return errors.New("目标用户不存在")
	}

	// 3. 检查是不是已经是好友了
	if repository.IsFriend(requesterID, receiverID) {
		return errors.New("你们已经是好友了")
	}

	// 4. 检查有没有未处理的申请 (防骚扰)
	if repository.CheckRequestExist(requesterID, receiverID) {
		return errors.New("已发送申请，请耐心等待")
	}

	// 5. 入库
	req := &model.FriendRequest{
		RequesterID: requesterID,
		ReceiverID:  receiverID,
		VerifyMsg:   verifyMsg,
		Status:      model.RequestStatusPending,
	}
	return repository.CreateFriendRequest(req)
}

// HandleFriendRequest 处理好友申请 (同意/拒绝)
// requestID: 申请记录的 ID (不是人的 ID)
// action: 1-同意, 2-拒绝
func HandleFriendRequest(userID, requestID uint, action int) (uint, error) {
	// 1. 查出这条申请
	req, err := repository.GetPendingRequest(requestID)
	if err != nil {
		return 0, errors.New("申请记录不存在")
	}

	// 2. 权限校验：确保这条申请是发给我的
	if req.ReceiverID != userID {
		return 0, errors.New("无权处理此申请")
	}

	// 3. 校验状态：必须是待处理
	if req.Status != model.RequestStatusPending {
		return 0, errors.New("该申请已被处理")
	}

	//如果是拒绝，直接删除记录，避免数据库堆积
	if action == model.RequestStatusRefused {
		req.Status = model.RequestStatusRefused
		return req.RequesterID, global.DB.Delete(req).Error
	}

	// --- 重点：如果是同意，开启事务 ---
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// A. 更新申请状态为“已同意”
		req.Status = model.RequestStatusAccepted
		if err := tx.Save(req).Error; err != nil {
			return err
		}
		// 先检查是否存在已删除的记录，如果存在则恢复，否则创建新记录

		// 检查并恢复/创建 我 -> 他 的关系
		var friendA model.Friend
		// 使用Unscoped()包含软删除记录
		result := tx.Unscoped().Where("user_id = ? AND friend_id = ?", req.ReceiverID, req.RequesterID).First(&friendA)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// 不存在记录，创建新记录
				friendA = model.Friend{
					UserID:   req.ReceiverID,  // 我
					FriendID: req.RequesterID, // 他
					Source:   model.SourceSearch,
				}
				if err := tx.Create(&friendA).Error; err != nil {
					return err
				}
			} else {
				// 其他错误
				return result.Error
			}
		} else {
			// 记录存在，恢复软删除
			if err := tx.Unscoped().Model(&friendA).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}

		// 检查并恢复/创建 他 -> 我的关系
		var friendB model.Friend
		// 使用Unscoped()包含软删除记录
		result = tx.Unscoped().Where("user_id = ? AND friend_id = ?", req.RequesterID, req.ReceiverID).First(&friendB)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// 不存在记录，创建新记录
				friendB = model.Friend{
					UserID:   req.RequesterID, // 他
					FriendID: req.ReceiverID,  // 我
					Source:   model.SourceSearch,
				}
				if err := tx.Create(&friendB).Error; err != nil {
					return err
				}
			} else {
				// 其他错误
				return result.Error
			}
		} else {
			// 记录存在，恢复软删除
			if err := tx.Unscoped().Model(&friendB).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}

		return nil // 提交事务
	})

	return req.RequesterID, err
}

// GetFriendList 获取好友列表（原始版本，不包含在线状态）
func GetFriendList(userID uint) ([]model.UserProfileResponse, error) {
	friends, err := repository.GetFriendList(userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// GetFriendListWithOnlineStatus 获取好友列表并包含在线状态
func GetFriendListWithOnlineStatus(userID uint, checker OnlineStatusChecker) ([]model.UserProfileResponse, error) {
	friends, err := repository.GetFriendList(userID)
	if err != nil {
		return nil, err
	}

	// 为每个好友添加在线状态
	if checker != nil {
		for i := range friends {
			friends[i].IsOnline = checker.IsUserOnline(friends[i].ID)
		}
	}

	return friends, nil
}

// DeleteFriendRecord 删除好友
func DeleteFriendRecord(userID, FriendID uint) error {
	//两边都删
	err := repository.DeleteFriendRecord(userID, FriendID)
	if err != nil {
		return err
	}
	err = repository.DeleteFriendRecord(FriendID, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetPendingRequests 获取待处理的好友请求
func GetPendingRequests(userID uint) ([]dto.FriendRequestResponse, error) {
	return repository.GetPendingRequests(userID)
}

// RecommendFriends 推荐好友算法
func RecommendFriends(userID uint) ([]dto.RecommendResponse, error) {
	// 1. 获取我的好友ID列表
	myFriendIDs, err := repository.GetFriendIDs(userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取我的信息 (为了拿 Tags)
	me, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 3. 获取二度好友 (Candidates) 及其出现频次
	// Map[FriendID]Count
	candidateMap, err := repository.GetSecondDegreeFriends(myFriendIDs, userID)
	if err != nil {
		return nil, err
	}

	// 4. 过滤掉已经是好友的人 (虽然 SQL 里过滤了，但为了双重保险和扩展性，这里再处理一下)
	// 同时准备 IDs 用于查询详情
	finalCandidateIDs := make([]uint, 0)
	friendSet := make(map[uint]bool)
	for _, id := range myFriendIDs {
		friendSet[id] = true
	}
	// 还要过滤掉自己
	friendSet[userID] = true

	for id := range candidateMap {
		if !friendSet[id] {
			finalCandidateIDs = append(finalCandidateIDs, id)
		}
	}

	// 5. 获取候选人详情 (Tags)
	users, err := repository.GetUsersByIDs(finalCandidateIDs)
	if err != nil {
		return nil, err
	}

	// 6. 计算得分
	var recommendations []dto.RecommendResponse

	for _, u := range users {
		commonFriendsCount := candidateMap[u.ID]

		// 计算共同 Tags
		commonTags := intersectTags(me.Tags, u.Tags)

		// 算法公式
		// Score = 共同好友 * 5 + 共同标签 * 3
		// 可以根据实际情况调整权重
		score := commonFriendsCount*5 + len(commonTags)*3

		// 构建推荐理由
		var reasons []string
		if commonFriendsCount > 0 {
			reasons = append(reasons, fmt.Sprintf("%d位共同好友", commonFriendsCount))
		}
		if len(commonTags) > 0 {
			reasons = append(reasons, fmt.Sprintf("共同兴趣: %s", strings.Join(commonTags, ",")))
		}
		reasonStr := strings.Join(reasons, " | ")
		if reasonStr == "" {
			reasonStr = "缘分推荐"
		}

		recommendations = append(recommendations, dto.RecommendResponse{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Tags:     u.Tags,
			Reason:   reasonStr,
			Score:    score,
		})
	}

	// 7. 排序 (按分数倒序)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// 8. 限制返回数量 (例如 Top 10)
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations, nil
}

// intersectTags 计算两个字符串切片的交集
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
