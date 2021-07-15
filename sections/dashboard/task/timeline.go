package task

import (
	"bytes"
	"gitmonitor/constants"
	"gitmonitor/services/svg2png"
	"gitmonitor/services/utils"
	"time"

	"github.com/gregoryv/draw/design"
	"github.com/gregoryv/draw/types/date"
)

type taskInformation struct {
	taskName     string
	startDateStr string
	taskStatus   constants.TaskStatus
	days         int
}

type timelineData struct {
	startDateStr string
	days         int
	tasks        []taskInformation
}

func initData(taskData TaskData) timelineData {
	var data timelineData

	if len(taskData.Tasks) > 0 {
		startDate := utils.GetBeginningOfMonthByTime(
			time.Unix(taskData.Tasks[0].StartDate, 0),
		)
		endDate := time.Unix(taskData.Tasks[0].EndDate, 0)

		data = timelineData{
			startDateStr: utils.GetStringFromDatetime(startDate),
			days:         utils.GetDayDifference(startDate, endDate),
			tasks:        make([]taskInformation, 0),
		}
		for _, v := range taskData.Tasks {
			startTask := time.Unix(v.StartDate, 0)
			endTask := time.Unix(v.EndDate, 0)

			if startTask.Before(startDate) {
				startDate = utils.GetBeginningOfMonthByTime(startTask)
				data.startDateStr = utils.GetStringFromDatetime(startDate)
				data.days = utils.GetDayDifference(startDate, endDate)
			}
			if endTask.After(endDate) {
				endDate = endTask
				data.days = utils.GetDayDifference(startDate, endDate)
			}

			taskInfo := taskInformation{
				startDateStr: utils.GetStringFromDatetime(startTask),
				days:         utils.GetDayDifference(startTask, endTask),
				taskName:     v.Name,
				taskStatus:   constants.TaskStatus(v.TaskStatus),
			}
			data.tasks = append(data.tasks, taskInfo)
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
	for _, value := range t.tasks {
		switch value.taskStatus {
		case constants.Waiting:
			bar = ganttChart.Add(value.taskName).Blue()
		case constants.InProgress:
			bar = ganttChart.Add(value.taskName).Red()
		case constants.Done:
			bar = ganttChart.Add(value.taskName)
		default:
			bar = ganttChart.Add(value.taskName)
		}
		ganttChart.Place(bar).At(date.String(value.startDateStr), value.days)
	}

	ganttChart.SetCaption("Project Schedule")
	imgBuffer := new(bytes.Buffer)
	styling := ganttChart.Diagram.Style
	styling.SetOutput(imgBuffer)
	ganttChart.WriteSVG(&styling)

	imgByte := svg2png.GetImage(imgBuffer.String())
	return imgByte
}
