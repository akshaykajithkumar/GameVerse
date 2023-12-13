package usecase

import (
	"errors"
	"main/pkg/helper"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (u *userUseCase) SignUp(user models.UserDetails) (models.TokenUser, error) {
	// Check whether the user already exist. If yes, show the error message, since this is signUp
	userExist := u.userRepo.CheckUserAvailability(user.Email)
	if userExist {
		return models.TokenUser{}, errors.New("user already exist, sign in")
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUser{}, errors.New("password does not match")
	}

	// Hash password since details are validated

	hashedPassword, err := helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUser{}, err
	}

	user.Password = hashedPassword

	// add user details to the database
	userData, err := u.userRepo.SignUp(user)
	if err != nil {
		return models.TokenUser{}, err
	}

	// crete a JWT token string for the user
	accessTokenString, refreshTokenString, err := helper.GenerateTokensUser(userData)
	if err != nil {
		return models.TokenUser{}, errors.New("could not create token due to some internal error")
	}

	return models.TokenUser{
		Username: user.Username,

		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	}, nil
}

func (u *userUseCase) Login(user models.UserLogin) (models.TokenUser, error) {
	// checking if a username exist with this email address
	ok := u.userRepo.CheckUserAvailability(user.Email)
	if !ok {
		return models.TokenUser{}, errors.New("the user does not exist")
	}

	permission, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUser{}, err
	}

	if !permission {
		return models.TokenUser{}, errors.New("user is blocked by admin")
	}

	// Get the user details in order to check the password, in this case ( The same function can be reused in future )
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUser{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user_details.Password), []byte(user.Password))
	if err != nil {
		return models.TokenUser{}, errors.New("password incorrect")
	}

	accessToken, refreshToken, err := helper.GenerateTokensUser(user_details)
	if err != nil {
		return models.TokenUser{}, errors.New("could not create token")
	}
	return models.TokenUser{
		Username:     user_details.Username,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil

}
func (i *userUseCase) ChangePassword(id int, old string, password string, repassword string) error {

	userPassword, err := i.userRepo.GetPassword(id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(old))
	if err != nil {
		return errors.New("password incorrect")
	}

	if password != repassword {
		return errors.New("passwords does not match")
	}

	newpassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("internal server error")
	}

	return i.userRepo.ChangePassword(id, string(newpassword))

}

const maxBioLength = 60

func (u *userUseCase) AddBio(id int, bio string) error {
	// Check if bio length exceeds the limit
	if len(bio) > maxBioLength {
		return errors.New("bio length exceeds the limit")
	}
	if err := u.userRepo.AddBio(id, bio); err != nil {
		return err
	}

	return nil
}
func (u *userUseCase) EditProfile(id int, name, email, username, phone, bio string) error {
	// Validate the bio field length
	if len(bio) > maxBioLength {
		return errors.New("bio length exceeds the limit")
	}

	// Call the repository to update the user profile
	if err := u.userRepo.EditProfile(id, name, email, username, phone, bio); err != nil {
		return err
	}

	return nil
}

// GetProfile retrieves the user profile details by user ID
func (u *userUseCase) GetProfile(id int) (*models.UserProfileResponse, error) {
	// Call the repository to fetch the user profile
	user, err := u.userRepo.GetProfileDetailsById(id)
	if err != nil {
		return nil, err
	}

	// Convert the user details to UserProfileResponse
	userProfile := &models.UserProfileResponse{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
		Bio:      user.Bio,
	}

	return userProfile, nil
}

const maxReasonLength = 60

func (u *userUseCase) ReportUser(reporterID, targetID int, reason string) error {
	// Additional validation if needed

	// Check if the reason length exceeds the limit
	if len(reason) > maxReasonLength {
		return errors.New("reason : length exceeds the limit")
	}

	// Call the user repository to store the report
	if err := u.userRepo.StoreReport(reporterID, targetID, reason); err != nil {
		return err
	}

	return nil
}
