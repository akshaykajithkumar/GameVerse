package interfaces

import "main/pkg/utils/models"

type SubscriptionUseCase interface {
	PurchasePlan(planID int, creatorID int, userID int) (string, error)
	MakePaymentRazorPay(planID string, userID int) (models.OrderPaymentDetails, error)
	VerifyPayment(paymentID string, razorID string, orderID string) error
	// StartSubscriptionUpdateJob()
}
