package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("Нет DB_PATH")
	}

	// Создаём каталог, если не существует
	os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("❌ Ошибка при подключении к базе: %v", err)
	}
	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			notify_time TEXT DEFAULT '09:00'
		);`,
		`CREATE TABLE IF NOT EXISTS birthdays (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			date TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS authorized_users (
			user_id INTEGER PRIMARY KEY
		);`,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}
