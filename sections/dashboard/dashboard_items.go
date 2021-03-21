package dashboard

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	DashboardItems = []*container.TabItem{
		container.NewTabItem("General", widget.NewLabel("General Data")),
		container.NewTabItem("Task", widget.NewLabel("Task Data")),
		container.NewTabItem("Contribution", widget.NewLabel("Contribution Data")),
	}
)
