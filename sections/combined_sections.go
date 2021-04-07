package sections

import (
	"gitmonitor/config"
	"gitmonitor/db"
	"gitmonitor/sections/dashboard"

	"gitmonitor/sections/profile"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetContent(w fyne.Window, dbConfig db.DBConfig) fyne.CanvasObject {
	gitConfig := config.GitConfig{}
	profile := profile.GetProfileWindow(w, &gitConfig)
	dashboard := dashboard.GetDashboardWindow(w, &gitConfig)
	return container.NewBorder(profile, nil, nil, nil, dashboard)
}
