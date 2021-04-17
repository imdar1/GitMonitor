package task

import (
	"gitmonitor/db"
	"gitmonitor/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FormBinding struct {
}

func getBranchesName(branches []models.Branch) []string {
	var branchesName []string
	for _, v := range branches {
		branchesName = append(branchesName, v.Name)
	}
	return branchesName
}

func getTaskForm(onSubmit func(), onClose func()) fyne.CanvasObject {
	taskNameEntry := widget.NewEntry()
	taskNameEntry.SetPlaceHolder("Improve X Feature")
	taskStartDateEntry := widget.NewEntry()
	taskStartDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	taskEndDateEntry := widget.NewEntry()
	taskEndDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	taskAssigneeNameEntry := widget.NewEntry()
	taskAssigneeNameEntry.SetPlaceHolder("Username of the task author")
	taskAssigneeEmailEntry := widget.NewEntry()
	taskAssigneeEmailEntry.SetPlaceHolder("Email of the task author")
	taskBranchEntry := widget.NewSelectEntry([]string{}) // TODO: read branches from db

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: taskNameEntry},
			{Text: "Start date", Widget: taskStartDateEntry},
			{Text: "End date", Widget: taskEndDateEntry},
			{Text: "Assignee", Widget: taskAssigneeNameEntry},
			{Text: "Assignee mail", Widget: taskAssigneeEmailEntry},
			{Text: "Associated branch:", Widget: taskBranchEntry},
		},
	}

	submitButton := widget.NewButton("Submit", onSubmit)
	closeButton := widget.NewButton("Close", onClose)

	wrapper := container.NewBorder(
		nil,
		container.NewHBox(layout.NewSpacer(), submitButton, closeButton),
		nil,
		nil,
		form,
	)

	return wrapper
}

func showTaskWindow(taskData TaskData, db *db.DBConfig) {
	w := fyne.CurrentApp().NewWindow("Add a new task")
	taskForm := getTaskForm(
		func() {
			// call add to db
			// TODO: add data binding to models task
			db.AddTask(models.Task{})

			// re-render task
			RenderTaskTab(taskData, db)
		},
		func() {
			w.Close()
		},
	)
	wrapper := widget.NewCard("", "", taskForm)
	w.SetContent(wrapper)
	w.CenterOnScreen()
	w.Show()
}
