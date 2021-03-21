package sections

import (
	"gitmonitor/sections/dashboard"

	"gitmonitor/sections/profile"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetContent() fyne.CanvasObject {
	profile := profile.GetProfileWindow()
	dashboard := dashboard.GetDashboardWindow()
	return container.NewBorder(profile, nil, nil, nil, dashboard)
}
