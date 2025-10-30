// package http

// import (
// 	"fmt"
// 	"net/http"

// 	"echoLink/internal/user/usecase"

// 	"github.com/labstack/echo/v4"
// )

// type Handler struct {
// 	usecase usecase.UserUsecase
// }

// func NewHandler(uc usecase.UserUsecase) *Handler {
// 	return &Handler{usecase: uc}
// }

// // ConnectTwilioReq defines the structure for the initial registration and Twilio setup request.
// type ConnectTwilioReq struct {
// 	FirstName      string `json:"first_name"`
// 	LastName       string `json:"last_name"`
// 	Email          string `json:"email"`
// 	AccountSID     string `json:"account_sid"`
// 	AuthToken      string `json:"auth_token"`
// 	PhoneNumberSID string `json:"phone_number_sid"`
// }

// // ConnectTwilio handles the combined registration and Twilio connection. It issues the JWT.
// func (h *Handler) ConnectTwilio(c echo.Context) error {
// 	req := new(ConnectTwilioReq)
// 	if err := c.Bind(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
// 	}

// 	res, err := h.usecase.ConnectTwilio(c.Request().Context(), usecase.ConnectTwilioRequest{
// 		FirstName:      req.FirstName,
// 		LastName:       req.LastName,
// 		Email:          req.Email,
// 		AccountSID:     req.AccountSID,
// 		AuthToken:      req.AuthToken,
// 		PhoneNumberSID: req.PhoneNumberSID,
// 	})
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

// // NOTE: This must match the string defined by ClaimsContextKey in main.go
// const CustomClaimsContextKey = "userID"

// // GetUserIDFromContext safely extracts the userID string from the context.
// // It assumes the custom middleware has successfully run and stored the ID as a string.
// // This function is correct and robust against interface conversion panics.
// func GetUserIDFromContext(c echo.Context) (string, error) {
// 	// The custom middleware stores the final userID string under the "userID" key.
// 	claimsValue := c.Get(CustomClaimsContextKey)
// 	if claimsValue == nil {
// 		// This will be returned if the JWT middleware failed to run (e.g., token was bad)
// 		return "", fmt.Errorf("user ID not found in context (middleware missing or failed)")
// 	}

// 	// ASSERTION FIX: The middleware stores a string, so we assert to a string.
// 	userID, ok := claimsValue.(string)
// 	if !ok {
// 		// This should not happen if the middleware is correct, but acts as a final safeguard.
// 		return "", fmt.Errorf("invalid claims format: stored value is not a string")
// 	}

// 	return userID, nil
// }

// // GetMyNumber retrieves the authenticated user's phone number details.
// func (h *Handler) GetMyNumber(c echo.Context) error {
// 	userID, err := GetUserIDFromContext(c)
// 	if err != nil {
// 		// This now returns the specific error message from GetUserIDFromContext
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
// 	}

// 	res, err := h.usecase.GetMyNumber(c.Request().Context(), userID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

// // ListBots is a protected endpoint that requires authentication.
// func (h *Handler) ListBots(c echo.Context) error {
// 	userID, err := GetUserIDFromContext(c)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
// 	}

// 	// Placeholder response until ListBots usecase is implemented
// 	return c.JSON(http.StatusOK, []interface{}{
// 		map[string]string{"id": "bot_1", "goal": fmt.Sprintf("Placeholder for user %s: ListBots is pending", userID)},
// 	})
// }

package userhttp

import (
	"fmt"
	"net/http"

	"echoLink/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

// This key MUST match the key set in the custom AuthJWTMiddleware in main.go
const CustomClaimsContextKey = "userID"

type Handler struct {
	usecase usecase.UserUsecase
}

func NewHandler(usecase usecase.UserUsecase) *Handler {
	return &Handler{usecase: usecase}
}

// ConnectTwilioReq contains fields for registration/update and Twilio setup.
type ConnectTwilioReq struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	AccountSID     string `json:"account_sid" validate:"required"`
	AuthToken      string `json:"auth_token" validate:"required"`
	PhoneNumberSID string `json:"phone_number_sid" validate:"required"`
}

// ConnectTwilioRes includes the JWT for subsequent authentication.
type ConnectTwilioRes struct {
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
	AccessToken string `json:"access_token"`
}

// POST /v1/connect-twilio (Public Endpoint)
func (h *Handler) ConnectTwilio(c echo.Context) error {
	req := new(ConnectTwilioReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	res, err := h.usecase.ConnectTwilio(c.Request().Context(), usecase.ConnectTwilioRequest{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		AccountSID:     req.AccountSID,
		AuthToken:      req.AuthToken,
		PhoneNumberSID: req.PhoneNumberSID,
	})
	if err != nil {
		// Log the error for internal diagnostics
		fmt.Printf("ConnectTwilio failed: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ConnectTwilioRes{
		Message:     res.Message,
		PhoneNumber: res.PhoneNumber,
		AccessToken: res.AccessToken,
	})
}

// Helper to safely get the string UserID from context
func GetUserIDFromContext(c echo.Context) (string, error) {
	// CustomMiddlewareManager stores the ID under CustomClaimsContextKey as a string
	claimsValue := c.Get(CustomClaimsContextKey)

	if claimsValue == nil {
		return "", fmt.Errorf("user ID not found in context (middleware failed)")
	}

	userID, ok := claimsValue.(string)
	if !ok {
		return "", fmt.Errorf("user ID stored in context is not a string")
	}

	return userID, nil
}

// GET /v1/my-number (Protected Endpoint)
func (h *Handler) GetMyNumber(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
	}

	res, err := h.usecase.GetMyNumber(c.Request().Context(), userID)
	if err != nil {
		// Log the error (e.g., DB failure, or bot count failure)
		fmt.Printf("GetMyNumber usecase failed: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve data"})
	}

	return c.JSON(http.StatusOK, res)
}
