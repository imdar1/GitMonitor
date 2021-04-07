package dashboard

import (
	"gitmonitor/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow(_ fyne.Window, _ *services.GitConfig) fyne.CanvasObject {
	content := container.NewAppTabs(getTabItems()...)
	return content
}
