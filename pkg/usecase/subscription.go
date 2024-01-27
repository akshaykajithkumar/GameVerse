package usecase

import (
	"fmt"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"strconv"

	"github.com/razorpay/razorpay-go"
)

type subscriptionUseCase struct {
	repository interfaces.SubscriptionRepository
}

func NewSubscriptionUseCase(repo interfaces.SubscriptionRepository) services.SubscriptionUseCase {
	return &subscriptionUseCase{
		repository: repo,
	}
}

// // PurchasePlan handles the purchase of a subscription plan
// func (i *subscriptionUseCase) PurchasePlan(planID int,creatorID) (string, error) {
//     // Get the subscription list ID directly from the repository
//     subscriptionListID, err := i.subscriptionRepository.GetSubscriptionListIDByPlanID(planID,creatorID)
//     if err != nil {
//         // Handle the error, plan not found
//         return "", err
//     }

//     link := fmt.Sprintf("http://localhost:1245/users/payment/razorpay?subscription_list_id=%d", subscriptionListID)

//	    return link, nil
//	}
//
// PurchasePlan handles the purchase of a subscription plan
func (i *subscriptionUseCase) PurchasePlan(planID int, creatorID int, userID int) (string, error) {
	// Check if the user already has an active subscription
	existingSubscription, err := i.repository.GetActiveSubscription(creatorID, userID)
	if err != nil {
		// Handle the error, e.g., plan not found or database error
		return "", err
	}

	// If the user is already subscribed, return a message
	if existingSubscription != nil {
		return "User already subscribed", nil
	}

	// Proceed with the subscription purchase process

	// Get the subscription list ID directly from the repository
	subscriptionListID, err := i.repository.GetSubscriptionListIDByPlanID(planID, creatorID, userID)
	if err != nil {
		// Handle the error, e.g., plan not found or database error
		return "", err
	}

	link := fmt.Sprintf("http://localhost:1245/users/plans/choose-plan/razorpay?subscription_list_id=%d", subscriptionListID)

	return link, nil
}

func (p *subscriptionUseCase) MakePaymentRazorPay(planID string, userID int) (models.OrderPaymentDetails, error) {
	var orderDetails models.OrderPaymentDetails

	newid, err := strconv.Atoi(planID)
	if err != nil {
		return models.OrderPaymentDetails{}, err
	}
	orderDetails.OrderID = newid

	orderDetails.UserID = userID

	//get username
	username, err := p.repository.FindUsername(userID)
	if err != nil {
		return models.OrderPaymentDetails{}, err
	}

	orderDetails.Username = username

	//get total
	newfinal, err := p.repository.FindPrice(newid)
	if err != nil {
		return models.OrderPaymentDetails{}, err
	}

	orderDetails.FinalPrice = newfinal

	client := razorpay.NewClient("rzp_test_hfgxceEmSupD3T", "9UrqCcAM6F5JMXqTFmJBQBGQ")

	data := map[string]interface{}{
		"amount":   int(orderDetails.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return models.OrderPaymentDetails{}, nil
	}

	razorPayOrderID := body["id"].(string)

	orderDetails.Razor_id = razorPayOrderID

	return orderDetails, nil
}
func (p *subscriptionUseCase) VerifyPayment(paymentID string, razorID string, orderID string) error {

	err := p.repository.UpdatePaymentDetails(orderID, paymentID, razorID)
	if err != nil {
		return err
	}

	return nil

}

// func (p *subscriptionUseCase) StartSubscriptionUpdateJob() {
// 	// Schedule the job using AddFunc
// 	c := cron.New()
// 	c.AddFunc("0 0 * * *", func() {
// 		// Inside the scheduled function, call UpdateSubscriptionStatus
// 		if err := p.repository.UpdateSubscriptionStatus(); err != nil {
// 			log.Println("Error updating subscription statuses:", err)
// 		}
// 	})

// 	// Start the scheduler
// 	c.Start()

//		// Block the main goroutine to keep the application running
//		select {}
//	}

func (u *subscriptionUseCase) GetAnalytics(userID int, startDate string, endDate string) (models.AnalyticsData, error) {
	//Fetch subscribers count for the given date range
	subscribersCount, err := u.repository.GetSubscribersCount(userID, startDate, endDate)
	if err != nil {
		return models.AnalyticsData{}, err
	}

	// Fetch revenue for the given date range
	revenue, err := u.repository.GetRevenue(userID, startDate, endDate)
	if err != nil {
		return models.AnalyticsData{}, err
	}

	// Create and return the analytics data
	analyticsData := models.AnalyticsData{
		SubscribersCount: subscribersCount,
		Revenue:          revenue,
		// Add more fields if needed
	}

	return analyticsData, nil
}
