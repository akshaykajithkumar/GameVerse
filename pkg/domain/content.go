package domain

import (
	"fmt"
	"strings"
	"time"
)

type Category struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Category string `json:"category" gorm:"unique;not null"`
}

//	type Video struct {
//		ID          uint     `json:"id" gorm:"unique;not null"`
//		UserID      uint     `json:"user_id" gorm:"not null"`
//		Title       string   `json:"title"`
//		Description string   `json:"description"`
//		URL         string   `json:"url"`
//		CategoryID  int      `json:"category_id"`
//		Category    Category `json:"category" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
//		Likes       int      `json:"likes" gorm:"default:0"`
//		Views       int      `json:"views" gorm:"default:0"`
//	}

type Tags []string

// Custom scanner for the Tags type
// Custom scanner for the Tags type
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if str, ok := value.(string); ok {
		*t = strings.Split(str[1:len(str)-1], ",")
		return nil
	}
	return fmt.Errorf("failed to scan tags")
}

// Video struct with a custom scanner for the tags column

type Video struct {
	ID          uint      `json:"id" gorm:"unique;not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CategoryID  int       `json:"category_id"`
	Category    Category  `json:"category" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	Likes       int       `json:"likes" gorm:"default:0"`
	Views       int       `json:"views" gorm:"default:0"`
	Exclusive   bool      `json:"exclusive" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
}

type VideoLikes struct {
	ID      uint `json:"id" gorm:"primaryKey"`
	UserID  uint `json:"user_id"`
	VideoID uint `json:"video_id"`
}

// Comment represents a comment on a video
type Comment struct {
	ID      uint   `json:"id" gorm:"unique;not null"`
	UserID  uint   `json:"user_id" gorm:"not null"`
	VideoID uint   `json:"video_id" gorm:"not null"`
	Content string `json:"content"`
}

// Tag represents a tags.
type Tag struct {
	ID  uint   `json:"id" gorm:"primaryKey"`
	Tag string `json:"tag" gorm:"not null"`
}
type VideoTags struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  int    `json:"user_id"`
	VideoID uint   `json:"video_id" gorm:"not null"`
	Tag     string `json:"tag" gorm:"not null"`
	// Tag     Tag  `json:"-" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

type SubscriptionList struct {
	ID               int              `json:"id" gorm:"primaryKey"`
	CreatorID        int              `json:"creator_id"`
	UserID           int              `json:"user_id"`
	User             User             `json:"-" gorm:"foreignKey:UserID"`
	SubscribedAt     time.Time        `json:"subscribed_at" gorm:"-"`
	PlanID           int              `json:"plan_id"`
	PaymentStatus    string           `json:"paymentStatus" gorm:"default:'Pending'"`
	SubscriptionPlan SubscriptionPlan `json:"-" gorm:"foreignKey:PlanID"`
	PaymentID        string           `json:"paymentID"`
	IsActive         bool             `json:"isActive" gorm:"default:false"`
}
