package repository

import (
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewsubscriptionRepository(DB *gorm.DB) interfaces.SubscriptionRepository {
	return &SubscriptionRepository{DB}
}

func (r *SubscriptionRepository) GetSubscriptionListIDByPlanID(planID, creatorID, userID int) (int, error) {
	// Check if the user already has an active subscription
	// existingSubscription, err := r.GetActiveSubscription(planID, creatorID, userID)
	// if err != nil {
	// 	return 0, err
	// }

	// if existingSubscription != nil {
	// 	return 0, errors.New("user already subscribed")
	// }

	// Create a new subscription
	subscriptionList := domain.SubscriptionList{
		CreatorID: creatorID,
		UserID:    userID,
		PlanID:    planID,
		//IsActive:  true, // Assuming a new subscription is active by default
	}

	result := r.DB.Create(&subscriptionList)

	if result.Error != nil {
		return 0, result.Error
	}

	return subscriptionList.ID, nil
}

func (r *SubscriptionRepository) GetActiveSubscription(creatorID, userID int) (*domain.SubscriptionList, error) {
	var subscription domain.SubscriptionList

	// Specify the fields you need in the SELECT statement
	rows, err := r.DB.Raw("SELECT creator_id, user_id, plan_id, is_active FROM subscription_lists WHERE creator_id = ? AND user_id = ? AND is_active = true ORDER BY subscribed_at DESC LIMIT 1", creatorID, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use Scan to populate the subscription struct
	if rows.Next() {
		if err := r.DB.ScanRows(rows, &subscription); err != nil {
			return nil, err
		}
	} else {
		// No active subscription found
		return nil, nil
	}

	return &subscription, nil
}

// func (r *SubscriptionRepository) GetSubscriptionListIDByPlanID(planID, creatorID, userID int) (int, error) {
// 	subscriptionList := domain.SubscriptionList{
// 		CreatorID: creatorID,
// 		UserID:    userID,
// 		PlanID:    planID,
// 	}

// 	result := r.DB.Create(&subscriptionList)

// 	if result.Error != nil {
// 		return 0, result.Error
// 	}

// 	return subscriptionList.ID, nil
// }

func (p *SubscriptionRepository) FindUsername(user_id int) (string, error) {
	var username string
	if err := p.DB.Raw("SELECT username FROM users WHERE id=?", user_id).Scan(&username).Error; err != nil {
		return "", err
	}

	return username, nil
}
func (p *SubscriptionRepository) FindPrice(orderID int) (float64, error) {
	// Step 1: Find the planID using the orderID
	var planID int
	if err := p.DB.Raw("SELECT plan_id FROM subscription_lists WHERE id=?", orderID).Scan(&planID).Error; err != nil {
		return 0, err
	}

	// Step 2: Find the price using the obtained planID
	var price float64
	if err := p.DB.Raw("SELECT price FROM subscription_plans WHERE id=?", planID).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil
}

// func (p *SubscriptionRepository) UpdatePaymentDetails(orderID, paymentID, razorID string) error {
// 	status := "PAID"
// 	subscribedAt := time.Now() // Update the SubscribedAt field to the current time

// 	if err := p.DB.Exec(`
// 		UPDATE subscription_lists
// 		SET payment_status = $1, payment_id = $3, subscribed_at = $4
// 		WHERE id = $2`, status, orderID, paymentID, subscribedAt).Error; err != nil {
// 		return err
// 	}

//		return nil
//	}
func (p *SubscriptionRepository) UpdatePaymentDetails(orderID, paymentID, razorID string) error {
	status := "PAID"
	subscribedAt := time.Now() // Update the SubscribedAt field to the current time

	if err := p.DB.Exec(`
		UPDATE subscription_lists 
		SET payment_status = $1, payment_id = $3, subscribed_at = $4, is_active = true
		WHERE id = $2`, status, orderID, paymentID, subscribedAt).Error; err != nil {
		return err
	}

	return nil
}

// func (s *SubscriptionRepository) FetchActiveSubscriptions() ([]domain.SubscriptionList, error) {
// 	var activeSubscriptions []domain.SubscriptionList
// 	err := s.DB.Where("is_active = ?", true).Find(&activeSubscriptions).Error
// 	return activeSubscriptions, err
// }

// func (s *SubscriptionRepository) UpdateSubscriptionInDatabase(subscription domain.SubscriptionList) error {
// 	return s.DB.Save(&subscription).Error
// }

// func (s *SubscriptionRepository) HasSubscriptionExpired(subscribedAt time.Time, duration int) bool {
// 	expirationTime := subscribedAt.Add(time.Duration(duration) * 24 * time.Hour)
// 	return time.Now().After(expirationTime)
// }
// func (s *SubscriptionRepository) UpdateSubscriptionStatus() error {

// 	// Fetch active subscriptions from the database
// 	activeSubscriptions, err := s.FetchActiveSubscriptions()
// 	if err != nil {
// 		return err
// 	}

// 	// Update subscription statuses based on expiration
// 	for _, subscription := range activeSubscriptions {
// 		if s.HasSubscriptionExpired(subscription.SubscribedAt, subscription.SubscriptionPlan.Duration) {
// 			subscription.IsActive = false
// 			// Update the subscription in the database
// 			err := s.UpdateSubscriptionInDatabase(subscription)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

//		return nil
//	}
func (r *SubscriptionRepository) GetRevenue(creatorID int, startDate string, endDate string) (float64, error) {
	var revenue float64

	query := r.DB.
		Model(&domain.SubscriptionList{}).
		Select("COALESCE(SUM(subscription_plans.price), 0)").
		Joins("JOIN subscription_plans ON subscription_lists.plan_id = subscription_plans.id").
		Where("subscription_lists.creator_id = ?", creatorID)

	if startDate != "" && endDate != "" {
		query = query.Where("subscription_lists.subscribed_at BETWEEN ? AND ?", startDate, endDate)
	}

	result := query.Scan(&revenue)
	if result.Error != nil {
		return 0, result.Error
	}

	return revenue, nil
}
func (r *SubscriptionRepository) GetSubscribersCount(creatorID int, startDate string, endDate string) (int, error) {
	var count int

	query := "SELECT COUNT(id) FROM subscription_lists WHERE creator_id = ?"
	args := []interface{}{creatorID}

	// Check if start and end date are provided
	if startDate != "" && endDate != "" {
		query += " AND subscribed_at BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	result := r.DB.Raw(query, args...).Scan(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
