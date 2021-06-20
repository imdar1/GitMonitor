package task

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/sections/auth"
	"gitmonitor/services/git"
	"gitmonitor/services/utils"
)

type TaskData struct {
	Project  models.Project
	Tasks    []models.Task
	Branches []models.Branch
}

func (t *TaskData) ReadTaskData(gitConfig git.GitConfig, db db.DBConfig) {
	const serviceName = "ReadTaskData"
	// Initialize tasks list
	tasks := db.GetTasksData(t.Project.ProjectId)

	// get remote branches and sync branches with db
	branches, err := gitConfig.GetRemoteBranches(auth.AskAuth)
	utils.CheckErr(serviceName, err)
	err = db.SyncBranches(t.Project.ProjectId, branches)
	utils.CheckErr(serviceName, err)

	// get branches stored in db
	branchModels, err := db.GetBranchesData(t.Project.ProjectId)
	utils.CheckErr(serviceName, err)

	// sync task and branch in db
	err = db.SyncTask(tasks, branchModels)
	utils.CheckErr(serviceName, err)

	// get updated tasks
	tasks = db.GetTasksData(t.Project.ProjectId)

	t.Tasks = tasks
	t.Branches = branchModels
}
