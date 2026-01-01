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

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=50"`
	Description string `json:"description"`
}

type TeamResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func ToTeamResponse(team *Team) *TeamResponse {
	return &TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
	}
}

type CreateMembershipRequest struct {
	TeamID uuid.UUID `json:"teamId" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
	Role   TeamRole  `json:"role" validate:"required,oneof=project_manager developer"`
}

type DeleteMembershipRequest struct {
	TeamID uuid.UUID `json:"teamId" validate:"required,uuid"`
	UserID uuid.UUID `json:"userId" validate:"required,uuid"`
}

type MemberResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"username"`
	Email string    `json:"email"`
	Role  TeamRole  `json:"role"`
}
