package task

import (
	"errors"
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/sections/data"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type formData struct {
	projectId     int64
	tempStartDate string
	tempEndDate   string
	tempBranch    string
	task          models.Task
	branches      []models.Branch
}

func getBranchesName(branches []models.Branch) []string {
	var branchesName []string
	for _, v := range branches {
		branchesName = append(branchesName, v.Name)
	}
	return branchesName
}

func getBranchName(branchId int, branches []models.Branch) string {
	name := ""
	for _, v := range branches {
		if v.BranchId == branchId {
			name = v.Name
			break
		}
	}
	return name
}

func getTaskForm(data *formData, onSubmit func(), onClose func()) fyne.CanvasObject {
	taskNameEntry := widget.NewEntry()
	taskNameEntry.SetText(data.task.Name)
	taskNameEntry.SetPlaceHolder("Improve X Feature")
	taskNameEntry.OnChanged = func(s string) {
		data.task.Name = s
	}

	taskStartDateEntry := widget.NewEntry()
	taskStartDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	if data.task.StartDate > 0 {
		taskStartDateEntry.SetText(time.Unix(data.task.StartDate, 0).Format("02/01/2006"))
	}
	taskStartDateEntry.OnChanged = func(s string) {
		data.tempStartDate = s
	}

	taskEndDateEntry := widget.NewEntry()
	taskEndDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	if data.task.EndDate > 0 {
		taskEndDateEntry.SetText(time.Unix(data.task.EndDate, 0).Format("02/01/2006"))
	}
	taskEndDateEntry.OnChanged = func(s string) {
		data.tempEndDate = s
	}

	taskAssigneeNameEntry := widget.NewEntry()
	taskAssigneeNameEntry.SetPlaceHolder("Username of the task author")
	taskAssigneeNameEntry.SetText(data.task.AssigneeName)
	taskAssigneeNameEntry.OnChanged = func(s string) {
		data.task.AssigneeName = s
	}

	taskAssigneeEmailEntry := widget.NewEntry()
	taskAssigneeEmailEntry.SetPlaceHolder("Email of the task author")
	taskAssigneeEmailEntry.SetText(data.task.AssigneeEmail)
	taskAssigneeEmailEntry.OnChanged = func(s string) {
		data.task.AssigneeEmail = s
	}

	branchesName := getBranchesName(data.branches)
	taskBranchEntry := widget.NewSelectEntry(branchesName) // TODO: read branches from db
	taskBranchEntry.SetText(getBranchName(data.task.BranchId, data.branches))
	taskBranchEntry.SetPlaceHolder("Feature branch to work on this task")
	taskBranchEntry.OnChanged = func(s string) {
		data.tempBranch = s
	}

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

func validateTask(task models.Task) bool {
	return task.AssigneeEmail != "" && task.AssigneeName != "" &&
		task.EndDate != 0 && task.StartDate != 0 &&
		task.Name != "" && task.ProjectId != 0
}

func showModifyTaskWindow(selectedTask models.Task, appData *data.AppData) {
	w := fyne.CurrentApp().NewWindow("Edit a task")
	// TODO
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showTaskWindow(taskWrapper fyne.CanvasObject, taskData TaskData, appData *data.AppData) {
	w := fyne.CurrentApp().NewWindow("Add a new task")

	data := &formData{
		projectId: taskData.Project.ProjectId,
		task:      models.Task{},
		branches:  taskData.Branches,
	}
	taskForm := getTaskForm(
		data,
		func() {
			// Validate start date and end date, then assign both to data
			startDate, err := time.Parse("02/01/2006", data.tempStartDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			endDate, err := time.Parse("02/01/2006", data.tempEndDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			data.task.StartDate = startDate.Unix()
			data.task.EndDate = endDate.Unix()

			// Get branch id for selected branch
			data.task.BranchId = appData.Database.GetBranchIdByName(data.tempBranch)
			data.task.TaskStatus = int(constants.Waiting)
			data.task.ProjectId = taskData.Project.ProjectId

			isValid := validateTask(data.task)
			if !isValid {
				dialog.ShowError(errors.New("invalid data"), w)
				return
			}

			err = appData.Database.AddTask(data.task)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Re-render task
			taskData.RefreshTasksFromTaskData(appData)
			RenderTaskTab(taskWrapper, taskData, appData)

			dialog.ShowInformation("Success", "Task was successfully added", w)
			w.Close()
		},
		func() {
			w.Close()
		},
	)
	wrapper := widget.NewCard("", "", taskForm)
	w.SetContent(wrapper)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(600, 300))
	w.Show()
}
