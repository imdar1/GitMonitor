package dashboard

import (
	"gitmonitor/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetDashboardWindow(_ fyne.Window, appState *state.AppState) fyne.CanvasObject {
	appState.DashboardState = state.InitTabItems()
	dashboardItems := []*container.TabItem{
		appState.DashboardState.GeneralTab,
		appState.DashboardState.TaskTab,
		appState.DashboardState.ContributionTab,
	}
	content := container.NewAppTabs(dashboardItems...)
	return content
}
