package helper

import (
	"context"
	"errors"
	"fmt"
	"main/pkg/domain"
	"main/pkg/utils/models"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

// AuthCustomClaims represents custom claims for JWT
type AuthCustomClaims struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateTokensAdmin(admin domain.Admin) (string, string, error) {
	accessTokenClaims := &AuthCustomClaims{
		Id:    uint(admin.ID),
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		Id:    uint(admin.ID),
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	fmt.Println("Admin tokens created")
	return accessTokenString, refreshTokenString, nil
}

func GenerateTokensUser(user models.UserResponse) (string, string, error) {
	accessTokenClaims := &AuthCustomClaims{
		Id:    uint(user.Id),
		Email: user.Email,
		Role:  "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		Id:    uint(user.Id),
		Email: user.Email,
		Role:  "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	fmt.Println("User tokens created")
	return accessTokenString, refreshTokenString, nil
}

/*
validateToken is for decrypting a jwt token using HMAC256 algorithm

Parameters:
- token: JWT token string.
*/
func ValidateToken(token string) (*jwt.Token, error) {
	fmt.Println("Token validating.........")
	jwttoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		secret := viper.GetString("KEY")
		return []byte(secret), nil
	})

	return jwttoken, err
}

// using for generating tokens when access token expires

func TokensFromRefreshToken(prevRefreshTokenString string) (string, string, error) {
	// Parse the previous refresh token
	prevRefreshToken, err := jwt.Parse(prevRefreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("KEY")), nil
	})

	if err != nil {
		return "", "", err
	}

	// Extract claims from the previous refresh token
	prevRefreshClaims, ok := prevRefreshToken.Claims.(jwt.MapClaims)
	if !ok || !prevRefreshToken.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	// Use the claims to generate a new access token
	newAccessTokenClaims := &AuthCustomClaims{
		Id:    uint(prevRefreshClaims["id"].(float64)),
		Email: prevRefreshClaims["email"].(string),
		Role:  prevRefreshClaims["role"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims)
	newAccessTokenString, err := newAccessToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	// Generate a new refresh token for the next cycle
	newRefreshTokenClaims := &AuthCustomClaims{
		Id:    uint(prevRefreshClaims["id"].(float64)),
		Email: prevRefreshClaims["email"].(string),
		Role:  prevRefreshClaims["role"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims)
	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(viper.GetString("KEY")))
	if err != nil {
		return "", "", err
	}

	return newAccessTokenString, newRefreshTokenString, nil
}

/*
GetUserID returns the userID stored in the context

Parameters:
- c: gin context

Returns:
- int: userID
- error: error is returned
*/
func GetUserID(c *gin.Context) (int, error) {
	var key models.UserKey = "userID"
	val := c.Request.Context().Value(key)

	// Check if the value is not nil
	if val == nil {
		return 0, errors.New("userID not found in context")
	}

	// Use type assertion to convert to the expected type
	userKey, ok := val.(models.UserKey)
	if !ok {
		return 0, errors.New("failed to convert userID to the expected type")
	}

	ID := userKey.String()
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return 0, errors.New("failed to convert userID to int")
	}

	return userID, nil
}

/*
PasswordHashing hashes a password.

Parameters:
- password: Password to be hashed.

Returns:
- string: Hashed Password.
- error: Error is returned if any.
*/
func PasswordHashing(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}

	hash := string(hashedPassword)
	return hash, nil
}

var client *twilio.RestClient

/*
TwilioSetup will setup the twillio.

Parameters:
- username: Twillio Username.
- password: Twillio Password.
*/
func TwilioSetup(username string, password string) {
	// log.Printf("username=%s,password=%s", username, password)
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

}

/*
TwilioSendOTP sends otp to the number provides from the specified service

Parameters:
- phone: Otp reciever phone number.
- serviceID: Twillio Service ID to choose the service.

Returns:
- string: The unique string that we created to identify the Verification resource.
- error: Error is returned if any.
*/
func TwilioSendOTP(phone string, serviceID string) (string, error) {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {

		return " ", err
	}

	return *resp.Sid, nil

}

/*
TwilioVerifyOTP verifies the otp sent to the number

Parameters:
- phone: Otp reciever phone number.
- serviceID: Twillio Service ID to choose the service.
- code: OTP.

Returns:
- error: Error is returned if any.
*/
func TwilioVerifyOTP(serviceID string, code string, phone string) error {

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + phone)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)

	if err != nil {
		return err
	}

	if *resp.Status == "approved" {
		return nil
	}

	return errors.New("failed to validate otp")

}

func AddImageToS3(file *multipart.FileHeader) (string, error) {
	// Set AWS credentials using environment variables
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAX2D5JXBMLEOAGJOW")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		fmt.Println("configuration error:", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	f, openErr := file.Open()
	if openErr != nil {
		fmt.Println("opening error:", openErr)
		return "", openErr
	}
	defer f.Close()

	result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("bucketforgameverse"),
		Key:    aws.String(file.Filename),
		Body:   f,
		//ACL:    "public-read",
	})

	if uploadErr != nil {
		fmt.Println("uploading error:", uploadErr)
		return "", uploadErr
	}

	return result.Location, nil
}
