package org

import (
	"testing"

	"github.com/StefanShivarov/gollab-backend/internal/common"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type userRepositoryMock struct {
	mock.Mock
}

func (m *userRepositoryMock) GetByID(id uuid.UUID) (*User, error) {
	args := m.Called(id)
	if u := args.Get(0); u != nil {
		return u.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *userRepositoryMock) Create(user *User) error {
	return m.Called(user).Error(0)
}

func (m *userRepositoryMock) Update(user *User) error {
	return m.Called(user).Error(0)
}

func (m *userRepositoryMock) DeleteByID(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *userRepositoryMock) List(offset, limit int) ([]User, int, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]User), args.Int(1), args.Error(2)
}

func setupUserServiceTest() (*UserService, *userRepositoryMock, *validator.Validate) {
	mockRepo := &userRepositoryMock{}
	v := validator.New()
	service := NewUserService(mockRepo, v)
	return service, mockRepo, v
}

func TestUserService_Create(t *testing.T) {
	service, repo, _ := setupUserServiceTest()

	req := CreateUserRequest{
		Name:     "testUser",
		Email:    "test@test.com",
		Password: "testPass123",
	}

	repo.On("Create", mock.AnythingOfType("*org.User")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Email, res.Email)

	repo.AssertExpectations(t)
}

func TestUserService_Create_InvalidRequest(t *testing.T) {
	service, _, _ := setupUserServiceTest()

	req := CreateUserRequest{
		Name:     "", // invalid: empty
		Email:    "invalid-email",
		Password: "short",
	}

	res, err := service.Create(req)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUserService_GetByID(t *testing.T) {
	service, repo, _ := setupUserServiceTest()
	id := uuid.New()
	user := &User{
		BaseEntity: common.BaseEntity{
			ID: id,
		},
		Email: "test@test.com",
		Name:  "testUser",
	}

	repo.On("GetByID", id).Return(user, nil)

	res, err := service.GetByID(id)

	assert.NoError(t, err)
	assert.Equal(t, id, res.ID)
	assert.Equal(t, user.Name, res.Name)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := new(userRepositoryMock)
	service := NewUserService(repo, validator.New())

	id := uuid.New()
	repo.On("GetByID", id).Return(nil, gorm.ErrRecordNotFound)

	resp, err := service.GetByID(id)

	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestUserService_UpdateByID(t *testing.T) {
	service, repo, _ := setupUserServiceTest()

	id := uuid.New()
	user := &User{BaseEntity: common.BaseEntity{ID: id}, Name: "Old"}

	repo.On("GetByID", id).Return(user, nil)
	repo.On("Update", user).Return(nil)

	resp, err := service.UpdateByID(id, UpdateUserRequest{Name: "New"})

	assert.NoError(t, err)
	assert.Equal(t, "New", resp.Name)
}

func TestUserService_UpdateByID_ValidationError(t *testing.T) {
	service, repo, _ := setupUserServiceTest()
	id := uuid.New()
	user := &User{BaseEntity: common.BaseEntity{ID: id}, Name: "Old"}
	repo.On("GetByID", id).Return(user, nil)
	repo.On("Update", user).Return(nil)

	req := UpdateUserRequest{Name: ""} // invalid: empty string

	// Should still pass validation because "omitempty,min=2" allows empty
	resp, err := service.UpdateByID(id, req)
	assert.NoError(t, err)
	assert.Equal(t, "Old", resp.Name)
}

func TestUserService_DeleteByID(t *testing.T) {
	service, repo, _ := setupUserServiceTest()
	id := uuid.New()

	repo.On("GetByID", id).Return(&User{BaseEntity: common.BaseEntity{ID: id}}, nil)
	repo.On("DeleteByID", id).Return(nil)

	err := service.DeleteByID(id)

	assert.NoError(t, err)
}

func TestUserService_List(t *testing.T) {
	service, repoMock, _ := setupUserServiceTest()

	users := []User{
		{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "Alice", Email: "alice@test.com"},
		{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "Bob", Email: "bob@test.com"},
	}
	repoMock.On("List", 0, 2).Return(users, 2, nil)

	resp, err := service.List(1, 2)
	assert.NoError(t, err)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, "Alice", resp.Items[0].Name)
	assert.Equal(t, "Bob", resp.Items[1].Name)
	repoMock.AssertExpectations(t)
}

type teamRepositoryMock struct {
	mock.Mock
}

func (m *teamRepositoryMock) GetByID(id uuid.UUID) (*Team, error) {
	args := m.Called(id)
	if t := args.Get(0); t != nil {
		return t.(*Team), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *teamRepositoryMock) Update(team *Team) error {
	return m.Called(team).Error(0)
}

func (m *teamRepositoryMock) DeleteByID(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *teamRepositoryMock) CreateTeamWithOwner(team *Team, creatorID uuid.UUID) error {
	return m.Called(team, creatorID).Error(0)
}

func (m *teamRepositoryMock) AddMembership(mem *Membership) error {
	return m.Called(mem).Error(0)
}

func (m *teamRepositoryMock) DeleteMembershipByTeamIDAndUserID(teamID, userID uuid.UUID) error {
	return m.Called(teamID, userID).Error(0)
}

func (m *teamRepositoryMock) ListMembers(teamID uuid.UUID) ([]MemberResponse, error) {
	args := m.Called(teamID)
	return args.Get(0).([]MemberResponse), args.Error(1)
}

func (m *teamRepositoryMock) List(offset, limit int) ([]Team, int, error) {
	args := m.Called(offset, limit)
	teams, _ := args.Get(0).([]Team)
	total := args.Int(1)
	return teams, total, args.Error(2)
}

func setupTeamServiceTest() (*TeamService, *teamRepositoryMock, *UserService, *userRepositoryMock, *validator.Validate) {
	v := validator.New()
	userRepoMock := &userRepositoryMock{}
	userService := NewUserService(userRepoMock, v)
	teamRepoMock := &teamRepositoryMock{}
	teamService := NewTeamService(teamRepoMock, userService, v)
	return teamService, teamRepoMock, userService, userRepoMock, v
}

func TestTeamService_Create(t *testing.T) {
	service, repo, _, _, _ := setupTeamServiceTest()

	creatorID := uuid.New()
	req := CreateTeamRequest{Name: "Team A"}

	repo.On("CreateTeamWithOwner", mock.AnythingOfType("*org.Team"), creatorID).Return(nil)

	resp, err := service.Create(creatorID, req)

	assert.NoError(t, err)
	assert.Equal(t, "Team A", resp.Name)
}

func TestTeamService_Create_ValidationError(t *testing.T) {
	service, _, _, _, _ := setupTeamServiceTest()
	creatorID := uuid.New()
	req := CreateTeamRequest{Name: ""} // invalid: required

	resp, err := service.Create(creatorID, req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestTeamService_GetByID(t *testing.T) {
	service, repoMock, _, _, _ := setupTeamServiceTest()
	team := &Team{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "TeamX", Description: "Desc"}
	repoMock.On("GetByID", team.ID).Return(team, nil)

	resp, err := service.GetByID(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "TeamX", resp.Name)
	repoMock.AssertExpectations(t)
}

func TestTeamService_GetByID_NotFound(t *testing.T) {
	service, repo, _, _, _ := setupTeamServiceTest()
	id := uuid.New()
	repo.On("GetByID", id).Return(nil, gorm.ErrRecordNotFound)

	resp, err := service.GetByID(id)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestTeamService_UpdateByID(t *testing.T) {
	service, repoMock, _, _, _ := setupTeamServiceTest()
	team := &Team{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "OldName", Description: "OldDesc"}
	repoMock.On("GetByID", team.ID).Return(team, nil)
	repoMock.On("Update", team).Return(nil)

	req := UpdateTeamRequest{Name: "NewName", Description: "NewDesc"}
	resp, err := service.UpdateByID(team.ID, req)
	assert.NoError(t, err)
	assert.Equal(t, "NewName", resp.Name)
	assert.Equal(t, "NewDesc", resp.Description)
	repoMock.AssertExpectations(t)
}

func TestTeamService_DeleteByID(t *testing.T) {
	service, repoMock, _, _, _ := setupTeamServiceTest()
	team := &Team{BaseEntity: common.BaseEntity{ID: uuid.New()}}
	repoMock.On("GetByID", team.ID).Return(team, nil)
	repoMock.On("DeleteByID", team.ID).Return(nil)

	err := service.DeleteByID(team.ID)
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestTeamService_List(t *testing.T) {
	service, repoMock, _, _, _ := setupTeamServiceTest()
	teams := []Team{
		{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "Team1"},
		{BaseEntity: common.BaseEntity{ID: uuid.New()}, Name: "Team2"},
	}
	repoMock.On("List", 0, 2).Return(teams, 2, nil)

	resp, total, err := service.Repo.List(0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Team1", resp[0].Name)
	assert.Equal(t, "Team2", resp[1].Name)
	repoMock.AssertExpectations(t)
}

func TestTeamService_AddMembership(t *testing.T) {
	service, teamRepo, _, userRepo, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()

	teamRepo.On("GetByID", teamID).Return(&Team{BaseEntity: common.BaseEntity{ID: teamID}}, nil)
	userRepo.On("GetByID", userID).Return(&User{BaseEntity: common.BaseEntity{ID: userID}}, nil)
	teamRepo.On("AddMembership", mock.AnythingOfType("*org.Membership")).Return(nil)

	err := service.AddMembership(CreateMembershipRequest{
		TeamID: teamID,
		UserID: userID,
		Role:   Developer,
	})

	assert.NoError(t, err)
}

func TestTeamService_AddMembership_TeamNotFound(t *testing.T) {
	service, teamRepo, _, userRepo, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()

	teamRepo.On("GetByID", teamID).Return(nil, gorm.ErrRecordNotFound)
	userRepo.On("GetByID", userID).Return(&User{BaseEntity: common.BaseEntity{ID: userID}}, nil)

	err := service.AddMembership(CreateMembershipRequest{
		TeamID: teamID,
		UserID: userID,
		Role:   Developer,
	})
	assert.Error(t, err)
}

func TestTeamService_AddMembership_UserNotFound(t *testing.T) {
	service, teamRepo, _, userRepo, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()

	teamRepo.On("GetByID", teamID).Return(&Team{BaseEntity: common.BaseEntity{ID: teamID}}, nil)
	userRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	err := service.AddMembership(CreateMembershipRequest{
		TeamID: teamID,
		UserID: userID,
		Role:   Developer,
	})
	assert.Error(t, err)
}

func TestTeamService_AddMembership_ValidationError(t *testing.T) {
	service, _, _, _, _ := setupTeamServiceTest()
	req := CreateMembershipRequest{}

	err := service.AddMembership(req)
	assert.Error(t, err)
}

func TestTeamService_RemoveMembership(t *testing.T) {
	service, teamRepo, _, userRepo, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()

	teamRepo.On("GetByID", teamID).Return(&Team{BaseEntity: common.BaseEntity{ID: teamID}}, nil)
	userRepo.On("GetByID", userID).Return(&User{BaseEntity: common.BaseEntity{ID: userID}}, nil)
	teamRepo.On("DeleteMembershipByTeamIDAndUserID", teamID, userID).Return(nil)

	err := service.RemoveMembership(teamID, userID)

	assert.NoError(t, err)
}

func TestTeamService_RemoveMembership_TeamNotFound(t *testing.T) {
	service, teamRepo, _, _, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()
	teamRepo.On("GetByID", teamID).Return(nil, gorm.ErrRecordNotFound)

	err := service.RemoveMembership(teamID, userID)
	assert.Error(t, err)
}

func TestTeamService_RemoveMembership_UserNotFound(t *testing.T) {
	service, teamRepo, _, userRepo, _ := setupTeamServiceTest()
	teamID := uuid.New()
	userID := uuid.New()
	teamRepo.On("GetByID", teamID).Return(&Team{BaseEntity: common.BaseEntity{ID: teamID}}, nil)
	userRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	err := service.RemoveMembership(teamID, userID)
	assert.Error(t, err)
}

func TestTeamService_ListMembers(t *testing.T) {
	service, repo, _, _, _ := setupTeamServiceTest()

	teamID := uuid.New()
	members := []MemberResponse{
		{ID: uuid.New(), Name: "Alice", Role: Developer},
	}

	repo.On("ListMembers", teamID).Return(members, nil)

	resp, err := service.ListMembers(teamID)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Alice", resp[0].Name)
}
