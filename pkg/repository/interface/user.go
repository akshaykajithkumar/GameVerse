package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
)

type UserRepository interface {
	SignUp(user models.UserDetails) (models.UserResponse, error)
	CheckUserAvailability(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserResponse, error)
	ChangePassword(id int, password string) error
	GetPassword(id int) (string, error)
	UserBlockStatus(email string) (bool, error)
	EditProfile(id int, name, email, username, phone, bio string, url string) error
	GetProfileDetailsById(id int) (*domain.User, error)
	StoreReport(reporterID, targetID int, reason string) error
	CheckFollowRelationship(followerID, followingID int) (bool, error)
	StoreFollow(followerID, followingID int) error
	RemoveFollow(followerID, followingID int) error
	SearchUsersByNameWithPagination(searchTerm string, page, limit int) ([]domain.User, error)
	GetFollowingListWithPagination(userID int, page, limit int) ([]models.FollowingUser, error)
	GetSubscriptionPlans() ([]domain.SubscriptionPlan, error)
	GetFollowersListWithPagination(userID int, page, limit int) ([]models.FollowerUser, error)
}
