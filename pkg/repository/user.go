package repository

import (
	"errors"
	"fmt"
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}
func (c *userDatabase) CheckUserAvailability(email string) bool {

	var count int
	query := fmt.Sprintf("select count(*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0
}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	var permission bool
	err := cr.DB.Raw("select permission from users where email = ?", email).Scan(&permission).Error
	if err != nil {
		return false, err
	}
	return permission, nil
}
func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserResponse, error) {

	var user_details models.UserResponse

	err := c.DB.Raw(`
		SELECT *
		FROM users where email = ? and permission = true
		`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserResponse{}, errors.New("error checking user details")
	}

	return user_details, nil
}
func (c *userDatabase) SignUp(user models.UserDetails) (models.UserResponse, error) {

	var userDetails models.UserResponse
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone, username) VALUES (?, ?, ?, ?,?) RETURNING id, name, email, phone", user.Name, user.Email, user.Password, user.Phone, user.Username).Scan(&userDetails).Error

	if err != nil {
		return models.UserResponse{}, err
	}

	return userDetails, nil
}
func (i *userDatabase) ChangePassword(id int, password string) error {

	err := i.DB.Exec("UPDATE users SET password=? WHERE id=?", password, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (i *userDatabase) GetPassword(id int) (string, error) {

	var userPassword string
	err := i.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}

// EditProfile updates the user profile fields in the repository
func (i *userDatabase) EditProfile(id int, name, email, username, phone, bio, url string) error {
	// Create a map to store the fields that need to be updated
	updateFields := make(map[string]interface{})

	// Check if each field is provided and update the map accordingly
	if name != "" {
		updateFields["name"] = name
	}

	if email != "" {
		updateFields["email"] = email
	}

	if username != "" {
		updateFields["username"] = username
	}

	if phone != "" {
		updateFields["phone"] = phone
	}

	if bio != "" {
		updateFields["bio"] = bio
	}
	if url != "" {
		updateFields["url"] = url
	}

	// Update the user's profile fields in the database
	err := i.DB.Model(&domain.User{}).Where("id = ?", id).Updates(updateFields).Error
	if err != nil {
		return err
	}

	return nil
}

// retrieves specific user details by ID from the database
func (i *userDatabase) GetProfileDetailsById(id int) (*domain.User, error) {
	user := &domain.User{}
	err := i.DB.Select("name, email, username, phone, bio,url").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// for storing the reports from userReports
func (i *userDatabase) StoreReport(reporterID, targetID int, reason string) error {
	report := domain.Reports{
		ReporterID: reporterID,
		TargetID:   targetID,
		Reason:     reason,
	}

	// Create a new report record in the database
	if err := i.DB.Create(&report).Error; err != nil {
		return err
	}

	return nil
}

// CheckFollowRelationship checks if a follow relationship exists between two users
func (i *userDatabase) CheckFollowRelationship(followerID, followingID int) (bool, error) {
	var count int64

	// Use raw SQL query to count the follow relationships
	err := i.DB.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ?", followerID, followingID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// StoreFollow creates a new follow relationship in the database
func (i *userDatabase) StoreFollow(followerID, followingID int) error {
	// Use raw SQL query to insert a new follow record
	err := i.DB.Exec("INSERT INTO follows (follower_id, following_id) VALUES (?, ?)", followerID, followingID).Error
	if err != nil {
		return err
	}

	return nil
}

// RemoveFollow removes a follow relationship from the database
func (i *userDatabase) RemoveFollow(followerID, followingID int) error {
	// Use raw SQL query to delete a follow record
	err := i.DB.Exec("DELETE FROM follows WHERE follower_id = ? AND following_id = ?", followerID, followingID).Error
	if err != nil {
		return err
	}

	return nil
}

// GetFollowingListWithPagination retrieves the paginated list of users (ID and username) that a given user is following
func (ur *userDatabase) GetFollowingListWithPagination(userID int, page, limit int) ([]models.FollowingUser, error) {
	var followingList []models.FollowingUser

	// Use a raw SQL query to retrieve the following list with pagination
	err := ur.DB.Raw("SELECT u.id, u.username FROM users u JOIN follows f ON u.id = f.following_id WHERE f.follower_id = ? ORDER BY u.username LIMIT ? OFFSET ?", userID, limit, (page-1)*limit).Scan(&followingList).Error
	if err != nil {
		return nil, err
	}

	return followingList, nil
}

// SearchUsersByNameWithPagination returns a paginated list of users matching the search term in alphabetical order
func (ur *userDatabase) SearchUsersByNameWithPagination(searchTerm string, page, limit int) ([]domain.User, error) {
	var searchResults []domain.User

	// Use a raw SQL query to retrieve users matching the search term in alphabetical order with pagination
	err := ur.DB.Raw("SELECT id, name, username FROM users WHERE LOWER(name) LIKE LOWER(?) ORDER BY name LIMIT ? OFFSET ?", "%"+searchTerm+"%", limit, (page-1)*limit).Scan(&searchResults).Error
	if err != nil {
		return nil, err
	}

	return searchResults, nil
}
func (ar *userDatabase) GetSubscriptionPlans() ([]domain.SubscriptionPlan, error) {
	var plans []domain.SubscriptionPlan

	// Fetch the list of subscription plans from the database
	err := ar.DB.Find(&plans).Error
	if err != nil {
		// Handle any error during the fetch operation, you can log or perform additional actions as needed
		return nil, err
	}

	// Return the list of subscription plans
	return plans, nil
}
func (ur *userDatabase) GetFollowersListWithPagination(userID int, page, limit int) ([]models.FollowerUser, error) {
	var followersList []models.FollowerUser

	// Use a raw SQL query to retrieve the followers list with pagination
	err := ur.DB.Raw("SELECT u.id, u.username FROM users u JOIN follows f ON u.id = f.follower_id WHERE f.following_id = ? ORDER BY u.username LIMIT ? OFFSET ?", userID, limit, (page-1)*limit).Scan(&followersList).Error
	if err != nil {
		return nil, err
	}

	return followersList, nil
}
