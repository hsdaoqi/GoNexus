package repository

import (
	"go-nexus/internal/model"
	"go-nexus/pkg/global"
)

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return global.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查找用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserById 根据ID查询用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := global.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

// SaveUser 存储用户
func SaveUser(user *model.User) error {
	return global.DB.Save(user).Error
}
