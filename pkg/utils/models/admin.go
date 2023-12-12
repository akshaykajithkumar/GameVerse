package models

type AdminLogin struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}
type TokenAdmin struct {
	Username     string
	RefreshToken string
	AccessToken  string
}
type UserDetailsAtAdmin struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Permission bool   `json:"permission"`
}
type AdminDetailsResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email" `
}
