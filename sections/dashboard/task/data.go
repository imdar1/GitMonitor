package task

import (
	"gitmonitor/db"
	"gitmonitor/models"
)

type TaskData struct {
	Project  models.Project
	Tasks    []models.Task
	Branches []models.Branch
}

func (t *TaskData) ReadTaskData(db db.DBConfig) {
	tasks := db.GetTasksData(t.Project.ProjectId)
	branches := db.GetBranchesData(t.Project.ProjectId)

	t.Tasks = tasks
	t.Branches = branches
}
