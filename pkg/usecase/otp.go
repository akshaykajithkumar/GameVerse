package usecase

import (
	"errors"
	"fmt"
	"log"
	"main/pkg/config"
	"main/pkg/helper"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"

	"golang.org/x/crypto/bcrypt"
)

type otpUseCase struct {
	cfg           config.Config
	otpRepository interfaces.OtpRepository
}

func NewOtpUseCase(cfg config.Config, repo interfaces.OtpRepository) services.OtpUseCase {
	return &otpUseCase{
		cfg:           cfg,
		otpRepository: repo,
	}
}
func (ot *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUser, error) {

	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(ot.cfg.SERVICESID, code.Code, code.PhoneNumber)
	if err != nil {
		//this guard clause catches the error code runs only until here
		return models.TokenUser{}, errors.New("error while verifying")
	}

	if err := ot.otpRepository.UpdateUserPermissionByPhone(code.PhoneNumber); err != nil {
		return models.TokenUser{}, err
	}
	// if user is authenticated using OTP send back user details
	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(code.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, err
	}

	accessTokenString, refreshTokenString, err := helper.GenerateTokensUser(userDetails)
	if err != nil {
		return models.TokenUser{}, err

	}
	return models.TokenUser{
		Username:     userDetails.Username,
		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	}, nil
}

func (ot *otpUseCase) SendOTP(phone string) error {

	log.Printf("ACCOUNTSID=%s, AUTHTOKEN=%s", ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	log.Printf("SERVICESID=%s", ot.cfg.SERVICESID)
	_, err := helper.TwilioSendOTP(phone, ot.cfg.SERVICESID)

	if err != nil {
		log.Printf("TwilioSendOTP error: %v", err)
		return errors.New("error occurred while generating OTP")
		//return errors.New("error ocurred while generating OTP")
	}

	return nil

}
func (ot *otpUseCase) ForgotPassword(details models.ForgotPasswordData) (models.TokenUser, error) {
	ok := ot.otpRepository.FindUserByMobileNumber(details.PhoneNumber)
	if !ok {
		return models.TokenUser{}, errors.New("the user does not exist")
	}
	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(ot.cfg.SERVICESID, details.Code, details.PhoneNumber)
	if err != nil {
		// This guard clause catches the error and runs only until here
		return models.TokenUser{}, errors.New("error while verifying")
	}

	if details.ConfirmPassword != details.Newpassword {
		return models.TokenUser{}, errors.New("passwords do not match")
	}

	// Generate a new password using bcrypt
	newPassword, err := bcrypt.GenerateFromPassword([]byte(details.ConfirmPassword), 10)
	if err != nil {
		return models.TokenUser{}, errors.New("internal server error")
	}

	// Change the password in the repository
	if err := ot.otpRepository.ChangePasswordByPhone(details.PhoneNumber, string(newPassword)); err != nil {
		return models.TokenUser{}, fmt.Errorf("failed to change password: %v", err)
	}

	// If the user is authenticated using OTP, send back user details
	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(details.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, err
	}

	// Generate tokens for the user
	accessTokenString, refreshTokenString, err := helper.GenerateTokensUser(userDetails)
	if err != nil {
		return models.TokenUser{}, err
	}

	return models.TokenUser{
		Username:     userDetails.Username,
		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	}, nil
}
