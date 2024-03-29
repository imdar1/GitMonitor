package settings

import (
	"errors"
	"gitmonitor/models"
	"gitmonitor/sections/dashboard/general"
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

func renderSettingsTab(data SettingsData, appData *data.AppData) {
	defaultBranchName := appData.SelectedProject.DefaultBranchName
	defaultRemoteName := appData.SelectedProject.DefaultRemoteName
	projectStartDate := time.Unix(appData.SelectedProject.ProjectStartDate, 0).Format("02/01/2006")
	projectEndDate := time.Unix(appData.SelectedProject.ProjectEndDate, 0).Format("02/01/2006")
	if appData.SelectedProject.ProjectStartDate == 0 {
		projectStartDate = ""
	}
	if appData.SelectedProject.ProjectEndDate == 0 {
		projectEndDate = ""
	}

	// Default branch name section
	availableRemoteBranches := models.GetAvailableBranches(data.RemoteBranches)
	remoteBranchesName := models.GetBranchesName(availableRemoteBranches)
	defaultBranchEntry := widget.NewSelect(remoteBranchesName, func(s string) {
		defaultBranchName = s
	})
	defaultBranchEntry.SetSelected(defaultBranchName)

	saveBranchButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		appData.SelectedProject.DefaultBranchName = defaultBranchName
		updateProjectDialogue(
			appData,
			"Success",
			"The default branch name was successfully updated",
		)
	})
	defaultBranchButtonWrapper := container.NewHBox(layout.NewSpacer(), saveBranchButton)
	defaultBranchContent := container.NewBorder(
		nil,
		defaultBranchButtonWrapper,
		nil,
		nil,
		defaultBranchEntry,
	)
	defaultBranchWrapper := widget.NewCard(
		"Default Branch",
		"Set Git default branch to monitor (default = master). Please reload the project to see the changes that have been applied.",
		defaultBranchContent,
	)

	// Default remote name section
	defaultRemoteNameEntry := widget.NewEntry()
	defaultRemoteNameEntry.SetText(defaultRemoteName)
	defaultRemoteNameEntry.OnChanged = func(s string) {
		defaultRemoteName = s
	}
	saveRemoteButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		appData.SelectedProject.DefaultRemoteName = defaultRemoteName
		updateProjectDialogue(
			appData,
			"Success",
			"The default remote name was successfully updated",
		)
	})
	defaultRemoteButtonWrapper := container.NewHBox(layout.NewSpacer(), saveRemoteButton)
	defaultRemoteNameContent := container.NewBorder(
		nil,
		defaultRemoteButtonWrapper,
		nil,
		nil,
		defaultRemoteNameEntry,
	)
	defaultRemoteNameWrapper := widget.NewCard(
		"Default Remote Name",
		"Set Git default remote name (default = origin). Please reload the project to see the changes that have been applied.",
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
	projectEndDateEntry.SetText(projectEndDate)
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

		// Re-render general tab
		generalRenderer := data.AdditionalRenderers[0].(general.GeneralData)
		generalRenderer.UpdateProjectStartDate(appData)
		generalRenderer.UpdateProjectEndDate(appData)
		generalRenderer.UpdateProjectTaskStatus(appData)
		generalRenderer.Render(appData)

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

	settingsWrapper := data.Wrapper.(*widget.Card)
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
