package task

import (
	"gitmonitor/models"
	"gitmonitor/sections/auth"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
)

type TaskData struct {
	Project             models.Project
	Tasks               []models.Task
	Branches            []models.Branch
	AdditionalRenderers []data.Renderer

	wrapper fyne.CanvasObject
}

func InitReadTaskData(wrapper fyne.CanvasObject, appData *data.AppData) TaskData {
	const serviceName = "InitReadTaskData"
	// Initialize tasks list
	tasks := appData.Database.GetTasksData(appData.SelectedProject.ProjectId)

	// get remote branches and sync branches with db
	branches, err := appData.Repo.GetRemoteBranches(auth.AskAuth)
	utils.CheckErr(serviceName, err)
	err = appData.Database.SyncBranches(appData.SelectedProject.ProjectId, branches)
	utils.CheckErr(serviceName, err)

	// get branches stored in db
	branchModels, err := appData.Database.GetBranchesData(appData.SelectedProject.ProjectId)
	utils.CheckErr(serviceName, err)

	// sync task and branch in db
	err = appData.Database.SyncTask(tasks, branchModels)
	utils.CheckErr(serviceName, err)

	// get updated tasks
	tasks = appData.Database.GetTasksData(appData.SelectedProject.ProjectId)

	return TaskData{
		Project:  appData.SelectedProject,
		Tasks:    tasks,
		Branches: branchModels,
		wrapper:  wrapper,
	}
}

func (t *TaskData) RefreshTasksFromTaskData(appData *data.AppData) {
	t.Tasks = appData.Database.GetTasksData(t.Project.ProjectId)
}

func (t TaskData) Render(data *data.AppData) {
	renderTaskTab(t, data)
}
