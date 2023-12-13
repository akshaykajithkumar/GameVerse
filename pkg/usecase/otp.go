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

	ok := ot.otpRepository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}
	log.Printf("ACCOUNTSID=%s, AUTHTOKEN=%s", ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)

	_, err := helper.TwilioSendOTP(phone, ot.cfg.SERVICESID)
	if err != nil {
		log.Printf("TwilioSendOTP error: %v", err)
		return fmt.Errorf("error occurred while generating OTP: %v", err)
		//return errors.New("error ocurred while generating OTP")
	}

	return nil

}
