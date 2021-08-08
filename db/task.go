package db

import (
	"fmt"
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/services/utils"
	"time"
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

func (db *DBConfig) DeleteTask(task models.Task) error {
	const serviceName = "DeleteTask"
	tx, err := db.Driver.Begin()
	if err != nil {
		utils.CheckErr(serviceName, err)
		return err
	}
	queryTemplate := "DELETE FROM task WHERE task_id=%d"
	query := fmt.Sprintf(queryTemplate, task.TaskId)
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DBConfig) UpdateTask(task models.Task) error {
	const serviceName = "UpdateTask"
	tx, err := db.Driver.Begin()
	if err != nil {
		utils.CheckErr(serviceName, err)
		return err
	}
	queryTemplate := `UPDATE task SET 
						branch_id=%d,
						name='%s',
						assignee_name='%s',
						assignee_email='%s',
						task_status=%d,
						start_date=%d,
						end_date=%d WHERE task_id=%d`
	query := fmt.Sprintf(
		queryTemplate,
		task.BranchId,
		task.Name,
		task.AssigneeName,
		task.AssigneeEmail,
		task.TaskStatus,
		task.StartDate,
		task.EndDate,
		task.TaskId,
	)
	_, err = tx.Exec(query)
	if err != nil {
		utils.CheckErr(serviceName, fmt.Errorf(fmt.Sprintf("error message: %s with query: %s", err, query)))
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DBConfig) UpdateTaskStatus(task models.Task, status constants.TaskStatus) error {
	const serviceName = "UpdateTaskStatus"
	tx, err := db.Driver.Begin()
	if err != nil {
		utils.CheckErr(serviceName, err)
		return err
	}
	queryTemplate := "UPDATE task SET task_status=%d WHERE task_id=%d"
	_, err = tx.Exec(fmt.Sprintf(queryTemplate, status, task.TaskId))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func getBranchIdList(branches []models.Branch) []int {
	var branchIdList []int
	for _, v := range branches {
		branchIdList = append(branchIdList, v.BranchId)
	}
	return branchIdList
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
			if v.TaskStatus == int(constants.Waiting) {
				query := fmt.Sprintf(queryTemplate, constants.InProgress, v.BranchId)
				_, err := tx.Exec(query)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		} else {
			if v.TaskStatus == int(constants.InProgress) {
				currentTime := time.Now()
				taskDeadline := time.Unix(v.EndDate, 0)

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

				taskStatus := constants.Done
				if currentTime.After(taskDeadline) {
					taskStatus = constants.DoneLate
				}
				query := fmt.Sprintf(queryTemplate, taskStatus, v.BranchId)
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
