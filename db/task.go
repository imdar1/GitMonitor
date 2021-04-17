package db

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services"
)

func (db *DBConfig) GetTasksData(projectId int64) []models.Task {
	var tasks []models.Task
	query := fmt.Sprintf("SELECT * FROM task WHERE project_id=%d ORDER BY start_date ASC;", projectId)
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

func (db *DBConfig) AddTask(task models.Task) error {
	insertQuery := fmt.Sprintf(
		`INSERT INTO task(
			project_id, 
			branch_id, 
			name, 
			assignee_name, 
			assignee_email, 
			task_status, 
			start_date,
			end_date
		)
		VALUES(%d, %d, '%s', '%s', '%s', %d, %d, %d);`,
		task.ProjectId,
		task.BranchId,
		task.Name,
		task.AssigneeName,
		task.AssigneeEmail,
		task.TaskStatus,
		task.StartDate,
		task.EndDate,
	)
	statement, err := db.Driver.Prepare(insertQuery)
	services.CheckErr(err)

	_, err = statement.Exec()
	services.CheckErr(err)

	return err
}
