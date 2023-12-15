package models

type VerifyData struct {
	PhoneNumber string `json:"phone,omitempty" validate:"required"`
	Code        string `json:"code,omitempty" validate:"required"`
}

type OTPData struct {
	PhoneNumber string `json:"phone,omitempty" validate:"required"`
}
type ForgotPasswordData struct {
	PhoneNumber     string `json:"phone,omitempty" validate:"required"`
	Code            string `json:"code,omitempty" validate:"required"`
	Newpassword     string `json:"password,omitempty" validate:"min=8,max=20"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"min=8,max=20"`
}
