package state

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralContent      fyne.CanvasObject
	TaskContent         fyne.CanvasObject
	ContributionContent fyne.CanvasObject
}

func (tabState *TabItemsState) OnDatabaseLoaded(db *db.DBConfig) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(appData AppData, project models.Project) {
	taskContent := tabState.TaskContent.(*widget.Card)

	taskData := task.TaskData{
		Project: project,
	}
	taskData.ReadTaskData(*appData.Database)
	task.RenderTaskTab(taskContent, taskData, appData.Database)

}

func InitTabItems() TabItemsState {

	tabItems := TabItemsState{
		GeneralContent:      widget.NewCard("", "", general.InitGeneralTab()),
		TaskContent:         widget.NewCard("", "", task.InitTaskTab()),
		ContributionContent: widget.NewCard("", "", contribution.InitContributionTab()),
	}
	return tabItems
}
