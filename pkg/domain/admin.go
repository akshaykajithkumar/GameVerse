package domain

// Admin represents an administrative user in the system.
type Admin struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Username string `json:"name" gorm:"validate:required"`
	Email    string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}
type SubscriptionPlan struct {
	ID       int     `gorm:"primaryKey" json:"id"`
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	Price    float64 `json:"price"`
}
