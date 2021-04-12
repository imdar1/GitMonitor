package services

import (
	"log"
	"time"
)

func GetFirstAndLastDayOfMonth() (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth, lastOfMonth
}

func CheckErr(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
}
