package service

import (
	"errors"
	"go-nexus/internal/model"
	"go-nexus/internal/model/vo"
	"go-nexus/internal/repository"
	"go-nexus/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// 在 service 层定义错误类型，Controller 直接用，不用解析字符串
var (
	ErrUserExist    = errors.New("该用户名已被注册")
	ErrUserNotFound = errors.New("用户不存在")
	ErrWrongPwd     = errors.New("密码错误")
	ErrTokenFailed  = errors.New("token生成失败")
)

type UserService struct{}

var UserSvc = &UserService{}

func (s *UserService) Register(username, password string) error {
	_, err := repository.UserRepo.GetByUsername(username)
	if err == nil {
		return ErrUserExist // 直接返回有类型的错误
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // 数据库错误直接往上抛，不包装字符串
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return repository.UserRepo.Create(&model.User{
		Username:      username,
		Password:      hashed,
		Nickname:      username,
		LastLoginTime: time.Now(),
		Birthday:      time.Now(),
	})
}

func (s *UserService) Login(username, password string) (*model.User, string, error) {
	user, err := repository.UserRepo.GetByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", ErrUserNotFound
	}
	if err != nil {
		return nil, "", err
	}
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", ErrWrongPwd
	}
	token, err := utils.GenToken(user.ID, user.Username)
	if err != nil {
		return nil, "", ErrTokenFailed
	}
	return user, token, nil
}

func (s *UserService) GetUserInfo(userID uint) (*model.UserProfileResponse, error) {
	user, err := repository.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return &model.UserProfileResponse{
		ID:            user.ID,
		Username:      user.Username,
		Nickname:      user.Nickname,
		CreatedAt:     user.CreatedAt,
		LastLoginTime: user.LastLoginTime,
		Avatar:        user.Avatar,
		Email:         user.Email,
		Tags:          user.Tags,
		Signature:     user.Signature,
		Birthday:      user.Birthday,
		Location:      user.Location,
	}, nil
}

func (s *UserService) UpdateUserInfo(info vo.UpdateUser, userID uint) error {
	user, err := repository.UserRepo.GetByID(userID)
	if err != nil {
		return err
	}
	user.Nickname = info.Nickname
	user.Email = info.Email
	user.Gender = info.Gender
	user.Birthday = info.Birthday
	user.Signature = info.Signature
	user.Location = info.Location
	if info.Avatar != "" {
		user.Avatar = info.Avatar
	}
	return repository.UserRepo.Save(user)
}
func (s *UserService) UpdateAvatar(userID uint, avatar string) error {
	user, err := repository.UserRepo.GetByID(userID)
	if err != nil {
		return err
	}
	user.Avatar = avatar
	return repository.UserRepo.Save(user)
}
