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
	// Update task content
	taskContent := tabState.TaskContent.(*widget.Card)
	taskData := task.TaskData{
		Project: project,
	}
	taskData.ReadTaskData(appData.Repo, *appData.Database)
	task.RenderTaskTab(taskContent, taskData, appData.Database)

	// Update general content
	generalContent := tabState.GeneralContent.(*widget.Card)
	generalData := general.InitGeneralData(project, appData.Repo)
	general.RenderGeneralTab(generalContent, generalData)

	// Update contributor content
	contributionContent := tabState.ContributionContent.(*widget.Card)
	contributionData := contribution.InitContributorData(appData.Repo)
	contribution.RenderContributorTab(contributionContent, contributionData)
}

func InitTabItems() TabItemsState {

	tabItems := TabItemsState{
		GeneralContent:      widget.NewCard("", "", general.InitGeneralTab()),
		TaskContent:         widget.NewCard("", "", task.InitTaskTab()),
		ContributionContent: widget.NewCard("", "", contribution.InitContributionTab()),
	}
	return tabItems
}
