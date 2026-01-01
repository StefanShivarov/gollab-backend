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

type TeamRepository struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{DB: db}
}

func (r *TeamRepository) GetByID(id uuid.UUID) (*Team, error) {
	var team Team
	if err := r.DB.First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) Update(team *Team) error {
	return r.DB.Save(team).Error
}

func (r *TeamRepository) DeleteByID(id uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Membership{}, "team_id = ?", id).Error; err != nil {
			return err
		}

		return tx.Delete(&Team{}, "id = ?", id).Error
	})
}

func (r *TeamRepository) List(offset, limit int) ([]Team, int, error) {
	var teams []Team
	var total int64
	r.DB.Model(&Team{}).Count(&total)
	if err := r.DB.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	return teams, int(total), nil
}

func (r *TeamRepository) AddMembership(m *Membership) error {
	return r.DB.Create(m).Error
}

func (r *TeamRepository) DeleteMembershipByTeamIDAndUserID(teamID, userID uuid.UUID) error {
	return r.DB.Delete(&Membership{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

func (r *TeamRepository) ListMembers(teamID uuid.UUID) ([]Membership, error) {
	var m []Membership
	err := r.DB.Where("team_id = ?", teamID).Find(&m).Error
	return m, err
}

func (r *TeamRepository) CreateTeamWithOwner(team *Team, creatorID uuid.UUID) error {
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
