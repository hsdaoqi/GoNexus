package v1

import (
	"go-nexus/internal/model/vo"
	"go-nexus/internal/repository"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// Register 用户注册接口
func Register(c *gin.Context) {
	// 1. 参数绑定和验证
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "参数错误："+err.Error())
		return
	}

	// 2. 调用业务逻辑
	err := service.Register(req.Username, req.Password)
	if err != nil {
		// 根据错误信息判断错误类型
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "user_exist"):
			response.Fail(c, response.ErrUserExist)
		case strings.Contains(errMsg, "db_error"):
			response.Error(c, http.StatusInternalServerError, response.ErrSystemError, "数据库错误")
		default:
			response.Error(c, http.StatusInternalServerError, response.ErrSystemError, "注册失败")
		}
		return
	}

	// 3. 注册成功
	response.Success(c, nil)
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	user, token, err := service.Login(req.Username, req.Password)
	if err != nil {
		// 这里根据错误类型返回不同的错误码会更完美，暂时统一返回密码错误
		response.FailWithMessage(c, response.ErrPasswordError, err.Error())
		return
	}

	// 返回数据：包含 Token 和 用户基础信息
	response.Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
	})
}

// GetUserInfo 获取当前登录用户的个人信息
func GetUserInfo(c *gin.Context) {
	// 1. 从上下文获取 userID
	value, exists := c.Get("userID")
	if !exists {
		// 这种情况一般不会发生，除非中间件出 Bug 了
		response.Fail(c, response.ErrSystemError)
		return
	}

	userID := value.(uint)

	// 2. 调用 Service (真正的业务逻辑)
	data, err := service.GetUserInfo(userID)
	if err != nil {
		// 查不到用户（可能Token有效，但号被删了）
		response.FailWithMessage(c, response.ErrUserNotExist, "用户不存在或已注销")
		return
	}
	// 3. 返回干净的数据
	response.Success(c, data)
}

func UpdateAvatar(c *gin.Context) {
	// 1. 获取文件 (前端字段名叫 "file")
	file, err := c.FormFile("avatar")
	if err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "请上传头像")
		return
	}

	// 2. 调用服务
	avatar, err := service.UploadFile(file)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "上传失败: "+err.Error())
		return
	}
	//3.更改用户avatar
	userID, _ := c.Get("userID")
	id := userID.(uint)
	user, _ := repository.GetUserByID(id)
	user.Avatar = avatar
	err = repository.SaveUser(user)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
	}
	// 3. 返回 URL
	response.Success(c, gin.H{"avatar": avatar})
}

func UpdateUserInfo(c *gin.Context) {
	userID, exits := c.Get("userID")
	if !exits {
		response.Fail(c, response.ErrSystemError)
		return
	}
	id := userID.(uint)
	userinfo := vo.UpdateUser{}
	err := c.ShouldBindJSON(&userinfo)
	if err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, err.Error())
		return
	}
	err = service.UpdateUserInfo(userinfo, id)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
	}
	response.Success(c, gin.H{"msg": "更新成功"})
}
