// internal/twilio/repository/call_repo.go
package repository

import (
	"context"
	"echoLink/internal/twilio/model"

	"github.com/uptrace/bun"
)

type CallRepository interface {
	Upsert(ctx context.Context, state *model.CallState) error
	Get(ctx context.Context, callSid string) (*model.CallState, error)
}

type callRepo struct {
	db *bun.DB
}

func NewCallRepository(db *bun.DB) CallRepository {
	return &callRepo{db: db}
}

func (r *callRepo) Upsert(ctx context.Context, state *model.CallState) error {
	_, err := r.db.NewInsert().
		Model(state).
		On("CONFLICT (call_sid) DO UPDATE").
		Set("bot_id = EXCLUDED.bot_id").
		Set("user_input = EXCLUDED.user_input").
		Set("assistant_resp = EXCLUDED.assistant_resp").
		Set("context = EXCLUDED.context").
		Set("step = EXCLUDED.step").
		Set("updated_at = NOW()").
		Exec(ctx)
	return err
}

func (r *callRepo) Get(ctx context.Context, callSid string) (*model.CallState, error) {
	var state model.CallState
	err := r.db.NewSelect().Model(&state).Where("call_sid = ?", callSid).Scan(ctx)
	return &state, err
}
