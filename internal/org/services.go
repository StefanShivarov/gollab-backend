package org

import (
	"errors"
	"fmt"

	"github.com/StefanShivarov/gollab-backend/internal/common"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Repo      *UserRepository
	Validator *validator.Validate
}

func NewUserService(repo *UserRepository, validator *validator.Validate) *UserService {
	return &UserService{
		Repo:      repo,
		Validator: validator,
	}
}

func (s *UserService) Create(req CreateUserRequest) (*UserResponse, error) {
	if err := s.Validator.Struct(req); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		Role:         Standard,
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (s *UserService) findByID(id uuid.UUID) (*User, error) {
	user, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFound(fmt.Sprintf("User with id %s was not found!", id))
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByID(id uuid.UUID) (*UserResponse, error) {
	user, err := s.findByID(id)
	if err != nil {
		return nil, err
	}
	return ToUserResponse(user), nil
}

func (s *UserService) UpdateByID(id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.findByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.Validator.Struct(req); err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.Repo.Update(user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (s *UserService) DeleteByID(id uuid.UUID) error {
	_, err := s.findByID(id)
	if err != nil {
		return err
	}
	return s.Repo.DeleteByID(id)
}

func (s *UserService) List(page, size int) (*common.PaginatedResponse[UserResponse], error) {
	offset := (page - 1) * size
	users, total, err := s.Repo.List(offset, size)
	if err != nil {
		return nil, err
	}

	res := make([]UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, *ToUserResponse(&u))
	}

	return &common.PaginatedResponse[UserResponse]{
		Items: res,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}
