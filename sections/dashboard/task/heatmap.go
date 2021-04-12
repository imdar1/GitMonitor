package task

import (
	"gitmonitor/db"
	"gitmonitor/services"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type heatmap struct {
	colors    []string
	days      []string
	month     string
	task      []string
	taskCount int
}

func initData(db db.DBConfig) heatmap {
	firstDate, lastDate := services.GetFirstAndLastDayOfMonth()
	_, month, firstDay := firstDate.Date()
	_, _, lastDay := lastDate.Date()
	var days []string
	for i := firstDay; i <= lastDay; i++ {
		days = append(days, strconv.Itoa(i))
	}

	h := heatmap{
		days:  days,
		month: month.String(),
	}
	return h
}

func (h *heatmap) timelineBase() *charts.HeatMap {
	timeline := charts.NewHeatMap()
	timeline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Timeline for " + h.month,
		}),
	)
	return timeline
}
