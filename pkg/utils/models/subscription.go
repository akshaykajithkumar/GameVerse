package models

type OrderPaymentDetails struct {
	UserID     int     `json:"user_id"`
	Username   string  `json:"username"`
	Razor_id   string  `josn:"razor_id"`
	OrderID    int     `json:"order_id"`
	FinalPrice float64 `json:"final_price"`
}
type AnalyticsData struct {
	SubscribersCount int     `json:"subscribers_count"`
	Revenue          float64 `json:"revenue"`
	// Add more fields as needed
}
