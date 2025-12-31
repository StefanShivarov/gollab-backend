package org

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"username" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=150"`
}

type UpdateUserRequest struct {
	Name string `json:"username" validate:"omitempty,min=2"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"username"`
}

func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
