package repository

import (
	"go-nexus/internal/model"
	"go-nexus/pkg/global"
)

type UserRepository struct{}

var UserRepo = &UserRepository{}

func (r *UserRepository) Create(user *model.User) error {
	return global.DB.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := global.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Save(user *model.User) error {
	return global.DB.Save(user).Error
}
