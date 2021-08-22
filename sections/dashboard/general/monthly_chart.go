package general

import (
	"bytes"
	"fmt"
	"gitmonitor/services/utils"
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wcharczuk/go-chart/v2"
)

func fillWithZerosNTimes(arr []int, n int, sizeLimit int) []int {
	for i := 0; i < n; i++ {
		if len(arr) == sizeLimit {
			return arr
		}
		arr = append(arr, 0)
	}

	return arr
}

func toChartValueAndGetMax(elements []int, diff int) ([]chart.Value, float64) {
	var chartValue []chart.Value
	max := float64(0)
	currentTime := time.Now()
	for _, element := range elements {
		chartValue = append(chartValue, chart.Value{
			Style: chart.Style{
				Hidden:    false,
				ClassName: fmt.Sprint(element),
			},
			Value: float64(element),
			Label: currentTime.Format("2 Jan"),
		})
		max = math.Max(max, float64(element))
		currentTime = currentTime.AddDate(0, 0, -diff)
	}
	return chartValue, max
}

func getLast30DayCommits(commits []*object.Commit) []int {
	commitsCount := []int{}
	itr := 0 // commits current index
	currentDate := time.Now()

	for {
		beginningOfDay := utils.BeginningOfDay(currentDate)
		dayDiff := utils.GetDayDifference(utils.BeginningOfDay(commits[itr].Author.When), beginningOfDay)
		beginningOfDay = beginningOfDay.AddDate(0, 0, -dayDiff)

		// Fill with zeros until reached the latest commit date
		commitsCount = fillWithZerosNTimes(commitsCount, dayDiff, 30)
		if len(commitsCount) == 30 {
			return commitsCount
		}

		currCommitCount := 0
		for utils.BeginningOfDay(commits[itr].Author.When).Equal(beginningOfDay) {
			currentDate = commits[itr].Author.When.AddDate(0, 0, -1)
			currCommitCount++
			itr++
		}
		commitsCount = append(commitsCount, currCommitCount)
	}
}

func getMonthlyChart(commits []*object.Commit) image.Image {
	thisMonthCommits := getLast30DayCommits(commits)
	chartValue, maxVal := toChartValueAndGetMax(thisMonthCommits, 1)
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
		Height:       512,
		BarWidth:     20,
		UseBaseValue: true,
		BaseValue:    0,
		Bars:         chartValue,
		XAxis: chart.Style{
			StrokeWidth:         1,
			TextRotationDegrees: 90,
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
	utils.CheckErr("getMonthlyChart", err)

	return img
}

func getMonthlyChartCanvas(commits []*object.Commit) fyne.CanvasObject {
	chart := canvas.NewImageFromImage(getMonthlyChart(commits))
	chart.SetMinSize(fyne.NewSize(float32(chart.Image.Bounds().Dx()), float32(chart.Image.Bounds().Dy())))
	return container.NewMax(chart)
}
