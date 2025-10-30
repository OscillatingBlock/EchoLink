// internal/twilio/model/call_state.go
package model

import (
	"github.com/uptrace/bun"
	"time"
)

// CallState tracks conversation state per call (Twilio CallSid)
type CallState struct {
	bun.BaseModel `bun:"table:call_states"`

	CallSid       string    `bun:",pk"` // Twilio CallSid
	BotID         string    `bun:",notnull"`
	UserInput     string    `bun:",nullzero"`
	AssistantResp string    `bun:",nullzero"`
	Context       string    `bun:",nullzero"`
	Step          int       `bun:",default:0"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
