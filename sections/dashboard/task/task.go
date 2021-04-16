package task

import (
	"bytes"
	"gitmonitor/models"
	"gitmonitor/services"
	"image"

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

func getTaskObject(taskData TaskData) fyne.CanvasObject {
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
	return list
}

func RenderTaskTab(taskData TaskData) fyne.CanvasObject {
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
	taskContentBottom := container.NewHSplit(
		getTaskObject(tDummy),
		widget.NewLabel("Infomasi Task"),
	)
	taskContent := container.NewVSplit(taskContentTop, taskContentBottom)
	addTaskButton := widget.NewButton("Add Task", func() {})
	setBranchButton := widget.NewButton("Set Branch", func() {})
	actionButton := container.NewHBox(layout.NewSpacer(), addTaskButton, setBranchButton)
	taskContentWrapper := container.NewBorder(nil, actionButton, nil, nil, taskContent)
	return taskContentWrapper
}
