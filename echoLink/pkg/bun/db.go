// pkg/bun/db.go
package bun

import (
	"context"
	"database/sql"
	"fmt"

	"echoLink/config"
	botmodel "echoLink/internal/bot/model"
	twiliomodel "echoLink/internal/twilio/model"
	usermodel "echoLink/internal/user/model"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func InitDB(cfg *config.Config) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Postgres.URL)))
	DB = bun.NewDB(sqldb, pgdialect.New())

	ctx := context.Background()

	// collect all models here
	models := []any{
		(*usermodel.User)(nil),
		(*botmodel.Bot)(nil),
		(*twiliomodel.CallState)(nil), // example – adjust if exists
	}

	// loop + auto-create tables if not exist
	for _, m := range models {
		if _, err := DB.NewCreateTable().
			Model(m).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx); err != nil {
			fmt.Printf("❌ failed creating table for %T: %v\n", m, err)
		} else {
			fmt.Printf("✅ ensured table for %T exists\n", m)
		}
	}

	// optional: create trigger function for updated_at
	if _, err := DB.ExecContext(ctx, `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`); err != nil {
		fmt.Println("❌ failed creating trigger func:", err)
	}

	// attach triggers to tables that have updated_at columns
	tables := []string{"users", "bots"} // add more if needed
	for _, t := range tables {
		sql := fmt.Sprintf(`
			DROP TRIGGER IF EXISTS update_%[1]s_updated_at ON %[1]s;
			CREATE TRIGGER update_%[1]s_updated_at
			BEFORE UPDATE ON %[1]s
			FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
		`, t)
		if _, err := DB.ExecContext(ctx, sql); err != nil {
			fmt.Printf("❌ failed creating trigger for %s: %v\n", t, err)
		}
	}
}
