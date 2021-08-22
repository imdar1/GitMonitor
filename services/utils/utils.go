package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/wcharczuk/go-chart/v2"
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
	return int(math.Floor(end.Sub(start).Hours() / 24))
}

// Convert string date with "DD/MM/YYYY" format into unix timestamp
func GetUnixTimeStampFromString(timeString string) (int64, error) {
	date, err := time.Parse("02/01/2006", timeString)
	if err != nil {
		return 0, err
	}
	return date.Unix(), nil
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

func calculateEffectiveBarSpacing(bc chart.BarChart, canvasBox chart.Box) int {
	totalWithBaseSpacing := len(bc.Bars) * (bc.GetBarWidth() + bc.GetBarSpacing())
	if totalWithBaseSpacing > canvasBox.Width() {
		lessBarWidths := canvasBox.Width() - (len(bc.Bars) * bc.GetBarWidth())
		if lessBarWidths > 0 {
			return int(math.Ceil(float64(lessBarWidths) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarSpacing()
}

func calculateEffectiveBarWidth(bc chart.BarChart, canvasBox chart.Box, spacing int) int {
	totalWithBaseWidth := len(bc.Bars) * (bc.GetBarWidth() + spacing)
	if totalWithBaseWidth > canvasBox.Width() {
		totalLessBarSpacings := canvasBox.Width() - (len(bc.Bars) * spacing)
		if totalLessBarSpacings > 0 {
			return int(math.Ceil(float64(totalLessBarSpacings) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarWidth()
}

func AddLabel(c *chart.BarChart, chartRange chart.Range) chart.Renderable {
	return func(r chart.Renderer, canvasBox chart.Box, chartDefaults chart.Style) {
		xoffset := canvasBox.Left

		spacing := calculateEffectiveBarSpacing(*c, canvasBox)
		width := calculateEffectiveBarWidth(*c, canvasBox, spacing)
		bs2 := spacing >> 1

		var bxl, bxr, by int
		for _, bar := range c.Bars {
			bxl = xoffset + bs2
			bxr = bxl + width
			xoffset += width + spacing
			if int(bar.Value) == 0 {
				continue
			}

			by = canvasBox.Bottom - chartRange.Translate(bar.Value)
			lx := bxl + ((bxr - bxl) / 2)
			ly := (canvasBox.Bottom + by) / 2

			legendDefaults := chart.Style{
				FillColor:   chart.ColorWhite,
				FontColor:   chart.DefaultTextColor,
				FontSize:    8.0,
				StrokeColor: chart.DefaultAxisColor,
				StrokeWidth: chart.DefaultAxisLineWidth,
			}
			bar.Style.InheritFrom(chartDefaults.InheritFrom(legendDefaults)).WriteToRenderer(r)
			tb := r.MeasureText(fmt.Sprintf("%.0f", bar.Value))
			lx = lx - (tb.Width() >> 1)
			ly = ly + (tb.Height() >> 1)

			if lx < 0 {
				lx = 0
			}
			if ly < 0 {
				lx = 0
			}

			r.Text(fmt.Sprintf("%.0f", bar.Value), lx, ly)
		}
	}
}
