package bot

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// MainKeyboard возвращает основное меню клавиатуры
func MainKeyboard() *telego.ReplyKeyboardMarkup {
	return tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("➕ Добавить день рождения"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("📋 Показать список"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("⏰ Установить время напоминаний"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("❌ Удалить день рождения"),
		),
	).WithResizeKeyboard()
}
