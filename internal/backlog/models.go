package backlog

import (
	"time"

	"github.com/StefanShivarov/gollab-backend/internal/common"
	"github.com/StefanShivarov/gollab-backend/internal/org"
	"github.com/google/uuid"
)

type Board struct {
	common.BaseEntity
	Name        string    `gorm:"not null;uniqueIndex:idx_team_board_name"`
	Description string    `gorm:"type:text;not null"`
	TeamID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_board_name"`
	Team        org.Team  `gorm:"foreignKey:TeamID"`
}

type ItemStatus string

const (
	NotPlanned ItemStatus = "not_planned"
	ToDo       ItemStatus = "to_do"
	InProgress ItemStatus = "in_progress"
	OnHold     ItemStatus = "on_hold"
	InReview   ItemStatus = "in_review"
	Done       ItemStatus = "done"
)

type Item struct {
	common.BaseEntity
	DueDate     *time.Time
	Title       string     `gorm:"type:varchar(50);not null"`
	Description string     `gorm:"type:text"`
	Status      ItemStatus `gorm:"type:varchar(20);not null;default:'not_planned';check: status IN ('not_planned', 'to_do', 'in_progress', 'on_hold', 'in_review', 'done')"`
	Priority    int        `gorm:"default:0"`
	AuthorID    uuid.UUID  `gorm:"type:uuid;not null"`
	Author      org.User   `gorm:"foreignKey:AuthorID"`
	BoardID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Board       Board      `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags        []Tag      `gorm:"many2many:items_tags"`
	Assignees   []org.User `gorm:"many2many:items_assignees"`
}

type Tag struct {
	common.BaseEntity
	TeamID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_tag_name"`
	Team   org.Team  `gorm:"foreignKey:TeamID"`
	Name   string    `gorm:"type:varchar(30);not null;uniqueIndex:idx_team_tag_name"`
	Color  string    `gorm:"type:varchar(20);not null"`
}

type Comment struct {
	common.BaseEntity
	Content string    `gorm:"type:text;not null"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;"`
	User    org.User  `gorm:"foreignKey:UserID;onUpdate:CASCADE,onDelete:CASCADE"`
	ItemID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Item    Item      `gorm:"foreignKey:ItemID;onUpdate:CASCADE,onDelete:CASCADE"`
}
