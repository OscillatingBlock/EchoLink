package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"echoLink/config"
	bot_http "echoLink/internal/bot/delivery/http" // Alias for clarity
	botrepo "echoLink/internal/bot/repository"
	"echoLink/internal/bot/usecase"
	twiliohttp "echoLink/internal/twilio/delivery/http"
	twiliorepo "echoLink/internal/twilio/repository"
	twiliousecase "echoLink/internal/twilio/usecase"
	userhttp "echoLink/internal/user/delivery/http"
	userrepo "echoLink/internal/user/repository"
	userusecase "echoLink/internal/user/usecase"
	bun "echoLink/pkg/bun"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// --- Placeholder for Middleware Management Structures ---

// CustomClaimsKey is used to store claims in the request context
type CustomClaimsKey string

const ClaimsContextKey CustomClaimsKey = "userID"

// CustomMiddlewareManager holds dependencies needed by middleware functions
type CustomMiddlewareManager struct {
	Secret string
	// In a real app, this would also hold a Logger and AuthUsecase
}

func NewCustomMiddlewareManager(secret string) *CustomMiddlewareManager {
	return &CustomMiddlewareManager{Secret: secret}
}

// --- Custom JWT Middleware Functions ---

func (mw *CustomMiddlewareManager) AuthJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerHeader := c.Request().Header.Get("Authorization")
		if bearerHeader == "" {
			// If no header is present, the request is unauthorized for protected routes
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Missing token"})
		}

		headerParts := strings.Split(bearerHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			// Check for "Bearer <token>" format
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Invalid token format"})
		}

		tokenString := headerParts[1]
		userID, err := mw.ValidateJWTToken(c.Request().Context(), tokenString)
		if err != nil {
			log.Printf("JWT Validation Failed: %v", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Invalid or expired token"})
		}

		// Inject the UserID into the Echo Context
		c.Set(string(ClaimsContextKey), userID)

		return next(c)
	}
}

func (mw *CustomMiddlewareManager) ValidateJWTToken(ctx context.Context, tokenString string) (string, error) {
	if tokenString == "" {
		return "", jwt.ErrInvalidKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := []byte(mw.Secret)
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrInvalidType
	}

	// 1. Check Expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return "", jwt.ErrTokenExpired
		}
	} else {
		// Tokens MUST have an expiration claim in this flow
		return "", fmt.Errorf("missing exp claim")
	}

	// 2. Extract User ID
	// Use fmt.Sprintf to robustly convert the claim value to a string
	userID := fmt.Sprintf("%v", claims["user_id"])
	if userID == "" {
		return "", fmt.Errorf("missing user_id claim")
	}

	return userID, nil
}

// --- Main Application Entry Point ---

func main() {
	cfg := config.Load()
	bun.InitDB(cfg)

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	// Repositories
	userRepo := userrepo.NewUserRepository(bun.DB)
	botRepo := botrepo.NewBotRepository(bun.DB)
	callRepo := twiliorepo.NewCallRepository(bun.DB)

	// Usecases (Dependency Injection)
	// Inject botRepo into userUC for bot counting
	userUC := userusecase.NewUserUsecase(userRepo, botRepo, []byte(cfg.JWT.Secret))
	botUC := usecase.NewBotUsecase(botRepo)
	twilioUC := twiliousecase.NewTwilioUsecase(botUC, callRepo, cfg.MLService.URL)

	// Initialize Custom Middleware Manager
	mw := NewCustomMiddlewareManager(cfg.JWT.Secret)

	// --- PUBLIC ENDPOINTS (No Auth Required) ---

	e.POST("/v1/voice", twiliohttp.NewHandler(twilioUC).HandleVoice)
	e.POST("/v1/voice-response", twiliohttp.NewHandler(twilioUC).HandleVoiceResponse)
	e.POST("/v1/connect-twilio", userhttp.NewHandler(userUC).ConnectTwilio)

	// --- PROTECTED API GROUP (Custom JWT Auth Required) ---
	api := e.Group("/v1")
	// Apply the custom AuthJWTMiddleware
	api.Use(mw.AuthJWTMiddleware)

	botHandler := bot_http.NewHandler(botUC)

	// User Endpoints
	api.GET("/my-number", userhttp.NewHandler(userUC).GetMyNumber)

	// Bot Collection Endpoints
	api.POST("/bots", botHandler.CreateBot)
	api.GET("/bots", botHandler.ListBots)

	// Bot Resource Endpoints (Fixes the 404)
	api.GET("/bots/:id", botHandler.GetBot)       // NEW: Get a single bot
	api.DELETE("/bots/:id", botHandler.DeleteBot) // NEW: Delete a single bot

	log.Println("Server on :" + cfg.Server.Port)
	e.Start(":" + cfg.Server.Port)
}
