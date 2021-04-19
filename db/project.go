package db

import (
	"database/sql"
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/utils"
)

func (db *DBConfig) GetProjects() []models.Project {
	var projects []models.Project
	rows, err := db.Driver.Query("SELECT * FROM project")
	utils.CheckErr(err)
	if rows != nil {
		for rows.Next() {
			var project models.Project
			err = rows.Scan(&project.ProjectId, &project.ProjectDir)
			utils.CheckErr(err)
			projects = append(projects, project)
		}
		rows.Close()
	}
	return projects
}

func (db *DBConfig) GetProjectByDir(dir string) models.Project {
	var project models.Project

	query := fmt.Sprintf("SELECT * FROM project WHERE project_dir='%s' LIMIT 1;", dir)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(&project.ProjectId, &project.ProjectDir)
	if err == sql.ErrNoRows {
		project = models.Project{
			ProjectDir: dir,
		}
		id := db.insertProject(project)
		project.ProjectId = id
		return project
	}
	utils.CheckErr(err)

	return project
}

func (db *DBConfig) insertProject(p models.Project) int64 {
	dir := p.ProjectDir
	insertQuery := fmt.Sprintf("INSERT INTO project(project_dir) VALUES('%s');", dir)
	statement, err := db.Driver.Prepare(insertQuery)
	utils.CheckErr(err)

	result, err := statement.Exec()
	utils.CheckErr(err)

	id, err := result.LastInsertId()
	utils.CheckErr(err)
	return id
}
