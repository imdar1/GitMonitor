package state

import (
	"gitmonitor/db"
	"gitmonitor/services"
)

type AppData struct {
	Repo     services.GitConfig
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
	selectedProject := a.ProfileState.OnRepositoryLoaded(a.Database)
	a.DashboardState.OnRepositoryLoaded(a.Repo, a.Database, selectedProject)
}
