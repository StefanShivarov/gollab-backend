package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/StefanShivarov/gollab-backend/internal/backlog"
	"github.com/StefanShivarov/gollab-backend/internal/config"
	"github.com/StefanShivarov/gollab-backend/internal/org"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})
}

func PerformMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&org.User{},
		&org.Team{},
		&org.Membership{},
		&backlog.Board{},
		&backlog.Item{},
		&backlog.Tag{},
		&backlog.Comment{},
	)
}
