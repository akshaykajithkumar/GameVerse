package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
	"mime/multipart"
)

type UserUseCase interface {
	SignUp(user models.UserDetails) (models.TokenUser, error)
	ChangePassword(id int, old string, password string, repassword string) error
	Login(user models.UserLogin) (models.TokenUser, error)
	EditProfile(id int, name, email, username, phone, bio string, image *multipart.FileHeader) error
	GetProfile(id int) (*models.UserProfileResponse, error)
	ReportUser(reporterID, targetID int, reason string) error
	ToggleFollow(followerID, followingID int) error
	GetFollowingListWithPagination(userID int, page, limit int) ([]models.FollowingUser, error)
	SearchUsersByNameWithPagination(searchTerm string, page, limit int) ([]domain.User, error)
	GetSubscriptionPlans() ([]domain.SubscriptionPlan, error)
	GetFollowersListWithPagination(userID int, page, limit int) ([]models.FollowerUser, error)
}
