package dashboard

import (
	"gitmonitor/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow(_ fyne.Window, _ *config.GitConfig) fyne.CanvasObject {
	content := container.NewAppTabs(getTabItems()...)
	return content
}
