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
	rows, err := db.Driver.Query("SELECT * FROM project")
	utils.CheckErr(serviceName, err)
	if rows != nil {
		for rows.Next() {
			var project models.Project
			err = rows.Scan(
				&project.ProjectId,
				&project.ProjectDir,
				&project.DefaultBranchName,
				&project.DefaultRemoteName,
			)
			utils.CheckErr(serviceName, err)
			projects = append(projects, project)
		}
		rows.Close()
	}
	return projects
}

func (db *DBConfig) GetProjectByDir(dir string) models.Project {
	var project models.Project
	const serviceName = "GetProjectByDir"

	query := fmt.Sprintf("SELECT * FROM project WHERE project_dir='%s' LIMIT 1;", dir)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(
		&project.ProjectId,
		&project.ProjectDir,
		&project.DefaultBranchName,
		&project.DefaultRemoteName,
	)
	if err == sql.ErrNoRows {
		project = models.Project{
			ProjectDir: dir,
		}
		id := db.insertProject(project)
		project.ProjectId = id
		return project
	}
	utils.CheckErr(serviceName, err)

	return project
}

func (db *DBConfig) insertProject(p models.Project) int64 {
	const serviceName = "insertProject"
	dir := p.ProjectDir
	insertQuery := fmt.Sprintf(
		`INSERT INTO project(project_dir, default_branch_name, default_remote_name) 
			VALUES('%s', 'master', 'origin');`,
		dir,
	)
	statement, err := db.Driver.Prepare(insertQuery)
	utils.CheckErr(serviceName, err)

	result, err := statement.Exec()
	utils.CheckErr(serviceName, err)

	id, err := result.LastInsertId()
	utils.CheckErr(serviceName, err)
	return id
}
