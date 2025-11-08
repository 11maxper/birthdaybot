package bot

import (
	"context"
	"database/sql"
	"log"

	"birthday-bot/internal/config"

	"github.com/mymmrac/telego"
)

func Start(ctx context.Context, cfg *config.Config, db *sql.DB) {
	b, err := telego.NewBot(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	Init(ctx, b, db, cfg)
}
