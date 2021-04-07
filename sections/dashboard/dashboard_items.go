package dashboard

import (
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"

	"fyne.io/fyne/v2/container"
)

type tabItems struct {
	GeneralTab      *container.TabItem
	TaskTab         *container.TabItem
	ContributionTab *container.TabItem
}

func getTabItems() []*container.TabItem {
	tabItems := tabItems{
		GeneralTab:      container.NewTabItem("General", general.InitGeneralTab()),
		TaskTab:         container.NewTabItem("Task", task.InitTaskTab()),
		ContributionTab: container.NewTabItem("Contribution", contribution.InitContributionTab()),
	}
	dashboardItems := []*container.TabItem{
		tabItems.GeneralTab,
		tabItems.TaskTab,
		tabItems.ContributionTab,
	}
	return dashboardItems
}
