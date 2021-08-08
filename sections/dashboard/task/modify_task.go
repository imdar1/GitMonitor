package task

import (
	"errors"
	"fmt"
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

func getTaskForm(data *formData, onSubmit func(), onClose func()) fyne.CanvasObject {
	taskNameEntry := widget.NewEntry()
	taskNameEntry.SetText(data.task.Name)
	taskNameEntry.SetPlaceHolder("Improve X Feature")
	taskNameEntry.OnChanged = func(s string) {
		data.task.Name = s
	}

	taskStartDateEntry := widget.NewEntry()
	taskStartDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	tempStartDate := time.Unix(data.task.StartDate, 0).Format("02/01/2006")
	data.tempStartDate = tempStartDate
	if data.task.StartDate > 0 {
		taskStartDateEntry.SetText(tempStartDate)
	}
	taskStartDateEntry.OnChanged = func(s string) {
		data.tempStartDate = s
	}

	taskEndDateEntry := widget.NewEntry()
	taskEndDateEntry.SetPlaceHolder("DD/MM/YYYY format")
	tempEndDate := time.Unix(data.task.EndDate, 0).Format("02/01/2006")
	data.tempEndDate = tempEndDate
	if data.task.EndDate > 0 {
		taskEndDateEntry.SetText(tempEndDate)
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

	branchesName := models.GetBranchesName(data.branches)
	taskBranchEntry := widget.NewSelectEntry(branchesName) // TODO: read branches from db
	tempBranch := models.GetBranchName(data.task.BranchId, data.branches)
	data.tempBranch = tempBranch
	taskBranchEntry.SetText(tempBranch)
	taskBranchEntry.SetPlaceHolder("Feature branch to work on this task")
	taskBranchEntry.OnChanged = func(s string) {
		data.tempBranch = s
	}

	taskStatusEntry := widget.NewSelect(constants.TaskStatusList, func(s string) {})
	taskStatusEntry.SetSelectedIndex(data.task.TaskStatus)
	taskStatusEntry.OnChanged = func(s string) {
		data.task.TaskStatus = taskStatusEntry.SelectedIndex()
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: taskNameEntry},
			{Text: "Start date", Widget: taskStartDateEntry},
			{Text: "End date", Widget: taskEndDateEntry},
			{Text: "Assignee", Widget: taskAssigneeNameEntry},
			{Text: "Assignee mail", Widget: taskAssigneeEmailEntry},
			{Text: "Status:", Widget: taskStatusEntry},
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

// Convert string date with "DD/MM/YYYY" format into unix timestamp
func getUnixTimeStampFromString(timeString string) (int64, error) {
	date, err := time.Parse("02/01/2006", timeString)
	if err != nil {
		return 0, err
	}
	return date.Unix(), nil
}

func showModifyTaskWindow(
	taskWrapper fyne.CanvasObject,
	selectedTask models.Task,
	taskData TaskData,
	appData *data.AppData,
) {
	w := fyne.CurrentApp().NewWindow("Edit a task")
	// TODO
	data := &formData{
		projectId: taskData.Project.ProjectId,
		task:      selectedTask,
		branches:  taskData.Branches,
	}
	taskForm := getTaskForm(
		data,
		func() {
			// Validate start date and end date, then assign both to data
			startDate, err := getUnixTimeStampFromString(data.tempStartDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			endDate, err := getUnixTimeStampFromString(data.tempEndDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			data.task.StartDate = startDate
			data.task.EndDate = endDate

			// Get branch id for selected branch
			data.task.BranchId = appData.Database.GetBranchIdByName(data.tempBranch)
			data.task.ProjectId = taskData.Project.ProjectId
			// data.task.TaskStatus = int(constants.Waiting)

			isValid := validateTask(data.task)
			if !isValid {
				dialog.ShowError(errors.New("invalid data"), w)
				return
			}

			currentTime := time.Now()
			taskDeadline := time.Unix(data.task.EndDate, 0)
			// set the deadline to the next day, 00.00
			taskDeadline = time.Date(
				taskDeadline.Year(),
				taskDeadline.Month(),
				taskDeadline.Day()+1,
				0,
				0,
				0,
				0,
				taskDeadline.Location(),
			)

			if data.task.TaskStatus == int(constants.Done) && currentTime.After(taskDeadline) {
				data.task.TaskStatus = int(constants.DoneLate)
			}

			err = appData.Database.UpdateTask(data.task)
			if err != nil {
				fmt.Println(err.Error())
				dialog.ShowError(err, w)
				return
			}

			// Re-render task
			taskData.RefreshTasksFromTaskData(appData)
			RenderTaskTab(taskWrapper, taskData, appData)

			dialog.ShowInformation("Success", "Task was successfully updated", w)
			w.Close()
		},
		func() {
			w.Close()
		},
	)
	wrapper := widget.NewCard("", "", taskForm)
	w.SetContent(wrapper)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showAddTaskWindow(taskWrapper fyne.CanvasObject, taskData TaskData, appData *data.AppData) {
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
			startDate, err := getUnixTimeStampFromString(data.tempStartDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			endDate, err := getUnixTimeStampFromString(data.tempEndDate)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			data.task.StartDate = startDate
			data.task.EndDate = endDate

			// Get branch id for selected branch
			data.task.BranchId = appData.Database.GetBranchIdByName(data.tempBranch)
			data.task.ProjectId = taskData.Project.ProjectId
			// data.task.TaskStatus = int(constants.Waiting)

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
