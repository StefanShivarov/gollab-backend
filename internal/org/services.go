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
		return nil, common.BadRequest(err.Error())
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

type TeamService struct {
	Repo        *TeamRepository
	UserService *UserService
	Validator   *validator.Validate
}

func NewTeamService(repo *TeamRepository, userService *UserService, validator *validator.Validate) *TeamService {
	return &TeamService{
		Repo:        repo,
		UserService: userService,
		Validator:   validator,
	}
}

func (s *TeamService) List(page, size int) (*common.PaginatedResponse[TeamResponse], error) {
	offset := (page - 1) * size
	teams, total, err := s.Repo.List(offset, size)
	if err != nil {
		return nil, err
	}

	res := make([]TeamResponse, 0, len(teams))
	for _, t := range teams {
		res = append(res, *ToTeamResponse(&t))
	}

	return &common.PaginatedResponse[TeamResponse]{
		Items: res,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (s *TeamService) Create(creatorID uuid.UUID, req CreateTeamRequest) (*TeamResponse, error) {
	if err := s.Validator.Struct(req); err != nil {
		return nil, common.BadRequest(err.Error())
	}

	team := &Team{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.Repo.CreateTeamWithOwner(team, creatorID); err != nil {
		return nil, err
	}

	return ToTeamResponse(team), nil
}

func (s *TeamService) UpdateByID(id uuid.UUID, req UpdateTeamRequest) (*TeamResponse, error) {
	team, err := s.findByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.Validator.Struct(req); err != nil {
		return nil, common.BadRequest(err.Error())
	}

	if req.Name != "" {
		team.Name = req.Name
	}

	if req.Description != "" {
		team.Description = req.Description
	}

	if err := s.Repo.Update(team); err != nil {
		return nil, err
	}

	return ToTeamResponse(team), nil
}

func (s *TeamService) DeleteByID(id uuid.UUID) error {
	_, err := s.findByID(id)
	if err != nil {
		return err
	}
	return s.Repo.DeleteByID(id)
}

func (s *TeamService) GetByID(id uuid.UUID) (*TeamResponse, error) {
	team, err := s.findByID(id)
	if err != nil {
		return nil, err
	}
	return ToTeamResponse(team), nil
}

func (s *TeamService) findByID(id uuid.UUID) (*Team, error) {
	team, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFound(fmt.Sprintf("Team with id %s was not found!", id))
		}
		return nil, err
	}
	return team, nil
}

func (s *TeamService) AddMembership(request CreateMembershipRequest) error {
	if err := s.Validator.Struct(request); err != nil {
		return common.BadRequest(err.Error())
	}

	if _, err := s.findByID(request.TeamID); err != nil {
		return err
	}

	if _, err := s.UserService.findByID(request.UserID); err != nil {
		return err
	}

	m := &Membership{
		TeamID: request.TeamID,
		UserID: request.UserID,
		Role:   request.Role,
	}

	return s.Repo.AddMembership(m)
}

func (s *TeamService) RemoveMembership(teamID, userID uuid.UUID) error {
	if _, err := s.findByID(teamID); err != nil {
		return err
	}

	if _, err := s.UserService.findByID(userID); err != nil {
		return err
	}
	return s.Repo.DeleteMembershipByTeamIDAndUserID(teamID, userID)
}

func (s *TeamService) ListMembers(teamID uuid.UUID) ([]MemberResponse, error) {
	memberships, err := s.Repo.ListMembers(teamID)
	if err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return []MemberResponse{}, nil
	}

	userIds := make([]uuid.UUID, 0, len(memberships))
	for _, m := range memberships {
		userIds = append(userIds, m.UserID)
	}

	var users []User
	if err := s.Repo.DB.Where("id IN ?", userIds).Find(&users).Error; err != nil {
		return nil, err
	}

	userMap := make(map[uuid.UUID]User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	res := make([]MemberResponse, len(memberships))
	for i, m := range memberships {
		user := userMap[m.UserID]
		res[i] = MemberResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  m.Role,
		}
	}

	return res, nil
}
