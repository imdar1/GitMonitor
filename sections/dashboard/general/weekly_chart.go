package general

import (
	"bytes"
	"fmt"
	"gitmonitor/services/utils"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wcharczuk/go-chart/v2"
)

func getCommitsCountByWeeks(commits []*object.Commit, weeks int) []int {
	commitsCount := []int{0}

	// get the beginning of current date
	now := time.Now()
	y, m, d := now.Date()
	beginningOfDay := time.Date(y, m, d, 0, 0, 0, 0, now.Location())

	// get the beginning of current week (date of the latest sunday)
	beginningOfWeek := beginningOfDay.AddDate(0, 0, int(beginningOfDay.Weekday())*-1)
	for _, v := range commits {
		if weeks == 0 {
			break
		}

		if v.Author.When.After(beginningOfWeek) {
			commitsCount[len(commitsCount)-1] += 1
		} else {
			beginningOfWeek = beginningOfWeek.AddDate(0, 0, -7)
			weeks -= 1
			commitsCount = append(commitsCount, 1)
		}
	}

	commitsCount = utils.Reverse(commitsCount)
	return commitsCount
}

func toChartValue(elements []int, placeholder string) []chart.Value {
	var chartValue []chart.Value
	for index, element := range elements {
		chartValue = append(chartValue, chart.Value{
			Value: float64(element),
			Label: fmt.Sprintf(placeholder, index+1),
		})
	}
	return chartValue
}

func getWeeklyChart(commits []*object.Commit) image.Image {
	last10WeeksCommitsCount := getCommitsCountByWeeks(commits, 10)
	chartValue := toChartValue(last10WeeksCommitsCount, "Week-%d")

	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     chartValue,
		YAxis: chart.YAxis{
			ValueFormatter: chart.IntValueFormatter,
		},
	}

	buf := new(bytes.Buffer)
	graph.Render(chart.PNG, buf)
	img, _, err := image.Decode(buf)
	utils.CheckErr(err)

	return img
}

func getWeeklyChartCanvas(commits []*object.Commit) fyne.CanvasObject {
	chart := canvas.NewImageFromImage(getWeeklyChart(commits))
	chart.SetMinSize(fyne.NewSize(float32(chart.Image.Bounds().Dx()), float32(chart.Image.Bounds().Dy())))
	return container.NewMax(chart)
}
