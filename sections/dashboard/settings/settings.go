package settings

import (
	"gitmonitor/models"
	"gitmonitor/sections/auth"
	"gitmonitor/sections/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitSettingsTab() fyne.CanvasObject {
	return widget.NewLabel("Settings")
}

func RenderSettingsTab(wrapper fyne.CanvasObject, data SettingsData, appData *data.AppData) {
	defaultBranchName := appData.SelectedProject.DefaultBranchName
	defaultRemoteName := appData.SelectedProject.DefaultRemoteName

	remoteBranchesName := models.GetBranchesName(data.RemoteBranches)
	defaultBranchEntry := widget.NewSelectEntry(remoteBranchesName)
	defaultBranchEntry.SetText(defaultBranchName)
	defaultBranchEntry.OnChanged = func(s string) {
		defaultBranchName = s
	}

	checkoutButton := widget.NewButtonWithIcon("Checkout", theme.ViewRefreshIcon(), func() {
		appData.Repo.Checkout(defaultBranchName)
	})
	saveBranchButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		err := appData.Database.UpdateDefaultBranchName(defaultBranchName, appData.SelectedProject.ProjectId)
		if err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		dialog.ShowInformation(
			"Success",
			"The default branch name was successfully updated",
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
	})
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
	defaultRemoteNameEntry.SetText(defaultRemoteName)
	defaultRemoteNameEntry.OnChanged = func(s string) {
		defaultRemoteName = s
	}
	fetchButton := widget.NewButtonWithIcon("Fetch", theme.ViewRefreshIcon(), func() {
		appData.Repo.Fetch(auth.AskAuth, defaultRemoteName)
	})
	saveRemoteButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		err := appData.Database.UpdateDefaultRemoteName(defaultRemoteName, appData.SelectedProject.ProjectId)
		if err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		dialog.ShowInformation(
			"Success",
			"The default remote name was successfully updated",
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
	})
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
