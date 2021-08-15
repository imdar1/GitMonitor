package getting_started

import (
	"errors"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GettingStartedForm(appData *data.AppData) {
	var projectStartDate string
	var projectEndDate string

	done := make(chan struct{})

	w := fyne.CurrentApp().NewWindow("Getting started with your new project")

	projectStartDateEntry := widget.NewEntry()
	projectStartDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	projectStartDateEntry.OnChanged = func(s string) {
		projectStartDate = s
	}

	projectEndDateEntry := widget.NewEntry()
	projectEndDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	projectEndDateEntry.OnChanged = func(s string) {
		projectEndDate = s
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Project Start Date", Widget: projectStartDateEntry},
			{Text: "Project End Date", Widget: projectEndDateEntry},
		},
	}
	okButton := widget.NewButton("OK", func() {
		startDate, err1 := utils.GetUnixTimeStampFromString(projectStartDate)
		endDate, err2 := utils.GetUnixTimeStampFromString(projectEndDate)
		if err1 != nil || err2 != nil {
			dialog.ShowError(errors.New("invalid date format, please use DD/MM/YYYY format"), w)
			return
		}

		appData.SelectedProject.ProjectStartDate = startDate
		appData.SelectedProject.ProjectEndDate = endDate
		appData.Database.UpdateProject(appData.SelectedProject)
		w.Close()
	})
	container := container.NewBorder(
		nil, container.NewHBox(layout.NewSpacer(), okButton), nil, nil, form,
	)
	w.SetContent(container)
	w.SetOnClosed(func() {
		close(done)
	})
	w.Resize(fyne.NewSize(500, w.Content().Size().Height))
	w.CenterOnScreen()
	w.Show()

	<-done
}
