package db

type Birthday struct {
	ID     int
	UserID int64
	Name   string
	Date   string
}

type User struct {
	UserID     int64
	NotifyTime string
}
