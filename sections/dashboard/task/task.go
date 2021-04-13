package task

import (
	"gitmonitor/models"

	"fyne.io/fyne/v2"
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
	// TODO: call init heatmap
	heatMap := initData(taskData)
	table := heatMap.getTable()

	taskContentTop := container.NewVScroll(table)
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
