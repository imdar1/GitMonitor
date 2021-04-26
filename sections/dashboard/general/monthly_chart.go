package general

import (
	"bytes"
	"gitmonitor/services/utils"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wcharczuk/go-chart/v2"
)

func getThisMonthCommits(commits []*object.Commit) []int {
	commitsCount := []int{0}

	firstDayOfMonth := utils.BeginningOfMonth()
	lastDayChecked := time.Now().Day()
	for _, v := range commits {
		if v.Author.When.Before(firstDayOfMonth) {
			break
		}

		if v.Author.When.Day() == lastDayChecked {
			commitsCount[len(commitsCount)-1] += 1
		} else if v.Author.When.Day() < lastDayChecked {
			for i := lastDayChecked; i > v.Author.When.Day(); i-- {
				commitsCount = append(commitsCount, 0)
			}
			commitsCount[len(commitsCount)-1] += 1
		}
		lastDayChecked = v.Author.When.Day()
	}

	commitsCount = utils.Reverse(commitsCount)
	return commitsCount
}

func getMonthlyChart(commits []*object.Commit) image.Image {
	thisMonthCommits := getThisMonthCommits(commits)
	chartValue, maxVal := toChartValueAndGetMax(thisMonthCommits, "Day-%d")

	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:       512,
		BarWidth:     20,
		UseBaseValue: true,
		BaseValue:    0,
		Bars:         chartValue,
		XAxis: chart.Style{
			TextRotationDegrees: 90,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: maxVal,
			},
		},
	}

	buf := new(bytes.Buffer)
	graph.Render(chart.PNG, buf)
	img, _, err := image.Decode(buf)
	utils.CheckErr(err)

	return img
}

func getMonthlyChartCanvas(commits []*object.Commit) fyne.CanvasObject {
	chart := canvas.NewImageFromImage(getMonthlyChart(commits))
	chart.SetMinSize(fyne.NewSize(float32(chart.Image.Bounds().Dx()), float32(chart.Image.Bounds().Dy())))
	return container.NewMax(chart)
}
