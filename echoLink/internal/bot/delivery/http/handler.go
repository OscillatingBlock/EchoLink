package bothttp

import (
	"fmt"
	"net/http"

	"echoLink/internal/bot/usecase"
	"github.com/labstack/echo/v4"
)

// This key MUST match the key set in the custom AuthJWTMiddleware in main.go
const CustomClaimsContextKey = "userID"

type Handler struct {
	usecase usecase.BotUsecase
}

func NewHandler(usecase usecase.BotUsecase) *Handler {
	return &Handler{usecase: usecase}
}

type CreateBotReq struct {
	Goal    string `json:"goal" validate:"required"`
	Webhook string `json:"webhook" validate:"required,url"`
	Context string `json:"context"`
}

type CreateBotRes struct {
	BotID string `json:"bot_id"`
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

// POST /v1/bots
func (h *Handler) CreateBot(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
	}
	req := new(CreateBotReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	res, err := h.usecase.CreateBot(c.Request().Context(), userID, usecase.CreateBotRequest{
		Goal:    req.Goal,
		Webhook: req.Webhook,
		Context: req.Context,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create bot"})
	}

	return c.JSON(http.StatusCreated, CreateBotRes{BotID: res.BotID})
}

// GET /v1/bots
func (h *Handler) ListBots(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
	}

	bots, err := h.usecase.ListBots(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, bots)
}

// GET /v1/bots/:id
func (h *Handler) GetBot(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
	}

	botID := c.Param("id")
	if botID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bot ID is required"})
	}

	// Use GetBotForCall since it checks both ID and UserID (security principle)
	bot, err := h.usecase.GetBotForCall(c.Request().Context(), botID, userID)
	if err != nil {
		// Assuming usecase returns an error if bot is not found or user doesn't own it
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bot not found"})
	}

	return c.JSON(http.StatusOK, bot)
}

// DELETE /v1/bots/:id
func (h *Handler) DeleteBot(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
	}

	botID := c.Param("id")
	if botID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bot ID is required"})
	}

	// First, check ownership by attempting to retrieve it
	_, err = h.usecase.GetBotForCall(c.Request().Context(), botID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bot not found or access denied"})
	}

	// If ownership is confirmed, proceed with deletion
	if err := h.usecase.DeleteBot(c.Request().Context(), botID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete bot"})
	}

	return c.NoContent(http.StatusNoContent)
}
