package main

import (
	"context"
	"log"

	"birthday-bot/internal/bot"
	"birthday-bot/internal/config"
	"birthday-bot/internal/db"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	database := db.InitDB()

	log.Println("🤖 Бот в работе...")
	bot.Start(ctx, cfg, database)
}
