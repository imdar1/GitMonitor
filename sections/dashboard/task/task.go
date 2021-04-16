package task

import (
	"bytes"
	"gitmonitor/constants"
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/services"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TaskData struct {
	Tasks    []models.Task
	Branches []models.Branch
}

func InitTaskTab() fyne.CanvasObject {
	return widget.NewLabel("Task Information")
}

func getTasksListCanvas(taskData TaskData, taskInfoCanvas fyne.CanvasObject, db *db.DBConfig) fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(taskData.Tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Task Name"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(taskData.Tasks[id].Name)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selectedTask := taskData.Tasks[id]
		selectedBranch := db.GetBranchById(selectedTask.ProjectId)
		taskDetail := taskInfoCanvas.(*container.Scroll)
		taskDetail.Content = getTaskDetailCanvas(selectedTask, selectedBranch)
		taskDetail.Refresh()
	}
	list.OnUnselected = func(id widget.ListItemID) {
		taskDetail := taskInfoCanvas.(*container.Scroll)
		taskDetail.Content = widget.NewLabel("Infomasi Task")
		taskDetail.Refresh()
	}
	return list
}

func getTaskDetailCanvas(selectedTask models.Task, selectedBranch models.Branch) fyne.CanvasObject {
	startDate := time.Unix(selectedTask.StartDate, 0)
	endDate := time.Unix(selectedTask.EndDate, 0)
	taskNameLabel := widget.NewLabel(selectedTask.Name)
	taskStartDateLabel := widget.NewLabel(startDate.Format("dd/mm/yyyy"))
	taskEndDateLabel := widget.NewLabel(endDate.Format("dd/mm/yyyy"))
	taskAssigneeNameLabel := widget.NewLabel(selectedTask.AssigneeName)
	taskAssigneeEmailLabel := widget.NewLabel(selectedTask.AssigneeEmail)
	taskBranchLabel := widget.NewLabel(selectedBranch.Name)
	taskStatusLabel := widget.NewLabel(constants.TaskStatusMap[int(selectedTask.TaskStatus)])

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: taskNameLabel},
			{Text: "Start date", Widget: taskStartDateLabel},
			{Text: "End date", Widget: taskEndDateLabel},
			{Text: "Assignee", Widget: taskAssigneeNameLabel},
			{Text: "Assignee mail", Widget: taskAssigneeEmailLabel},
			{Text: "Associated branch:", Widget: taskBranchLabel},
			{Text: "Status", Widget: taskStatusLabel},
		},
	}
	return form
}

func RenderTaskTab(taskData TaskData, db *db.DBConfig) fyne.CanvasObject {
	// timeData := initData(taskData)
	timeData := initDummy()
	svgString := timeData.getGanttChartImage()

	var ganttChartCanvas fyne.CanvasObject
	if len(svgString) == 0 {
		ganttChartCanvas = widget.NewLabel("No task found")
	} else {
		byteImg := timeData.getGanttChartImage()
		ganttChartImg, _, err := image.Decode(bytes.NewReader(byteImg))
		services.CheckErr(err)

		ganttChartObj := canvas.NewImageFromImage(ganttChartImg)
		ganttChartObj.FillMode = canvas.ImageFillOriginal
		ganttChartCanvas = ganttChartObj
	}

	tDummy := TaskData{
		Tasks: []models.Task{
			{
				Name: "abc",
			}, {
				Name: "def",
			}, {
				Name: "ghi",
			},
		},
	}
	taskContentTop := container.NewVScroll(ganttChartCanvas)
	taskDetail := container.NewVScroll(widget.NewLabel("Infomasi Task"))
	taskContentBottom := container.NewHSplit(
		getTasksListCanvas(tDummy, taskDetail, db),
		taskDetail,
	)
	taskContent := container.NewVSplit(taskContentTop, taskContentBottom)
	addTaskButton := widget.NewButton("Add Task", func() {})
	setBranchButton := widget.NewButton("Set Branch", func() {})
	actionButton := container.NewHBox(layout.NewSpacer(), addTaskButton, setBranchButton)
	taskContentWrapper := container.NewBorder(nil, actionButton, nil, nil, taskContent)
	return taskContentWrapper
}
