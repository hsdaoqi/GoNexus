package service

import (
	"errors"
	"fmt"
	"go-nexus/internal/model"
	"go-nexus/internal/model/vo"
	"go-nexus/internal/repository"
	"go-nexus/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// Register 注册业务
func Register(username, password string) error {
	//1.检查用户名是否已存在
	_, err := repository.GetUserByUsername(username)
	if err == nil { //没报错说明用户存在
		return fmt.Errorf("user_exist") // 使用特定错误标识
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("db_error: %w", err) // 数据库错误
	}

	//2，密码加密
	password, err = utils.HashPassword(password)
	if err != nil {
		return err
	}

	//3.创建用户
	user := &model.User{
		Username:      username,
		Password:      password,
		Nickname:      username, //昵称默认为用户名
		LastLoginTime: time.Now(),
		Birthday:      time.Now(),
	}

	//4.存库
	return repository.CreateUser(user)
}

// Login 用户登录
func Login(username, password string) (*model.User, string, error) {
	//根据用户名找到用户
	user, err := repository.GetUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", errors.New("用户不存在,请检查用户名是否正确")
	}
	if err != nil {
		return nil, "", err
	}
	//2.校验密码
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", errors.New("密码错误")
	}
	//3.生成token
	token, err := utils.GenToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("token生成失败")
	}
	return user, token, nil
}

// GetUserInfo 获取个人信息服务
func GetUserInfo(userID uint) (*model.UserProfileResponse, error) {
	// 1. 调库查数据
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 2. 数据清洗 (Model -> DTO)
	// 这一步非常关键！我们不返回 user.Password
	resp := &model.UserProfileResponse{
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
	}
	return resp, nil
}

func UpdateUserInfo(userinfo vo.UpdateUser, userID uint) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.Nickname = userinfo.Nickname
	user.Email = userinfo.Email
	user.Gender = userinfo.Gender
	user.Birthday = userinfo.Birthday
	user.Signature = userinfo.Signature
	user.Location = userinfo.Location
	err = repository.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}
