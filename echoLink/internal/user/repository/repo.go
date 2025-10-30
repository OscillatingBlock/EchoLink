package repository

import (
	"context"
	"echoLink/internal/user/model"
	"fmt"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByPhoneSID(ctx context.Context, sid string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error) // NEW
}

type userRepo struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Upsert(ctx context.Context, user *model.User) error {
	// Upsert is now based on Email for identification, and ID as PK
	// We use the ID if present, or rely on the unique Email constraint for conflict resolution.
	_, err := r.db.NewInsert().
		Model(user).
		On("CONFLICT (id) DO UPDATE"). // Using ID for conflict resolution
		Set("first_name = EXCLUDED.first_name").
		Set("last_name = EXCLUDED.last_name").
		Set("email = EXCLUDED.email"). // Email is unique and is updated here
		Set("client_id = EXCLUDED.client_id").
		Set("client_secret = EXCLUDED.client_secret").
		Set("twilio_sid = EXCLUDED.twilio_sid").
		Set("twilio_token = EXCLUDED.twilio_token").
		Set("phone_number = EXCLUDED.phone_number").
		Set("phone_number_sid = EXCLUDED.phone_number_sid").
		Set("updated_at = NOW()").
		Exec(ctx)

	// Note: If ID is not provided on insert, bun should use gen_random_uuid().
	// If you prefer to rely on the unique Email constraint for the upsert logic:
	// On("CONFLICT (email) DO UPDATE")...
	// However, keeping existing ID-based logic and letting the usecase handle the registration vs update check is cleaner for now.

	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	return &user, err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	return &user, err
}

func (r *userRepo) GetByPhoneSID(ctx context.Context, sid string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("phone_number_sid = ?", sid).Scan(ctx)
	return &user, err
}

