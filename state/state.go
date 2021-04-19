package state

import (
	"gitmonitor/db"
	"gitmonitor/services/git"
)

type AppData struct {
	Repo     git.GitConfig
	Database *db.DBConfig
}

type AppState struct {
	AppData
	ProfileState   ProfileState
	DashboardState TabItemsState
}

func (a *AppState) OnDatabaseLoaded() {
	a.ProfileState.OnDatabaseLoaded(a.Database)
	a.DashboardState.OnDatabaseLoaded(a.Database)
}

func (a *AppState) OnRepositoryLoaded() {
	selectedProject := a.ProfileState.OnRepositoryLoaded(a.AppData)
	a.DashboardState.OnRepositoryLoaded(a.AppData, selectedProject)
}
