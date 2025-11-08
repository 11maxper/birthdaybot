package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Token      string
	AdminID    int64
	AccessCode string
}

func Load() *Config {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("❌ TOKEN не найден в .env")
	}

	adminID, err := strconv.ParseInt(os.Getenv("ADMINID"), 10, 64)
	if err != nil {
		log.Fatal("❌ ADMINID должен быть числом")
	}

	code := os.Getenv("ACCESSCODE")
	if code == "" {
		log.Fatal("❌ ACCESSCODE не найден")
	}

	return &Config{
		Token:      token,
		AdminID:    adminID,
		AccessCode: code,
	}
}
