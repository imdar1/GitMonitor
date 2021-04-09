package state

import (
	"gitmonitor/db"
	"gitmonitor/services"
)

type State interface {
	OnDatabaseLoaded()
	OnRepositoryLoaded()
}

type AppData struct {
	Repo     services.GitConfig
	Database db.DBConfig
}

type AppState struct {
	AppData
	ProfileState   ProfileState
	DashboardState TabItemsState
}

func (a *AppState) OnDatabaseLoaded() {
	a.ProfileState.OnDatabaseLoaded()
	a.DashboardState.OnDatabaseLoaded()
}

func (a *AppState) OnRepositoryLoaded() {
	a.ProfileState.OnRepositoryLoaded()
	a.DashboardState.OnRepositoryLoaded()
}
