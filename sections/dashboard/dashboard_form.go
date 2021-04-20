package dashboard

import (
	"gitmonitor/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow(_ fyne.Window, appState *state.AppState) fyne.CanvasObject {
	appState.DashboardState = state.InitTabItems()
	dashboardItems := []*container.TabItem{
		container.NewTabItem("Task", appState.DashboardState.TaskContent),
		container.NewTabItem("General", appState.DashboardState.GeneralContent),
		container.NewTabItem("Contribution", appState.DashboardState.ContributionContent),
	}
	content := container.NewAppTabs(dashboardItems...)
	return content
}
