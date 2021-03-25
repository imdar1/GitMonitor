package sections

import (
	"gitmonitor/db"
	"gitmonitor/sections/dashboard"

	"gitmonitor/sections/profile"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetContent(w fyne.Window, dbConfig db.DBConfig) fyne.CanvasObject {
	profile := profile.GetProfileWindow(w)
	dashboard := dashboard.GetDashboardWindow(w)
	return container.NewBorder(profile, nil, nil, nil, dashboard)
}
