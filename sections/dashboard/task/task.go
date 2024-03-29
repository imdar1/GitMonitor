package task

import (
	"bytes"
	"errors"
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitTaskTab() fyne.CanvasObject {
	return widget.NewLabel("Task Information")
}

func getTasksListCanvas(
	taskInfoCanvas fyne.CanvasObject,
	taskData TaskData,
	appData *data.AppData,
	selectedTaskIndex binding.Int,
) fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(taskData.Tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Task Name"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(taskData.Tasks[id].Name)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selectedTask := taskData.Tasks[id]
		selectedBranch := appData.Database.GetBranchById(taskData.Tasks[id].BranchId)
		taskDetail := taskInfoCanvas.(*fyne.Container)
		taskDetail.Remove(taskDetail.Objects[0])
		taskDetail.Add(container.NewScroll(getTaskDetailCanvas(selectedTask, selectedBranch)))
		selectedTaskIndex.Set(id)
		taskDetail.Refresh()
	}
	list.OnUnselected = func(id widget.ListItemID) {
		taskDetail := taskInfoCanvas.(*fyne.Container)
		taskDetail.Remove(taskDetail.Objects[0])
		taskDetail.Add(widget.NewLabel("Infomasi Task"))
		selectedTaskIndex.Set(-1)
		taskDetail.Refresh()
	}
	return list
}

func getTaskDetailCanvas(selectedTask models.Task, selectedBranch models.Branch) fyne.CanvasObject {
	startDate := time.Unix(selectedTask.StartDate, 0)
	endDate := time.Unix(selectedTask.EndDate, 0)
	taskNameLabel := widget.NewLabel(selectedTask.Name)
	taskStartDateLabel := widget.NewLabel(startDate.Format("02/01/2006"))
	taskEndDateLabel := widget.NewLabel(endDate.Format("02/01/2006"))
	taskAssigneeNameLabel := widget.NewLabel(selectedTask.AssigneeName)
	taskAssigneeEmailLabel := widget.NewLabel(selectedTask.AssigneeEmail)
	taskBranchLabel := widget.NewLabel(selectedBranch.Name)
	taskStatusLabel := widget.NewLabel(constants.TaskStatusMap[int(selectedTask.TaskStatus)])

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: taskNameLabel},
			{Text: "Start date", Widget: taskStartDateLabel},
			{Text: "End date", Widget: taskEndDateLabel},
			{Text: "Assignee", Widget: taskAssigneeNameLabel},
			{Text: "Assignee mail", Widget: taskAssigneeEmailLabel},
			{Text: "Associated branch", Widget: taskBranchLabel},
			{Text: "Status", Widget: taskStatusLabel},
		},
	}
	return form
}

func invokeAdditionalRenderers(taskData TaskData, appData *data.AppData) {
	const serviceName = "invokeAdditionalRenderers"

	// Re-render general, contribution, and task tabs
	taskData.RefreshTasksFromTaskData(appData)
	renderTaskTab(taskData, appData)
	for _, renderer := range taskData.AdditionalRenderers {
		switch rendererWithType := renderer.(type) {
		case general.GeneralData:
			generalRenderer := general.InitGeneralData(
				rendererWithType.Wrapper,
				taskData.Tasks,
				appData,
			)
			generalRenderer.Render(appData)
		case contribution.ContributorData:
			contributionRenderer := rendererWithType
			contributionRenderer.SetTasks(taskData.Tasks)
			contributionRenderer.Render(appData)
		default:
			utils.CheckErr(serviceName, errors.New("unidentified type"))
		}
	}
}

// Render task to taskWrapper from given taskData and other operations needed to perform from db
func renderTaskTab(
	taskData TaskData,
	appData *data.AppData,
) {
	timeData := initData(taskData)
	svgString := timeData.getGanttChartImage()
	taskWrapperCard := taskData.wrapper.(*widget.Card)

	var ganttChartCanvas fyne.CanvasObject
	if len(svgString) == 0 {
		ganttChartCanvas = widget.NewLabel("No task found")
	} else {
		ganttChartImg, _, err := image.Decode(bytes.NewReader(svgString))
		utils.CheckErr("RenderTaskTab", err)

		ganttChartObj := canvas.NewImageFromImage(ganttChartImg)
		ganttChartObj.FillMode = canvas.ImageFillOriginal
		ganttChartCanvas = ganttChartObj
	}

	selectedTaskIndex := binding.NewInt()
	selectedTaskIndex.Set(-1)
	taskContentTop := container.NewScroll(ganttChartCanvas)

	taskDetail := container.NewBorder(nil, nil, nil, nil, widget.NewLabel("Informasi task"))
	taskContentBottom := container.NewHSplit(
		container.NewVScroll(getTasksListCanvas(taskDetail, taskData, appData, selectedTaskIndex)),
		taskDetail,
	)
	taskContent := container.NewVSplit(taskContentTop, taskContentBottom)
	addTaskButton := widget.NewButton("Add Task", func() {
		showAddTaskWindow(taskWrapperCard, taskData, appData)
	})

	editTaskButton := widget.NewButton("Edit Task", func() {
		// get Task index from selected task
		taskIndex, err := selectedTaskIndex.Get()
		if err != nil || taskIndex == -1 {
			dialog.ShowError(errors.New("please select a valid task"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		selectedTask := taskData.Tasks[taskIndex]
		showModifyTaskWindow(selectedTask, taskData, appData)
	})

	deleteTaskButton := widget.NewButton("Delete Task", func() {
		// get Task index from selected task
		taskIndex, err := selectedTaskIndex.Get()
		if err != nil || taskIndex == -1 {
			dialog.ShowError(errors.New("please select a valid task"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		selectedTask := taskData.Tasks[taskIndex]
		dialog.ShowConfirm(
			"Confirm Deleting a Task",
			"The selected task is "+selectedTask.Name+
				". Are you sure you want to delete the selected task?",
			func(isOk bool) {
				if isOk {
					err = appData.Database.DeleteTask(selectedTask)
					if err != nil {
						dialog.ShowError(
							errors.New("error when deleting task: "+err.Error()),
							fyne.CurrentApp().Driver().AllWindows()[0],
						)
					}

					// Re-render general, contribution, and task tabs
					invokeAdditionalRenderers(taskData, appData)
				}
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
	})

	actionButton := container.NewHBox(layout.NewSpacer(), addTaskButton, editTaskButton, deleteTaskButton)
	taskContentWrapper := container.NewBorder(nil, actionButton, nil, nil, taskContent)

	taskWrapperCard.SetContent(taskContentWrapper)
	taskWrapperCard.Refresh()
}
