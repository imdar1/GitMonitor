package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow() fyne.CanvasObject {
	content := container.NewAppTabs(DashboardItems...)
	return content
}
