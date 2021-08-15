package settings

import (
	"errors"
	"gitmonitor/models"
	"gitmonitor/sections/auth"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"
	"time"

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

func updateProjectDialogue(appData *data.AppData, title string, description string) {
	err := appData.Database.UpdateProject(appData.SelectedProject)
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	dialog.ShowInformation(
		title,
		description,
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func RenderSettingsTab(wrapper fyne.CanvasObject, data SettingsData, appData *data.AppData) {
	defaultBranchName := appData.SelectedProject.DefaultBranchName
	defaultRemoteName := appData.SelectedProject.DefaultRemoteName
	projectStartDate := time.Unix(appData.SelectedProject.ProjectStartDate, 0).Format("02/01/2006")
	projectEndDate := time.Unix(appData.SelectedProject.ProjectEndDate, 0).Format("02/01/2006")

	// Default branch name section
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
		appData.SelectedProject.DefaultBranchName = defaultBranchName
		updateProjectDialogue(
			appData,
			"Success",
			"The default branch name was successfully updated",
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

	// Default remote name section
	defaultRemoteNameEntry := widget.NewEntry()
	defaultRemoteNameEntry.SetText(defaultRemoteName)
	defaultRemoteNameEntry.OnChanged = func(s string) {
		defaultRemoteName = s
	}
	fetchButton := widget.NewButtonWithIcon("Fetch", theme.ViewRefreshIcon(), func() {
		appData.Repo.Fetch(auth.AskAuth, defaultRemoteName)
	})
	saveRemoteButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		appData.SelectedProject.DefaultRemoteName = defaultRemoteName
		updateProjectDialogue(
			appData,
			"Success",
			"The default remote name was successfully updated",
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

	// Set project start and end date section
	projectStartDateEntry := widget.NewEntry()
	projectStartDateEntry.SetText(projectStartDate)
	projectStartDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	projectStartDateEntry.OnChanged = func(s string) {
		projectStartDate = s
	}

	projectEndDateEntry := widget.NewEntry()
	projectEndDateEntry.SetText(projectStartDate)
	projectEndDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	projectEndDateEntry.OnChanged = func(s string) {
		projectEndDate = s
	}

	saveProjectDateButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		startDate, err1 := utils.GetUnixTimeStampFromString(projectStartDate)
		endDate, err2 := utils.GetUnixTimeStampFromString(projectEndDate)

		if err1 != nil || err2 != nil {
			dialog.ShowError(
				errors.New("invalid date format, please use DD/MM/YYYY format"),
				fyne.CurrentApp().Driver().AllWindows()[0],
			)
			return
		}

		appData.SelectedProject.ProjectStartDate = startDate
		appData.SelectedProject.ProjectEndDate = endDate
		updateProjectDialogue(
			appData,
			"Success",
			"Project start date and end date was successfully updated",
		)
	})

	projectDateForm := widget.NewForm(
		[]*widget.FormItem{
			{Text: "Project Start Date", Widget: projectStartDateEntry},
			{Text: "Project End Date", Widget: projectEndDateEntry},
		}...,
	)

	projectDateButtonWrapper := container.NewHBox(
		layout.NewSpacer(),
		saveProjectDateButton,
	)

	projectDateContent := container.NewBorder(
		nil,
		projectDateButtonWrapper,
		nil,
		nil,
		projectDateForm,
	)
	projectDateWrapper := widget.NewCard(
		"Project Start and End Date",
		"Planned project start and end dates",
		projectDateContent,
	)

	settingsWrapper := wrapper.(*widget.Card)
	settingsWrapper.SetContent(
		container.NewVScroll(
			container.NewVBox(
				defaultBranchWrapper,
				defaultRemoteNameWrapper,
				projectDateWrapper,
			),
		),
	)
}
