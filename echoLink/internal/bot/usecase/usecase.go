// internal/bot/usecase/bot_usecase.go
package usecase

import (
	"context"
	"echoLink/internal/bot/model"
	"echoLink/internal/bot/repository"
)

type BotUsecase interface {
	CreateBot(ctx context.Context, userID string, req CreateBotRequest) (*CreateBotResponse, error)
	ListBots(ctx context.Context, userID string) ([]BotSummary, error)
	GetBotForCall(ctx context.Context, botID, userID string) (*model.Bot, error)
	DeleteBot(ctx context.Context, botID string) error
	GetBotPublic(ctx context.Context, botID string) (*model.Bot, error)
}

type botUsecase struct {
	repo repository.BotRepository
}

func NewBotUsecase(repo repository.BotRepository) BotUsecase {
	return &botUsecase{repo: repo}
}

type CreateBotRequest struct {
	Goal    string
	Webhook string
	Context string
}

type CreateBotResponse struct {
	BotID string
}

type BotSummary struct {
	BotID     string `json:"bot_id"`
	Goal      string `json:"goal"`
	Webhook   string `json:"webhook"`
	CreatedAt string `json:"created_at"`
}

func (u *botUsecase) CreateBot(ctx context.Context, userID string, req CreateBotRequest) (*CreateBotResponse, error) {
	bot := &model.Bot{
		UserID:  userID,
		Goal:    req.Goal,
		Webhook: req.Webhook,
		Context: req.Context,
		Voice:   "man",
	}

	if err := u.repo.Create(ctx, bot); err != nil {
		return nil, err
	}

	return &CreateBotResponse{BotID: bot.ID}, nil
}

func (u *botUsecase) ListBots(ctx context.Context, userID string) ([]BotSummary, error) {
	bots, err := u.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]BotSummary, len(bots))
	for i, b := range bots {
		res[i] = BotSummary{
			BotID:     b.ID,
			Goal:      b.Goal,
			Webhook:   b.Webhook,
			CreatedAt: b.CreatedAt.Format("2006-01-02 15:04"),
		}
	}
	return res, nil
}

// NEW: DeleteBot implementation
func (u *botUsecase) DeleteBot(ctx context.Context, botID string) error {
	return u.repo.Delete(ctx, botID)
}

// GetBotPublic implements the fix for Twilio webhooks.
// It performs a simple lookup by ID only.
func (u *botUsecase) GetBotPublic(ctx context.Context, botID string) (*model.Bot, error) {
	// This calls the simple repository lookup, which ignores ownership.
	return u.repo.GetByID(ctx, botID)
}

// GetBotForCall is used by API endpoints (e.g., DELETE /bots/:id) and enforces ownership.
// It calls the repository method that checks both ID and UserID.
func (u *botUsecase) GetBotForCall(ctx context.Context, botID, userID string) (*model.Bot, error) {
	// This ensures only the owner can access the bot details via the protected API.
	return u.repo.GetByIDAndUserID(ctx, botID, userID)
}
