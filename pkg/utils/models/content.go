package models

type VideoResponse struct {
	CategoryID  int    `form:"CategoryID" binding:"required"`
	Title       string `form:"Title" binding:"required"`
	Description string `form:"Description" binding:"required"`
}
type Video struct {
	ID          uint   `json:"id" gorm:"unique;not null"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Likes       int    `json:"likes"`
	Views       int    `json:"views"`
}
type VideoResponses struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	CategoryID  int    `json:"category_id"`
}
type RecommendationListResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type EditVideoDetails struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
