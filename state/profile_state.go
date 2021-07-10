package state

import (
	"gitmonitor/models"
	"gitmonitor/sections/data"

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

func (p *ProfileState) OnWindowLoaded(appData *data.AppData) {
	projectList := appData.Database.GetProjects()
	if projectList != nil {
		p.projects = projectList
		projectName := getProjectName(p.projects)
		p.ProjectEntry.SetOptions(projectName)
	}
}

func (p *ProfileState) OnRepositoryLoaded(appData *data.AppData) {
	project := appData.Database.GetProjectByDir(p.ProjectEntry.Text)
	projectName := getProjectName(p.projects)
	p.ProjectEntry.SetOptions(projectName)
	appData.SelectedProject = project
}
