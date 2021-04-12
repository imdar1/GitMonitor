package state

import (
	"gitmonitor/db"
	"gitmonitor/models"

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
	projectList := db.GetProjects()
	if projectList != nil {
		p.projects = projectList
		projectName := getProjectName(p.projects)
		p.ProjectEntry.SetOptions(projectName)
	}
}

func (p *ProfileState) OnRepositoryLoaded(db *db.DBConfig) models.Project {
	project := db.GetProjectByDir(p.ProjectEntry.Text)
	projectName := getProjectName(p.projects)
	p.ProjectEntry.SetOptions(projectName)
	return project
}
