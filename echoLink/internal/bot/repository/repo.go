// internal/bot/repository/bot_repo.go
package repository

import (
	"context"
	"echoLink/internal/bot/model"

	"github.com/uptrace/bun"
)

type BotRepository interface {
	Create(ctx context.Context, bot *model.Bot) error
	ListByUserID(ctx context.Context, userID string) ([]model.Bot, error)
	GetByID(ctx context.Context, id string) (*model.Bot, error)
	GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bot, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, id string) error
}

type botRepo struct {
	db *bun.DB
}

func NewBotRepository(db *bun.DB) BotRepository {
	return &botRepo{db: db}
}

func (r *botRepo) Create(ctx context.Context, bot *model.Bot) error {
	_, err := r.db.NewInsert().Model(bot).Exec(ctx)
	return err
}

func (r *botRepo) ListByUserID(ctx context.Context, userID string) ([]model.Bot, error) {
	var bots []model.Bot
	err := r.db.NewSelect().
		Model(&bots).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx)
	return bots, err
}

func (r *botRepo) GetByID(ctx context.Context, id string) (*model.Bot, error) {
	var bot model.Bot
	err := r.db.NewSelect().Model(&bot).Where("id = ?", id).Scan(ctx)
	return &bot, err
}

func (r *botRepo) GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bot, error) {
	var bot model.Bot
	err := r.db.NewSelect().
		Model(&bot).
		Where("id = ? AND user_id = ?", id, userID).
		Scan(ctx)
	return &bot, err
}

func (r *botRepo) CountByUserID(ctx context.Context, userID string) (int, error) {
	count, err := r.db.NewSelect().
		Model((*model.Bot)(nil)). // Select on the table type without fetching data
		Where("user_id = ?", userID).
		Count(ctx)
	return count, err
}

// NEW: Delete implementation
func (r *botRepo) Delete(ctx context.Context, id string) error {
	// Simple hard delete for now. In production, consider soft-delete (setting a status field).
	_, err := r.db.NewDelete().
		Model((*model.Bot)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

