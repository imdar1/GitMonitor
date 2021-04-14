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
	"fyne.io/fyne/v2/widget"
)

type TaskData struct {
	Tasks    []models.Task
	Branches []models.Branch
}

func InitTaskTab() fyne.CanvasObject {
	return widget.NewLabel("Task Information")
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

	taskContentTop := container.NewVScroll(ganttChartCanvas)
	taskContentBottom := container.NewHSplit(
		widget.NewLabel("Daftar Task"),
		widget.NewLabel("Infomasi Task"),
	)
	taskContent := container.NewVSplit(taskContentTop, taskContentBottom)
	addTaskButton := widget.NewButton("Add Task", func() {})
	setBranchButton := widget.NewButton("Set Branch", func() {})
	actionButton := container.NewHBox(layout.NewSpacer(), addTaskButton, setBranchButton)
	taskContentWrapper := container.NewBorder(nil, actionButton, nil, nil, taskContent)
	return taskContentWrapper
}
