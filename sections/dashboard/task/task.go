package task

import (
	"gitmonitor/models"

	"fyne.io/fyne/v2"
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
	// TODO: Call generate heatmap function using task data, return new container
	return widget.NewLabel("Dummy data")
}
