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
}

type UserProfile struct {
	ID             int    `json:"id"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
	UserID         uint   `json:"user_id"`
	User           User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
