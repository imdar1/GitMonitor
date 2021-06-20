package state

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/services/git"
)

type AppData struct {
	Repo            git.GitConfig
	Database        *db.DBConfig
	SelectedProject models.Project
}

type AppState struct {
	AppData
	ProfileState   ProfileState
	DashboardState TabItemsState
}

func (a *AppState) OnWindowLoaded() {
	a.ProfileState.OnWindowLoaded(&a.AppData)
	a.DashboardState.OnWindowLoaded(&a.AppData)
}

func (a *AppState) OnRepositoryLoaded() {
	a.ProfileState.OnRepositoryLoaded(&a.AppData)
	a.DashboardState.OnRepositoryLoaded(&a.AppData)
}
