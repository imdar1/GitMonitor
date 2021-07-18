package settings

import (
	"gitmonitor/models"
	"gitmonitor/sections/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitSettingsTab() fyne.CanvasObject {
	return widget.NewLabel("Settings")
}

func RenderSettingsTab(wrapper fyne.CanvasObject, data SettingsData, appData *data.AppData) {
	remoteBranchesName := models.GetBranchesName(data.RemoteBranches)
	defaultBranchEntry := widget.NewSelectEntry(remoteBranchesName)
	defaultBranchEntry.SetText(appData.SelectedProject.DefaultBranchName)
	checkoutButton := widget.NewButtonWithIcon("Checkout", theme.ViewRefreshIcon(), func() {})
	saveBranchButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {})
	defaultBranchButtonWrapper := container.NewHBox(layout.NewSpacer(), checkoutButton, saveBranchButton)
	defaultBranchContent := container.NewBorder(
		nil,
		defaultBranchButtonWrapper,
		nil,
		nil,
		defaultBranchEntry,
	)
	defaultBranchWrapper := widget.NewCard(
		"Default Branch",
		"Set Git default branch to monitor (default = master)",
		defaultBranchContent,
	)

	defaultRemoteNameEntry := widget.NewEntry()
	defaultRemoteNameEntry.SetText(appData.SelectedProject.DefaultRemoteName)
	fetchButton := widget.NewButtonWithIcon("Fetch", theme.ViewRefreshIcon(), func() {})
	saveRemoteButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {})
	defaultRemoteButtonWrapper := container.NewHBox(layout.NewSpacer(), fetchButton, saveRemoteButton)
	defaultRemoteNameContent := container.NewBorder(
		nil,
		defaultRemoteButtonWrapper,
		nil,
		nil,
		defaultRemoteNameEntry,
	)
	defaultRemoteNameWrapper := widget.NewCard(
		"Default Remote Name",
		"Set Git default remote name (default = origin)",
		defaultRemoteNameContent,
	)

	settingsWrapper := wrapper.(*widget.Card)
	settingsWrapper.SetContent(container.NewVBox(defaultBranchWrapper, defaultRemoteNameWrapper))
}
