package state

import (
	"fmt"
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/services"

	"fyne.io/fyne/v2/widget"
)

type ProfileState struct {
	ProjectEntry *widget.SelectEntry
	projects     []models.Project
}

func getProjectName(p []models.Project) []string {
	var projectName []string
	for _, value := range p {
		projectName = append(projectName, value.ProjectDir)
	}
	return projectName
}

func (p *ProfileState) OnDatabaseLoaded(db *db.DBConfig) {
	rows, err := db.Driver.Query("SELECT * FROM project")
	services.CheckErr(err)
	if rows != nil {
		for rows.Next() {
			var project models.Project
			err = rows.Scan(&project.ProjectId, &project.ProjectDir)
			services.CheckErr(err)
			p.projects = append(p.projects, project)
		}
		rows.Close()

		projectName := getProjectName(p.projects)
		p.ProjectEntry.SetOptions(projectName)
	}
}

func (p *ProfileState) OnRepositoryLoaded(repo services.GitConfig, db *db.DBConfig) {
	newDir := p.ProjectEntry.Text
	insertQuery := fmt.Sprintf("INSERT INTO project(project_dir) VALUES('%s');", newDir)
	statement, err := db.Driver.Prepare(insertQuery)
	services.CheckErr(err)
	_, err = statement.Exec()
	services.CheckErr(err)
	projectName := getProjectName(p.projects)
	p.ProjectEntry.SetOptions(projectName)
}
