package models

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenUser struct {
	Username     string
	RefreshToken string
	AccessToken  string
}
type UserResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// signup
type UserDetails struct {
	Name            string `json:"name"`
	Email           string `json:"email" validate:"email"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
}

// user details shown after logging in
type UserDetailsResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UserProfile struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
}

type UserSettings struct {
	ID                      int  `json:"id"`
	UserID                  int  `json:"user_id"`
	NotificationPreferences bool `json:"notification_preferences"`
}

type EditUserProfile struct {
	ID           int `json:"id"`
	UserID       int `json:"user_id"`
	User         User
	UserProfile  UserProfile
	UserSettings UserSettings
}

type UserKey string

func (k UserKey) String() string {
	return string(k)
}
