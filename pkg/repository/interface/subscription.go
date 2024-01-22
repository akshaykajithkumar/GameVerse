package interfaces

import "main/pkg/domain"

type SubscriptionRepository interface {
	GetSubscriptionListIDByPlanID(planID, creatorID, userID int) (int, error)
	FindUsername(user_id int) (string, error)
	FindPrice(orderID int) (float64, error)
	UpdatePaymentDetails(orderID, paymentID, razorID string) error
	GetActiveSubscription(creatorID, userID int) (*domain.SubscriptionList, error)

	// FetchActiveSubscriptions() ([]domain.SubscriptionList, error)
	// UpdateSubscriptionInDatabase(subscription domain.SubscriptionList) error
	// HasSubscriptionExpired(subscribedAt time.Time, duration int) bool
	// UpdateSubscriptionStatus() error
	// GetActiveSubscription(planID, creatorID, userID int) (*domain.SubscriptionList, error)
}
