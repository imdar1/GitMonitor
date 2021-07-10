package state

import "gitmonitor/sections/data"

type AppState struct {
	data.AppData
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
