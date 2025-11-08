package services

import (
	"database/sql"
	"time"
)

func GetTodaysBirthdays(db *sql.DB, userID int64) ([]string, error) {
	now := time.Now().Format("01-02")
	rows, err := db.Query(`SELECT name, date FROM birthdays WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name, date string
		rows.Scan(&name, &date)
		t, err := time.Parse("2006-01-02", date)
		if err == nil && t.Format("01-02") == now {
			names = append(names, name)
		}
	}
	return names, nil
}
