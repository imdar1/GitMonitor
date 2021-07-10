package task

import (
	"bytes"
	"fmt"
	"gitmonitor/services/svg2png"
	"gitmonitor/services/utils"
	"time"

	"github.com/gregoryv/draw/design"
	"github.com/gregoryv/draw/types/date"
)

type taskInformation struct {
	startDateStr string
	days         int
}

type timelineData struct {
	taskInformation
	tasks map[string]taskInformation
}

func initData(taskData TaskData) timelineData {
	var data timelineData

	if len(taskData.Tasks) > 0 {
		startDate := time.Unix(taskData.Tasks[0].StartDate, 0)
		endDate := time.Unix(taskData.Tasks[0].EndDate, 0)

		data = timelineData{
			taskInformation: taskInformation{
				startDateStr: utils.GetStringFromDatetime(startDate),
				days:         utils.GetDayDifference(startDate, endDate),
			},
			tasks: make(map[string]taskInformation),
		}
		for _, v := range taskData.Tasks {
			startTask := time.Unix(v.StartDate, 0)
			endTask := time.Unix(v.EndDate, 0)

			if startTask.Before(startDate) {
				startDate = startTask
				data.taskInformation.startDateStr = utils.GetStringFromDatetime(startDate)
			}
			if endTask.After(endDate) {
				endDate = endTask
				data.taskInformation.days = utils.GetDayDifference(startDate, endDate)
			}

			taskInfo := taskInformation{
				startDateStr: utils.GetStringFromDatetime(startTask),
				days:         utils.GetDayDifference(startTask, endTask),
			}
			data.tasks[v.Name] = taskInfo
		}
	}
	return data
}

func (t *timelineData) getGanttChartImage() []byte {
	if len(t.tasks) == 0 {
		return []byte{}
	}

	var bar *design.Task
	ganttChart := design.NewGanttChart(date.String(t.startDateStr), t.days+1)
	fmt.Printf("task Start date:%s, %d\n", t.startDateStr, t.days+1)
	for key, value := range t.tasks {
		bar = ganttChart.Add(key)
		ganttChart.Place(bar).At(date.String(value.startDateStr), value.days)
		// fmt.Println("cekkk")
		fmt.Printf("task name: %s , date: %s , days: %d\n", key, value.startDateStr, value.days)
	}
	// ganttChart.SetCaption("Project Schedule")
	imgBuffer := new(bytes.Buffer)
	styling := ganttChart.Diagram.Style
	styling.SetOutput(imgBuffer)
	ganttChart.WriteSVG(&styling)

	imgByte := svg2png.GetImage(imgBuffer.String())
	return imgByte
}
