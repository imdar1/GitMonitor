package db

import (
	"fmt"
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/services/utils"
)

func (db *DBConfig) GetTasksData(projectId int64) []models.Task {
	var tasks []models.Task
	const serviceName = "GetTasksData"
	query := fmt.Sprintf("SELECT * FROM task WHERE project_id=%d ORDER BY start_date ASC;", projectId)
	rows, err := db.Driver.Query(query)
	utils.CheckErr(serviceName, err)

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
			utils.CheckErr(serviceName, err)
			tasks = append(tasks, task)
		}
		rows.Close()
	}
	return tasks
}

func (db *DBConfig) AddTask(task models.Task) error {
	const serviceName = "AddTasks"
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
	utils.CheckErr(serviceName, err)

	_, err = statement.Exec()
	utils.CheckErr(serviceName, err)

	return err
}

func (db *DBConfig) UpdateTask(task models.Task) error {
	return nil
}

func getBranchIdList(branches []models.Branch) []int {
	var branchIdList []int
	for _, v := range branches {
		branchIdList = append(branchIdList, v.BranchId)
	}
	return branchIdList
}

func (db *DBConfig) taskStatusIsInProgress(branchId int) bool {
	var taskStatus int
	const serviceName = "taskStatusIsInProgress"
	query := fmt.Sprintf("SELECT task_status FROM task WHERE branch_id=%d", branchId)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(&taskStatus)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return false
	}

	return taskStatus == int(constants.InProgress)
}

func (db *DBConfig) SyncTask(tasks []models.Task, branches []models.Branch) error {
	branchIdList := getBranchIdList(branches)
	tx, err := db.Driver.Begin()
	if err != nil {
		return err
	}

	queryTemplate := "UPDATE task SET task_status=%d WHERE branch_id=%d"
	for _, v := range tasks {
		if utils.IsExistInt(v.BranchId, branchIdList) {
			query := fmt.Sprintf(queryTemplate, constants.InProgress, v.BranchId)
			_, err := tx.Exec(query)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if db.taskStatusIsInProgress(v.BranchId) {
				query := fmt.Sprintf(queryTemplate, constants.Done, v.BranchId)
				_, err := tx.Exec(query)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit()
}
