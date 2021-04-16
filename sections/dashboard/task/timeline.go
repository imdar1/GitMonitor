package task

import (
	"bytes"
	"gitmonitor/services"
	"time"

	"github.com/gregoryv/draw/design"
	"github.com/gregoryv/draw/types/date"
)

type taskInformation struct {
	startDateStr string
	days         int
}

type timelineData struct {
	startDateStr string
	days         int
	tasks        map[string]taskInformation
}

func initData(taskData TaskData) timelineData {
	var data timelineData

	if len(taskData.Tasks) > 0 {
		taskCount := len(taskData.Tasks)
		startDate := time.Unix(taskData.Tasks[0].StartDate, 0)
		endDate := time.Unix(taskData.Tasks[taskCount-1].EndDate, 0)

		data = timelineData{
			startDateStr: services.GetStringFromDatetime(startDate),
			days:         services.GetDayDifference(startDate, endDate),
		}
		for _, v := range taskData.Tasks {
			startTask := time.Unix(v.StartDate, 0)
			endTask := time.Unix(v.EndDate, 0)
			taskInfo := taskInformation{
				startDateStr: services.GetStringFromDatetime(startTask),
				days:         services.GetDayDifference(startTask, endTask),
			}
			data.tasks[v.Name] = taskInfo
		}
	}
	return data
}

func initDummy() timelineData {
	tasks := map[string]taskInformation{
		"abc": {
			startDateStr: "20210401",
			days:         10,
		},
		"def": {
			startDateStr: "20210411",
			days:         5,
		},
		"ghi": {
			startDateStr: "20210413",
			days:         8,
		},
	}
	dummy := timelineData{
		startDateStr: "20210401",
		days:         30,
		tasks:        tasks,
	}
	return dummy
}

func (t *timelineData) getGanttChartImage() []byte {
	if len(t.tasks) == 0 {
		return []byte{}
	}

	var bar *design.Task
	ganttChart := design.NewGanttChart(date.String(t.startDateStr), t.days)
	for key, value := range t.tasks {
		bar = ganttChart.Add(key)
		ganttChart.Place(bar).At(date.String(value.startDateStr), value.days)
	}
	ganttChart.SetCaption("Project Schedule")
	imgBuffer := new(bytes.Buffer)
	styling := ganttChart.Diagram.Style
	styling.SetOutput(imgBuffer)
	ganttChart.WriteSVG(&styling)

	imgByte := services.GetImage(imgBuffer.String())
	return imgByte
}

// func (h *heatmap) getTable() fyne.CanvasObject {
// 	t := widget.NewTable(
// 		func() (int, int) {
// 			rows, cols := len(h.tasks)+1, len(h.days)+1
// 			return rows, cols
// 		},
// 		func() fyne.CanvasObject {
// 			rect := canvas.NewRectangle(color.White)
// 			return container.NewMax(widget.NewLabel(""), rect)
// 		},
// 		func(id widget.TableCellID, cell fyne.CanvasObject) {
// 			container := cell.(*fyne.Container)
// 			label := container.Objects[0].(*widget.Label)
// 			rect := container.Objects[1].(*canvas.Rectangle)
// 			if id.Col == 0 || id.Row == 0 {
// 				rect.FillColor = color.Transparent
// 			}
// 			if id.Col == 0 && id.Row >= 1 {
// 				label.SetText(h.tasks[id.Row-1])
// 			}
// 			if id.Row == 0 && id.Col >= 1 {
// 				label.SetText(h.days[id.Col-1])
// 			}
// 		},
// 	)
// 	t.SetColumnWidth(0, 102)
// 	return t
// }
