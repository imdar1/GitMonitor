package utils

import (
	"log"
	"sort"
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

func CheckErr(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
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

// Find element from a given array of string. Condition: list is already sorted
func IsExist(element string, list []string) bool {
	index := sort.Search(
		len(list),
		func(i int) bool {
			// get list with the value >= element
			return list[i] >= element
		},
	)
	return index < len(list) && list[index] == element
}
