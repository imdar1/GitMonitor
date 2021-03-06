package utils

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func GetFirstAndLastDayOfMonth() (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth, lastOfMonth
}

func GetStringFromDatetime(t time.Time) string {
	return t.Format("20060102")
}

func GetDayDifference(start time.Time, end time.Time) int {
	return int(end.Sub(start).Hours()/24) + 1
}

func CheckErr(serviceName string, err error) {
	if err != nil {
		log.Printf("[%s] Error: %s", serviceName, err.Error())
	}
}

func CreateBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}

func BeginningOfMonth() time.Time {
	now := time.Now()
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

func GetBeginningOfMonthByTime(inputTime time.Time) time.Time {
	y, m, _ := inputTime.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, inputTime.Location())
}

func BeginningOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

// Find element from a given array of string. Condition: list is already sorted
func IsExistStr(element string, list []string) bool {
	for _, v := range list {
		if element == v {
			return true
		}
	}
	return false
}

func IsExistInt(element int, list []int) bool {
	for _, v := range list {
		if element == v {
			return true
		}
	}

	return false
}

func Reverse(input []int) []int {
	a := input
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}
