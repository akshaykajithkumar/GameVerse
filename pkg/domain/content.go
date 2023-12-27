package domain

type Category struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Category string `json:"category" gorm:"unique;not null"`
}

type Video struct {
	ID          uint     `json:"id" gorm:"unique;not null"`
	UserID      uint     `json:"user_id" gorm:"not null"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	CategoryID  int      `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
}
