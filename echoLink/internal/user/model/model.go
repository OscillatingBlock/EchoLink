// internal/user/model/user.go
package model

import (
	"time"

	"github.com/uptrace/bun"
)

// User owns Twilio number and bots
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID string `bun:",pk,default:gen_random_uuid()"`

	// NEW REGISTRATION FIELDS
	FirstName string `bun:"first_name,notnull"`
	LastName  string `bun:"last_name,notnull"`
	Email     string `bun:"email,notnull,unique"` // Used for identification/upsert

	// NEW INTERNAL AUTH FIELDS (Long-lived key pair)
	// Note: ClientSecret MUST be hashed before saving to DB
	ClientID     string `bun:"client_id,notnull,unique"`
	ClientSecret string `bun:"client_secret,notnull"`

	// EXISTING TWILIO FIELDS
	TwilioSID      string `bun:"twilio_sid,notnull"`
	TwilioToken    string `bun:"twilio_token,notnull"`
	PhoneNumber    string `bun:"phone_number,notnull,unique"`
	PhoneNumberSID string `bun:"phone_number_sid,notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
