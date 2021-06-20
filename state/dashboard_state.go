package state

import (
	"gitmonitor/sections/auth"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralContent      fyne.CanvasObject
	TaskContent         fyne.CanvasObject
	ContributionContent fyne.CanvasObject
}

func (tabState *TabItemsState) OnWindowLoaded(appData *AppData) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(appData *AppData) {
	// Fetch and checkout
	err := appData.Repo.FetchAndCheckout(auth.AskAuth)
	utils.CheckErr("OnRepositoryLoaded", err)

	// Update task content
	taskContent := tabState.TaskContent.(*widget.Card)
	taskData := task.TaskData{
		Project: appData.SelectedProject,
	}
	taskData.ReadTaskData(appData.Repo, *appData.Database)
	task.RenderTaskTab(taskContent, taskData, appData.Database)

	// Update general content
	generalContent := tabState.GeneralContent.(*widget.Card)
	generalData := general.InitGeneralData(appData.SelectedProject, appData.Repo)
	general.RenderGeneralTab(generalContent, generalData)

	// Update contributor content
	contributionContent := tabState.ContributionContent.(*widget.Card)
	contributionData := contribution.InitContributorData(generalData.Commits, appData.Repo)
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
