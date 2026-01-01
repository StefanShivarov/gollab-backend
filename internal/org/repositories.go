package org

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uuid.UUID) (*User, error)
	Create(user *User) error
	Update(user *User) error
	DeleteByID(id uuid.UUID) error
	List(offset, limit int) ([]User, int, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) Update(user *User) error {
	return r.DB.Save(user).Error
}

func (r *userRepository) DeleteByID(id uuid.UUID) error {
	return r.DB.Delete(&User{}, "id = ?", id).Error
}

func (r *userRepository) List(offset, limit int) ([]User, int, error) {
	var users []User
	var total int64
	r.DB.Model(&User{}).Count(&total)
	if err := r.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

type TeamRepository interface {
	GetByID(id uuid.UUID) (*Team, error)
	Update(team *Team) error
	DeleteByID(id uuid.UUID) error
	List(offset, limit int) ([]Team, int, error)
	CreateTeamWithOwner(team *Team, creatorId uuid.UUID) error
	AddMembership(membership *Membership) error
	DeleteMembershipByTeamIDAndUserID(teamID uuid.UUID, userID uuid.UUID) error
	ListMembers(teamID uuid.UUID) ([]MemberResponse, error)
}

type teamRepository struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{DB: db}
}

func (r *teamRepository) GetByID(id uuid.UUID) (*Team, error) {
	var team Team
	if err := r.DB.First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Update(team *Team) error {
	return r.DB.Save(team).Error
}

func (r *teamRepository) DeleteByID(id uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Membership{}, "team_id = ?", id).Error; err != nil {
			return err
		}

		return tx.Delete(&Team{}, "id = ?", id).Error
	})
}

func (r *teamRepository) List(offset, limit int) ([]Team, int, error) {
	var teams []Team
	var total int64
	r.DB.Model(&Team{}).Count(&total)
	if err := r.DB.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	return teams, int(total), nil
}

func (r *teamRepository) AddMembership(m *Membership) error {
	return r.DB.Create(m).Error
}

func (r *teamRepository) DeleteMembershipByTeamIDAndUserID(teamID, userID uuid.UUID) error {
	return r.DB.Delete(&Membership{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

func (r *teamRepository) ListMembers(teamID uuid.UUID) ([]MemberResponse, error) {
	var res []MemberResponse
	err := r.DB.
		Table("memberships").
		Select("users.id as user_id, users.name, users.email, memberships.role").
		Joins("JOIN users ON users.id = memberships.user_id").
		Where("memberships.team_id = ?", teamID).
		Scan(&res).Error
	return res, err
}

func (r *teamRepository) CreateTeamWithOwner(team *Team, creatorID uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		m := &Membership{
			TeamID: team.ID,
			UserID: creatorID,
			Role:   ProjectManager,
		}
		return tx.Create(m).Error
	})
}
