package services

import "database/sql"

func IsAuthorized(db *sql.DB, userID int64) bool {
	var id int64
	err := db.QueryRow(`SELECT user_id FROM authorized_users WHERE user_id = ?`, userID).Scan(&id)
	return err == nil
}

func AuthorizeUser(db *sql.DB, userID int64) {
	db.Exec(`INSERT OR IGNORE INTO authorized_users (user_id) VALUES (?)`, userID)
}

func RevokeUser(db *sql.DB, userID int64) {
	db.Exec(`DELETE FROM authorized_users WHERE user_id = ?`, userID)
}

func GetAllUsers(db *sql.DB) ([]int64, error) {
	rows, err := db.Query(`SELECT user_id FROM authorized_users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}
