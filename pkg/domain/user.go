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
	URL        string `json:"url"`
}

type Reports struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	ReporterID int    `json:"reporter_id"`
	TargetID   int    `json:"target_id"`
	Reason     string `json:"reason"`
}
type UserTags struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID int  `json:"user_id"`
	TagID  uint `json:"tag_id" gorm:"not null"`
	Tag    Tag  `json:"-" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// Follow struct represents a user following another user
type Follow struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	FollowerID  int  `json:"follower_id"`
	FollowingID int  `json:"following_id"`
}

// You can include additional fields if needed, such as timestamps.
