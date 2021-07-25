package state

import (
	"gitmonitor/sections/auth"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/settings"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralContent      fyne.CanvasObject
	TaskContent         fyne.CanvasObject
	ContributionContent fyne.CanvasObject
	SettingsContent     fyne.CanvasObject
}

func (tabState *TabItemsState) OnWindowLoaded(appData *data.AppData) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(appData *data.AppData) {
	// Fetch and checkout
	err := appData.Repo.FetchAndCheckout(
		auth.AskAuth,
		appData.SelectedProject.DefaultBranchName,
		appData.SelectedProject.DefaultRemoteName,
	)
	utils.CheckErr("OnRepositoryLoaded", err)

	// Update task content
	taskContent := tabState.TaskContent.(*widget.Card)
	taskData := task.InitReadTaskData(appData)
	task.RenderTaskTab(taskContent, taskData, appData)

	// Update general content
	generalContent := tabState.GeneralContent.(*widget.Card)
	generalData := general.InitGeneralData(appData)
	general.RenderGeneralTab(generalContent, generalData)

	// Update contributor content
	contributionContent := tabState.ContributionContent.(*widget.Card)
	contributionData := contribution.InitContributorData(
		generalData.Commits,
		taskData.Tasks,
		appData.SelectedProject.DefaultBranchName,
		appData.SelectedProject.DefaultRemoteName,
	)
	contribution.RenderContributorTab(contributionContent, contributionData, appData)

	// Update settings tab
	settingsContent := tabState.SettingsContent.(*widget.Card)
	settingsData := settings.SettingsData{
		RemoteBranches: taskData.Branches,
	}
	settings.RenderSettingsTab(settingsContent, settingsData, appData)
}

func InitTabItems() TabItemsState {

	tabItems := TabItemsState{
		GeneralContent:      widget.NewCard("", "", general.InitGeneralTab()),
		TaskContent:         widget.NewCard("", "", task.InitTaskTab()),
		ContributionContent: widget.NewCard("", "", contribution.InitContributionTab()),
		SettingsContent:     widget.NewCard("", "", settings.InitSettingsTab()),
	}
	return tabItems
}
