package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"birthday-bot/internal/services"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// StartBirthdayChecker запускает фоновый процесс проверки дней рождения
func StartBirthdayChecker(ctx context.Context, bot *telego.Bot, db *sql.DB) {
	go func() {
		for {
			loc, _ := time.LoadLocation("Europe/Moscow") // можно поменять на свой регион
			now := time.Now().In(loc)

			rows, err := db.Query(`SELECT user_id, notify_time FROM users`)
			if err != nil {
				log.Println("DB error:", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			for rows.Next() {
				var userID int64
				var notifyTime string
				rows.Scan(&userID, &notifyTime)

				if now.Format("15:04") == notifyTime {
					sendBirthdayNotifications(ctx, bot, db, userID)
				}
			}
			rows.Close()

			time.Sleep(60 * time.Second)
		}
	}()
}

// sendBirthdayNotifications проверяет у конкретного пользователя дни рождения на сегодня
func sendBirthdayNotifications(ctx context.Context, bot *telego.Bot, db *sql.DB, userID int64) {
	today := time.Now().Format("01-02") // формат MM-DD для сравнения
	rows, err := db.Query(`SELECT name, date FROM birthdays WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("DB error:", err)
		return
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name, dateStr string
		rows.Scan(&name, &dateStr)

		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if t.Format("01-02") == today {
			names = append(names, fmt.Sprintf("%s — %s (%s)", name, t.Format("02.01.2006"), services.CalculateAge(t)))
		}
	}

	if len(names) == 0 {
		return
	}

	msg := "🎉 Сегодня день рождения у:\n"
	for _, n := range names {
		msg += "🎂 " + n + "\n"
	}

	bot.SendMessage(ctx, tu.Message(tu.ID(userID), msg))
}
