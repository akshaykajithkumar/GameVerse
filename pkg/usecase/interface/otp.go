package interfaces

import "main/pkg/utils/models"

type OtpUseCase interface {
	VerifyOTP(code models.VerifyData) (models.TokenUser, error)
	SendOTP(phone string) error
}
