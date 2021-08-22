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

func getCommitsCountByWeeks(commits []*object.Commit, weeks int) []int {
	commitsCount := []int{}

	// get the beginning of current date
	now := time.Now()
	beginningOfDay := utils.BeginningOfDay(now)

	// get the beginning of the current week
	beginningOfWeek := beginningOfDay.AddDate(0, 0, -7)
	for _, v := range commits {
		if weeks <= 0 {
			break
		}

		if v.Author.When.After(beginningOfWeek) {
			if len(commitsCount) == 0 {
				commitsCount = append(commitsCount, 1)
			} else {
				commitsCount[len(commitsCount)-1] += 1
			}
		} else {
			isInsert := false
			for v.Author.When.Before(beginningOfWeek) && weeks > 0 {
				if isInsert {
					commitsCount = append(commitsCount, 0)
					weeks -= 1
				}
				beginningOfWeek = beginningOfWeek.AddDate(0, 0, -7)
				isInsert = true
			}
			if weeks > 0 {
				commitsCount = append(commitsCount, 1)
			}
			weeks -= 1
		}
	}
	commitsCount = fillWithZerosNTimes(commitsCount, 10, 10)
	return commitsCount
}

func getWeeklyChart(commits []*object.Commit) image.Image {
	last10WeeksCommitsCount := getCommitsCountByWeeks(commits, 9)
	chartValue, maxVal := toChartValueAndGetMax(last10WeeksCommitsCount, 7, 7)
	// prevent runtime error whenever each element is 0
	maxVal += 1
	chartRange := &chart.ContinuousRange{
		Min: 0,
		Max: maxVal,
	}

	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     chartValue,
		XAxis: chart.Style{
			StrokeWidth: 1,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				StrokeWidth: 1,
			},
			Range: chartRange,
		},
	}

	graph.Elements = []chart.Renderable{
		utils.AddLabel(&graph, chartRange),
	}

	buf := new(bytes.Buffer)
	graph.Render(chart.PNG, buf)
	img, _, err := image.Decode(buf)
	utils.CheckErr("getWeeklyChart", err)

	return img
}

func getWeeklyChartCanvas(commits []*object.Commit) fyne.CanvasObject {
	chart := canvas.NewImageFromImage(getWeeklyChart(commits))
	chart.SetMinSize(fyne.NewSize(float32(chart.Image.Bounds().Dx()), float32(chart.Image.Bounds().Dy())))
	return container.NewMax(chart)
}
