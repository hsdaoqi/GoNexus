package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-nexus/internal/model/vo"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, err.Error())
		return
	}
	if err := service.UserSvc.Register(req.Username, req.Password); err != nil {
		switch {
		case errors.Is(err, service.ErrUserExist):
			response.Fail(c, response.ErrUserExist)
		default:
			response.Fail(c, response.ErrSystemError)
		}
		return
	}
	response.Success(c, nil)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	user, token, err := service.UserSvc.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Fail(c, response.ErrUserNotExist)
		case errors.Is(err, service.ErrWrongPwd):
			response.Fail(c, response.ErrPasswordError)
		default:
			response.Fail(c, response.ErrSystemError)
		}
		return
	}
	response.Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	data, err := service.UserSvc.GetUserInfo(userID)
	if err != nil {
		response.Fail(c, response.ErrUserNotExist)
		return
	}
	response.Success(c, data)
}

func UpdateAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "请上传头像")
		return
	}
	avatar, err := service.UploadFile(file)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "上传失败: "+err.Error())
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.UserSvc.UpdateAvatar(userID, avatar); err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, gin.H{"avatar": avatar})
}
func UpdateUserInfo(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var info vo.UpdateUser
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, err.Error())
		return
	}
	if err := service.UserSvc.UpdateUserInfo(info, userID); err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, gin.H{"msg": "更新成功"})
}
