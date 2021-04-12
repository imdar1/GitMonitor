package state

import (
	"gitmonitor/db"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralContent      *fyne.Container
	TaskContent         *fyne.Container
	ContributionContent *fyne.Container
}

func (tabState *TabItemsState) OnDatabaseLoaded(db *db.DBConfig) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(repo services.GitConfig, db *db.DBConfig) {
	tabState.GeneralContent.Objects = nil
	tabState.GeneralContent.Add(widget.NewLabel("Repository loaded"))
	tabState.GeneralContent.Refresh()
}

func InitTabItems() TabItemsState {
	tabItems := TabItemsState{
		GeneralContent:      container.NewVBox(general.InitGeneralTab()),
		TaskContent:         container.NewVBox(task.InitTaskTab()),
		ContributionContent: container.NewVBox(contribution.InitContributionTab()),
	}
	return tabItems
}
