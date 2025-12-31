package org

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) DeleteByID(id uuid.UUID) error {
	return r.DB.Delete(&User{}, "id = ?", id).Error
}

func (r *UserRepository) List(offset, limit int) ([]User, int, error) {
	var users []User
	var total int64
	r.DB.Model(&User{}).Count(&total)
	if err := r.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}
