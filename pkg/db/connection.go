package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "main/pkg/config"
	domain "main/pkg/domain"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.Reports{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Video{})
	db.AutoMigrate(&domain.VideoLikes{})
	db.AutoMigrate(&domain.Tag{})
	db.AutoMigrate(&domain.UserTags{})
	db.AutoMigrate(&domain.Comment{})
	db.AutoMigrate(&domain.VideoTags{})
	db.AutoMigrate(&domain.SubscriptionPlan{})
	return db, dbErr
}
