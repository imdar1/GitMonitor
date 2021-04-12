package db

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services"
)

func (db *DBConfig) GetProjects() []models.Project {
	var projects []models.Project
	rows, err := db.Driver.Query("SELECT * FROM project")
	services.CheckErr(err)
	if rows != nil {
		for rows.Next() {
			var project models.Project
			err = rows.Scan(&project.ProjectId, &project.ProjectDir)
			services.CheckErr(err)
			projects = append(projects, project)
		}
		rows.Close()
	}
	return projects
}

func (db *DBConfig) GetProjectByDir(dir string) models.Project {
	query := fmt.Sprintf("SELECT * FROM project WHERE project_dir='%s' LIMIT 1;", dir)
	rows, err := db.Driver.Query(query)
	services.CheckErr(err)

	var project models.Project
	if rows == nil {
		project = models.Project{
			ProjectDir: dir,
		}
		id := db.insertProject(project)
		project.ProjectId = id
		return project
	}

	err = rows.Scan(&project.ProjectId, &project.ProjectDir)
	services.CheckErr(err)
	return project
}

func (db *DBConfig) insertProject(p models.Project) int64 {
	dir := p.ProjectDir
	insertQuery := fmt.Sprintf("INSERT INTO project(project_dir) VALUES('%s');", dir)
	statement, err := db.Driver.Prepare(insertQuery)
	services.CheckErr(err)

	result, err := statement.Exec()
	services.CheckErr(err)

	id, err := result.LastInsertId()
	services.CheckErr(err)
	return id
}
