package state

import (
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralTab      *container.TabItem
	TaskTab         *container.TabItem
	ContributionTab *container.TabItem
}

func (tabState *TabItemsState) OnDatabaseLoaded() {
	tabState.GeneralTab.Content = widget.NewLabel("Database loaded")
	tabState.GeneralTab.Content.Refresh()
}

func (tabState *TabItemsState) OnRepositoryLoaded() {
	tabState.GeneralTab.Content = widget.NewLabel("Repository loaded")
	tabState.GeneralTab.Content.Refresh()
}

func InitTabItems() TabItemsState {
	tabItems := TabItemsState{
		GeneralTab:      container.NewTabItem("General", general.InitGeneralTab()),
		TaskTab:         container.NewTabItem("Task", task.InitTaskTab()),
		ContributionTab: container.NewTabItem("Contribution", contribution.InitContributionTab()),
	}
	return tabItems
}
