package bot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"birthday-bot/internal/config"
	"birthday-bot/internal/services"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Init(ctx context.Context, bot *telego.Bot, db *sql.DB, cfg *config.Config) {
	StartBirthdayChecker(ctx, bot, db)

	kb := MainKeyboard()

	waitAdd := make(map[int64]bool)
	waitDelete := make(map[int64]bool)
	waitTime := make(map[int64]bool)
	userDeleteMap := make(map[int64][]int)

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := strings.TrimSpace(update.Message.Text)
		userID := update.Message.Chat.ID

		// --- Проверка авторизации ---
		if !services.IsAuthorized(db, userID) && !strings.HasPrefix(text, "/access") {
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"🔒 Этот бот приватный.\nЧтобы получить доступ, введите команду:\n\n`/access <код>`").
				WithParseMode(telego.ModeMarkdown))
			continue
		}

		// --- Обработка команд ---
		switch {

		// --- Доступ по коду ---
		case strings.HasPrefix(text, "/access"):
			parts := strings.SplitN(text, " ", 2)
			if len(parts) != 2 {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"❌ Используйте формат: `/access <код>`").WithParseMode(telego.ModeMarkdown))
				continue
			}
			code := strings.TrimSpace(parts[1])
			if code != cfg.AccessCode {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID), "🚫 Неверный код доступа."))
				continue
			}

			services.AuthorizeUser(db, userID)
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"✅ Доступ разрешён! Добро пожаловать.").
				WithReplyMarkup(kb))

		// --- Список авторизованных пользователей (только для админа) ---
		case text == "/list_users":
			if userID != cfg.AdminID {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID), "🚫 У вас нет прав для этой команды."))
				continue
			}

			ids, err := services.GetAllUsers(db)
			if err != nil || len(ids) == 0 {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID), "📭 Нет авторизованных пользователей."))
				continue
			}

			var msg strings.Builder
			for i, id := range ids {
				chat, err := bot.GetChat(ctx, (&telego.GetChatParams{}).WithChatID(tu.ID(id)))
				name := fmt.Sprintf("ID: %d", id)
				if err == nil {
					full := chat.FirstName
					if chat.LastName != "" {
						full += " " + chat.LastName
					}
					if chat.Username != "" {
						full += fmt.Sprintf(" (@%s)", chat.Username)
					}
					name = full
				}
				msg.WriteString(fmt.Sprintf("%d. 👤 %s %d\n", i+1, name, id))
			}
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"📋 Список пользователей:\n\n"+msg.String()))

		// --- Отзыв доступа ---
		case strings.HasPrefix(text, "/revoke"):
			if userID != cfg.AdminID {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID), "🚫 Нет прав."))
				continue
			}
			parts := strings.SplitN(text, " ", 2)
			if len(parts) != 2 {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"Используйте формат: `/revoke <user_id>`").WithParseMode(telego.ModeMarkdown))
				continue
			}
			targetID, _ := strconv.ParseInt(parts[1], 10, 64)
			services.RevokeUser(db, targetID)
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				fmt.Sprintf("🗑 Пользователь %d удалён из списка доступа.", targetID)))

		// --- Старт ---
		case text == "/start":
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"👋 Привет! Я помогу тебе помнить дни рождения.\n"+
					"Используй кнопки ниже для управления.").
				WithReplyMarkup(kb))

		// --- Показать список ---
		case text == "📋 Показать список":
			rows, _ := db.Query(`SELECT name, date FROM birthdays WHERE user_id = ? ORDER BY date`, userID)

			var msg string
			var i int
			for rows.Next() {
				var name, date string
				rows.Scan(&name, &date)
				t, _ := time.Parse("2006-01-02", date)
				msg += fmt.Sprintf("%d. 🎂 %s — %s (%s)\n",
					i+1, name, t.Format("02.01.2006"), services.CalculateAge(t))
				i++
			}
			if msg == "" {
				msg = "📭 Пока нет сохранённых дней рождения."
			}
			bot.SendMessage(ctx, tu.Message(tu.ID(userID), msg).WithReplyMarkup(kb))

		// --- Добавить день рождения ---
		case text == "➕ Добавить день рождения":
			waitAdd[userID] = true
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"Введите имя и дату в формате:\n`Имя Фамилия ДД.ММ.ГГГГ`").
				WithParseMode(telego.ModeMarkdown))

		case waitAdd[userID]:
			parts := strings.Split(text, " ")
			if len(parts) != 3 {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"❌ Неверный формат. Пример: `Иван Петров 25.01.2006`").
					WithParseMode(telego.ModeMarkdown).
					WithReplyMarkup(kb))
				waitAdd[userID] = false
				continue
			}

			fName, sName, dateStr := parts[0], parts[1], parts[2]
			t, err := time.Parse("02.01.2006", dateStr)
			if err != nil {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"❌ Неверная дата. Пример: 25.01.2006").
					WithReplyMarkup(kb))
				waitAdd[userID] = false
				continue
			}

			name := fmt.Sprintf("%s %s", fName, sName)
			db.Exec(`INSERT INTO birthdays (user_id, name, date) VALUES (?, ?, ?)`, userID, name, t.Format("2006-01-02"))
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				fmt.Sprintf("✅ Добавлено: %s (%s)", name, t.Format("02.01.2006"))).
				WithReplyMarkup(kb))
			waitAdd[userID] = false

		// --- Удалить день рождения ---
		case text == "❌ Удалить день рождения":
			waitDelete[userID] = true
			rows, _ := db.Query(`SELECT id, name, date FROM birthdays WHERE user_id = ? ORDER BY date`, userID)

			var msg string
			var ids []int
			i := 1
			for rows.Next() {
				var id int
				var name, date string
				rows.Scan(&id, &name, &date)
				t, _ := time.Parse("2006-01-02", date)
				msg += fmt.Sprintf("%d. 🎂 %s — %s\n", i, name, t.Format("02.01.2006"))
				ids = append(ids, id)
				i++
			}
			rows.Close()

			if msg == "" {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID), "📭 Список пуст.").WithReplyMarkup(kb))
				waitDelete[userID] = false
				continue
			}

			userDeleteMap[userID] = ids
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"Введите номер для удаления:\n\n"+msg))

		case waitDelete[userID]:
			n, err := strconv.Atoi(text)
			if err != nil || n < 1 || n > len(userDeleteMap[userID]) {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"❌ Неверный номер.").WithReplyMarkup(kb))
				waitDelete[userID] = false
				continue
			}

			idToDelete := userDeleteMap[userID][n-1]
			db.Exec(`DELETE FROM birthdays WHERE id = ?`, idToDelete)
			bot.SendMessage(ctx, tu.Message(tu.ID(userID), "🗑 Удалено.").WithReplyMarkup(kb))
			waitDelete[userID] = false

		// --- Настроить время уведомлений ---
		case text == "⏰ Установить время напоминаний":
			waitTime[userID] = true
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"Введите время в формате HH:MM (например, 09:00 или 18:30):"))

		case waitTime[userID]:
			if _, err := time.Parse("15:04", text); err != nil {
				bot.SendMessage(ctx, tu.Message(tu.ID(userID),
					"❌ Неверный формат. Пример: 08:30").
					WithReplyMarkup(kb))
				waitTime[userID] = false
				continue
			}
			db.Exec(`INSERT INTO users (user_id, notify_time)
			         VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET notify_time = excluded.notify_time`,
				userID, text)
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				fmt.Sprintf("✅ Время напоминаний установлено на %s", text)).
				WithReplyMarkup(kb))
			waitTime[userID] = false

		default:
			bot.SendMessage(ctx, tu.Message(tu.ID(userID),
				"Выбери действие из меню ниже 👇").
				WithReplyMarkup(kb))
		}
	}
}
