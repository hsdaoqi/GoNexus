package v1

import (
	"go-nexus/internal/model/dto"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost 发布动态
func CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if err := service.MomentServiceApp.CreatePost(userID, req); err != nil {
		response.FailWithMessage(c, response.ErrBusiness, "发布失败: "+err.Error())
		return
	}

	response.Success(c, "发布成功")
}

// GetMoments 获取动态列表
func GetMoments(c *gin.Context) {
	var req dto.GetMomentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误")
		return
	}

	// 默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Type == "" {
		req.Type = "all"
	}

	userID := c.GetUint("userID")
	list, total, err := service.MomentServiceApp.GetMoments(userID, req)
	if err != nil {
		response.FailWithMessage(c, response.ErrBusiness, "获取失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// CreateComment 发表评论
func CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误")
		return
	}

	userID := c.GetUint("userID")
	res, err := service.MomentServiceApp.CreateComment(userID, req)
	if err != nil {
		response.FailWithMessage(c, response.ErrBusiness, "评论失败: "+err.Error())
		return
	}

	response.Success(c, res)
}

// GetComments 获取评论
func GetComments(c *gin.Context) {
	postIDStr := c.Query("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误")
		return
	}

	list, err := service.MomentServiceApp.GetComments(uint(postID))
	if err != nil {
		response.FailWithMessage(c, response.ErrBusiness, "获取评论失败")
		return
	}

	response.Success(c, list)
}

// ToggleLike 点赞/取消
func ToggleLike(c *gin.Context) {
	var req dto.LikePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误")
		return
	}

	userID := c.GetUint("userID")
	isLiked, likeCount, err := service.MomentServiceApp.ToggleLike(userID, req)
	if err != nil {
		response.FailWithMessage(c, response.ErrBusiness, "操作失败: "+err.Error())
		return
	}

	msg := "点赞成功"
	if !isLiked {
		msg = "取消点赞"
	}
	response.Success(c, gin.H{"is_liked": isLiked, "like_count": likeCount, "msg": msg})
}
