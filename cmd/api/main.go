package main

import (
	"log"
	"main/cmd/api/docs"
	"main/pkg/config"
	"main/pkg/db"
	di "main/pkg/di"
	"main/pkg/domain"
	"time"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func main() {

	docs.SwaggerInfo.Title = "GameVerse"

	docs.SwaggerInfo.Version = "1.0"

	//docs.SwaggerInfo.Host = "localhost:1245"
	docs.SwaggerInfo.Host = "gameverse.cloud"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)

	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {

		server.Start()
	}
	database, dbErr := db.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("cannot connect to the database: ", dbErr)
	}
	// Start the cron job for updating subscription status
	startSubscriptionUpdateJob(database)
}

// cron job functions for subscription updates...
func startSubscriptionUpdateJob(database *gorm.DB) {
	// Schedule the job using AddFunc
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		// Inside the scheduled function, call UpdateSubscriptionStatus
		if err := UpdateSubscriptionStatus(database); err != nil {
			log.Println("Error updating subscription statuses:", err)
		}
	})

	// Start the scheduler
	c.Start()

	// Block the main goroutine to keep the application running
	select {}
}

func UpdateSubscriptionStatus(database *gorm.DB) error {
	// Fetch active subscriptions from the database
	// Use the 'database' instance directly
	activeSubscriptions, err := FetchActiveSubscriptions(database)
	if err != nil {
		return err
	}

	// Update subscription statuses based on expiration
	for _, subscription := range activeSubscriptions {
		if HasSubscriptionExpired(subscription.SubscribedAt, subscription.SubscriptionPlan.Duration) {
			subscription.IsActive = false
			// Update the subscription in the database
			err := UpdateSubscriptionInDatabase(database, subscription)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func FetchActiveSubscriptions(database *gorm.DB) ([]domain.SubscriptionList, error) {
	var activeSubscriptions []domain.SubscriptionList
	err := database.Where("is_active = ?", true).Find(&activeSubscriptions).Error
	return activeSubscriptions, err
}

func UpdateSubscriptionInDatabase(database *gorm.DB, subscription domain.SubscriptionList) error {
	return database.Save(&subscription).Error
}

func HasSubscriptionExpired(subscribedAt time.Time, duration int) bool {
	expirationTime := subscribedAt.Add(time.Duration(duration) * 24 * time.Hour)
	return time.Now().After(expirationTime)
}
