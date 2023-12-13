package domain

// User represents a user in the system.

type User struct {
	ID         int    `gorm:"primaryKey"`
	Name       string `json:"name"`
	Email      string `gorm:"unique" json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Phone      string `gorm:"unique" json:"phone"`
	Permission bool   `gorm:"default:false" json:"permission"`
	Bio        string `json:"bio"`
}

type Reports struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	ReporterID int    `json:"reporter_id"`
	TargetID   int    `json:"target_id"`
	Reason     string `json:"reason"`
}
