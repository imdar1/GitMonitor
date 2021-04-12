package db

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services"
)

func (db *DBConfig) GetTasksData(projectId int64) []models.Task {
	var tasks []models.Task
	query := fmt.Sprintf("SELECT * FROM task WHERE project_id=%d;", projectId)
	rows, err := db.Driver.Query(query)
	services.CheckErr(err)

	if rows != nil {
		for rows.Next() {
			var task models.Task
			err = rows.Scan(
				&task.TaskId,
				&task.ProjectId,
				&task.BranchId,
				&task.Name,
				&task.AssigneeName,
				&task.AssigneeEmail,
				&task.TaskStatus,
				&task.StartDate,
				&task.EndDate,
			)
			services.CheckErr(err)
			tasks = append(tasks, task)
		}
		rows.Close()
	}
	return tasks
}