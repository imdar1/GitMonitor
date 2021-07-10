package task

import (
	"gitmonitor/models"
	"gitmonitor/sections/auth"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"
)

type TaskData struct {
	Project  models.Project
	Tasks    []models.Task
	Branches []models.Branch
}

func InitReadTaskData(appData *data.AppData) TaskData {
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
	}
}

func (t *TaskData) RefreshTaskData(appData *data.AppData) {

}
