package usecase

import (
	"context"
	"echoLink/internal/user/model"
	"echoLink/internal/user/repository"
	"fmt"
	"time"

	botrepository "echoLink/internal/bot/repository" // IMPORTED
	"github.com/golang-jwt/jwt/v5"
	"github.com/twilio/twilio-go"
	twilioapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type UserUsecase interface {
	ConnectTwilio(ctx context.Context, req ConnectTwilioRequest) (*ConnectTwilioResponse, error)
	GetMyNumber(ctx context.Context, userID string) (*GetMyNumberResponse, error)
}

type userUsecase struct {
	repo      repository.UserRepository
	botRepo   botrepository.BotRepository // ADDED: Bot repository for counting
	jwtSecret []byte
}

// Updated Constructor to accept BotRepository
func NewUserUsecase(repo repository.UserRepository, botRepo botrepository.BotRepository, secret []byte) UserUsecase {
	return &userUsecase{
		repo:      repo,
		botRepo:   botRepo, // Initialize bot repo
		jwtSecret: secret,
	}
}

// Assuming the following structs are defined here or imported:
type ConnectTwilioRequest struct {
	FirstName      string
	LastName       string
	Email          string
	AccountSID     string
	AuthToken      string
	PhoneNumberSID string
}

type ConnectTwilioResponse struct {
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
	AccessToken string `json:"access_token"`
}

type GetMyNumberResponse struct {
	PhoneNumber string `json:"phone_number"`
	BotsCount   int    `json:"bots_count"`
}

// ConnectTwilio (JWT Generation Logic)
func (u *userUsecase) ConnectTwilio(ctx context.Context, req ConnectTwilioRequest) (*ConnectTwilioResponse, error) {
	// ... (Twilio validation logic remains the same) ...

	// 1. Create Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: req.AccountSID,
		Password: req.AuthToken,
	})

	// 2. Fetch phone number from Twilio
	params := &twilioapi.FetchIncomingPhoneNumberParams{}
	phone, err := client.Api.FetchIncomingPhoneNumber(req.PhoneNumberSID, params)
	if err != nil {
		return nil, fmt.Errorf("invalid Twilio credentials or number: %w", err)
	}
	if phone.PhoneNumber == nil {
		return nil, fmt.Errorf("phone number not found")
	}

	// 3. Save/Update User
	user, _ := u.repo.GetByEmail(ctx, req.Email)
	if user == nil {
		// Mock ID/Secret generation for new user
		user = &model.User{ID: fmt.Sprintf("user_%d", time.Now().UnixNano()), ClientID: "mock", ClientSecret: "mock"}
	}
	user.TwilioSID = req.AccountSID
	user.TwilioToken = req.AuthToken
	user.PhoneNumber = *phone.PhoneNumber
	user.PhoneNumberSID = req.PhoneNumberSID
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	if err := u.repo.Upsert(ctx, user); err != nil {
		return nil, err
	}

	// 4. Generate JWT
	claims := jwt.MapClaims{"userID": user.ID, "exp": time.Now().Add(time.Hour * 24).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &ConnectTwilioResponse{
		Message:     "Connected",
		PhoneNumber: *phone.PhoneNumber,
		AccessToken: signedToken,
	}, nil
}

// FIX APPLIED HERE: GetMyNumber now calls the bot repository to get the count.
func (u *userUsecase) GetMyNumber(ctx context.Context, userID string) (*GetMyNumberResponse, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fetch bot count using the injected repository
	botsCount, err := u.botRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bot count: %w", err)
	}

	return &GetMyNumberResponse{
		PhoneNumber: user.PhoneNumber,
		BotsCount:   botsCount, // Uses the actual count
	}, nil
}

