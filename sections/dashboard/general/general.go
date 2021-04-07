package general

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func InitGeneralTab() fyne.CanvasObject {
	return widget.NewLabel("General Information")
}
