package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow(_ fyne.Window) fyne.CanvasObject {
	content := container.NewAppTabs(getTabItems()...)
	return content
}
