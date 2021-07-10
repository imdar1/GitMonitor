package sections

import (
	"gitmonitor/db"
	"gitmonitor/sections/dashboard"
	"gitmonitor/sections/data"
	"gitmonitor/state"

	"gitmonitor/sections/profile"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetContent(w fyne.Window, dbConfig *db.DBConfig) fyne.CanvasObject {
	appState := state.AppState{
		AppData: data.AppData{
			Database: dbConfig,
		},
	}
	profile := profile.GetProfileWindow(w, &appState)
	dashboard := dashboard.GetDashboardWindow(w, &appState)
	appState.OnWindowLoaded()

	return container.NewBorder(profile, nil, nil, nil, dashboard)
}
