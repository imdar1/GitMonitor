package task

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func InitTaskTab() fyne.CanvasObject {
	return widget.NewLabel("Task Information")
}
