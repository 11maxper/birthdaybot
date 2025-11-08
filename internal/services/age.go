package services

import (
	"fmt"
	"time"
)

func CalculateAge(birthDate time.Time) string {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	var suffix string
	lastDigit := age % 10
	lastTwo := age % 100

	if lastTwo >= 11 && lastTwo <= 14 {
		suffix = "лет"
	} else {
		switch lastDigit {
		case 1:
			suffix = "год"
		case 2, 3, 4:
			suffix = "года"
		default:
			suffix = "лет"
		}
	}
	return fmt.Sprintf("%d %s", age, suffix)
}
