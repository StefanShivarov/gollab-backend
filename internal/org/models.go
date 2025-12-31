package org

import (
	"github.com/StefanShivarov/gollab-backend/internal/common"
	"github.com/google/uuid"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	Standard UserRole = "standard"
)

type User struct {
	common.BaseEntity
	Email        string   `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name         string   `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string   `gorm:"type:varchar(150);not null"`
	Role         UserRole `gorm:"type:varchar(20);not null;default:'standard';check:role IN ('admin', 'standard')"`
}

type TeamRole string

const (
	ProjectManager TeamRole = "project_manager"
	Developer      TeamRole = "developer"
)

type Team struct {
	common.BaseEntity
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text"`
}

type Membership struct {
	common.BaseEntity
	UserID uuid.UUID `gorm:"type:uuid; not null;index;uniqueIndex:idx_user_team_membership"`
	User   User      `gorm:"foreignKey:UserID"`
	TeamID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_team_membership"`
	Team   Team      `gorm:"foreignKey:TeamID"`
	Role   TeamRole  `gorm:"type:varchar(20);not null;default:'developer';check: role IN ('project_manager', 'developer')"`
}
