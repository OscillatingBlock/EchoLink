// internal/bot/model/bot.go
package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Bot struct {
	bun.BaseModel `bun:"table:bots"`

	ID        string    `bun:",pk,default:gen_random_uuid()"`
	UserID    string    `bun:",notnull"`
	Goal      string    `bun:",notnull"`
	Webhook   string    `bun:",notnull,url"`
	Context   string    `bun:",nullzero"`
	Voice     string    `bun:",notnull,default:'man'"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

// Prompt returns full LLM context
func (b *Bot) Prompt() string {
	if b.Context == "" {
		return b.Goal
	}
	return b.Goal + "\n\nContext: " + b.Context
}
