package state

import (
	"gitmonitor/sections/auth"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/settings"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/sections/data"
	"gitmonitor/sections/getting_started"
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
	if appData.SelectedProject.IsFirstTime {
		getting_started.GettingStartedForm(appData)
	}

	// Fetch and checkout
	err := appData.Repo.FetchAndCheckout(
		auth.AskAuth,
		appData.SelectedProject.DefaultBranchName,
		appData.SelectedProject.DefaultRemoteName,
	)
	utils.CheckErr("OnRepositoryLoaded", err)

	// Init task content
	taskContent := tabState.TaskContent.(*widget.Card)
	taskRenderer := task.InitReadTaskData(taskContent, appData)

	// Init general content
	generalContent := tabState.GeneralContent.(*widget.Card)
	generalRenderer := general.InitGeneralData(generalContent, taskRenderer.Tasks, appData)

	// Init contributor content
	contributionContent := tabState.ContributionContent.(*widget.Card)
	contributionRenderer := contribution.InitContributorData(
		contributionContent,
		generalRenderer.Commits,
		taskRenderer.Tasks,
		appData.SelectedProject.DefaultBranchName,
		appData.SelectedProject.DefaultRemoteName,
	)

	// Update settings tab
	settingsContent := tabState.SettingsContent.(*widget.Card)
	settingsRenderer := settings.InitSettingsData(settingsContent, taskRenderer.Branches)

	// Add necessary renderers
	taskRenderer.AdditionalRenderers = []data.Renderer{
		generalRenderer,
		contributionRenderer,
	}
	settingsRenderer.AdditionalRenderers = []data.Renderer{
		generalRenderer,
	}

	// Render general, task, contribution, and settings
	generalRenderer.Render(appData)
	taskRenderer.Render(appData)
	contributionRenderer.Render(appData)
	settingsRenderer.Render(appData)
}

func initRenderer(caption string) fyne.CanvasObject {
	return widget.NewLabel(caption)
}

func InitTabItems() TabItemsState {
	tabItems := TabItemsState{
		GeneralContent:      widget.NewCard("", "", initRenderer("General Information")),
		TaskContent:         widget.NewCard("", "", initRenderer("Task Information")),
		ContributionContent: widget.NewCard("", "", initRenderer("Contribution Information")),
		SettingsContent:     widget.NewCard("", "", initRenderer("Settings")),
	}
	return tabItems
}
