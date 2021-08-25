package db

import (
	"database/sql"
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/utils"
)

func (db *DBConfig) GetProjects() []models.Project {
	const serviceName = "GetProjects"
	var projects []models.Project
	var projectStartDate sql.NullInt64
	var projectEndDate sql.NullInt64

	rows, err := db.Driver.Query("SELECT * FROM project")
	utils.CheckErr(serviceName, err)
	if rows != nil {
		for rows.Next() {
			var project models.Project
			err = rows.Scan(
				&project.ProjectId,
				&project.ProjectDir,
				&projectStartDate,
				&projectEndDate,
				&project.DefaultBranchName,
				&project.DefaultRemoteName,
			)
			utils.CheckErr(serviceName, err)
			project.ProjectStartDate = projectStartDate.Int64
			project.ProjectEndDate = projectEndDate.Int64
			projects = append(projects, project)
		}
		rows.Close()
	}
	return projects
}

func (db *DBConfig) GetProjectByDir(dir string) models.Project {
	var project models.Project
	const serviceName = "GetProjectByDir"
	var projectStartDate sql.NullInt64
	var projectEndDate sql.NullInt64

	query := fmt.Sprintf("SELECT * FROM project WHERE project_dir='%s' LIMIT 1;", dir)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(
		&project.ProjectId,
		&project.ProjectDir,
		&projectStartDate,
		&projectEndDate,
		&project.DefaultBranchName,
		&project.DefaultRemoteName,
	)
	project.ProjectStartDate = projectStartDate.Int64
	project.ProjectEndDate = projectEndDate.Int64

	if err == sql.ErrNoRows {
		project = models.Project{
			ProjectDir: dir,
		}
		id := db.insertProject(project)
		project.ProjectId = id
		project.IsFirstTime = true
		return project
	}
	utils.CheckErr(serviceName, err)
	project.IsFirstTime = false

	return project
}

func (db *DBConfig) insertProject(p models.Project) int64 {
	const serviceName = "insertProject"
	dir := p.ProjectDir
	insertQuery := fmt.Sprintf(
		`INSERT INTO project(project_dir, project_start_date, project_end_date, default_branch_name, default_remote_name) 
			VALUES('%s', %d, %d, 'master', 'origin');`,
		dir,
		p.ProjectStartDate,
		p.ProjectEndDate,
	)
	statement, err := db.Driver.Prepare(insertQuery)
	utils.CheckErr(serviceName, err)

	result, err := statement.Exec()
	utils.CheckErr(serviceName, err)

	id, err := result.LastInsertId()
	utils.CheckErr(serviceName, err)
	return id
}

func (db *DBConfig) UpdateProject(project models.Project) error {
	const serviceName = "UpdateProject"
	query := fmt.Sprintf(
		`UPDATE project 
		SET project_dir='%s', 
			project_start_date=%d,
			project_end_date=%d,
			default_branch_name='%s',
			default_remote_name='%s'
		WHERE project_id=%d;`,
		project.ProjectDir,
		project.ProjectStartDate,
		project.ProjectEndDate,
		project.DefaultBranchName,
		project.DefaultRemoteName,
		project.ProjectId,
	)
	tx, err := db.Driver.Begin()
	if err != nil {
		utils.CheckErr(serviceName, err)
		return err
	}
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		utils.CheckErr(serviceName, err)
		return err
	}
	return tx.Commit()
}
