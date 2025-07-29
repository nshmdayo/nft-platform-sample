package database

import (
	"github.com/nshmdayo/nft-platform-sample/internal/models"
	"github.com/nshmdayo/nft-platform-sample/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(databaseURL string) error {
	var err error

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return err
	}

	logger.Info("Successfully connected to database")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Paper{},
		&models.Review{},
		&models.NFTMetadata{},
	)

	if err != nil {
		logger.Error("Failed to auto migrate database", "error", err)
		return err
	}

	logger.Info("Database migration completed successfully")
	return nil
}
