package repository

import (
	"user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Count() int64 {
	var count int64
	r.DB.Model(&model.User{}).Count(&count)
	return count
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}