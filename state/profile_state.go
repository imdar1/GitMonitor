package state

import (
	"gitmonitor/models"
	"gitmonitor/sections/data"

	"fyne.io/fyne/v2/widget"
)

type ProfileState struct {
	ProjectEntry *widget.SelectEntry
	projects     []string
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
		p.projects = getProjectName(projectList)
		p.ProjectEntry.SetOptions(p.projects)
	}
}

func (p *ProfileState) OnRepositoryLoaded(appData *data.AppData) {
	project := appData.Database.GetProjectByDir(p.ProjectEntry.Text)
	projectList := appData.Database.GetProjects()
	p.projects = getProjectName(projectList)
	p.ProjectEntry.SetOptions(p.projects)
	appData.SelectedProject = project
}
