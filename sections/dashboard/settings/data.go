package settings

import (
	"gitmonitor/models"
	"gitmonitor/sections/data"

	"fyne.io/fyne/v2"
)

type SettingsData struct {
	RemoteBranches      []models.Branch
	AdditionalRenderers []data.Renderer
	Wrapper             fyne.CanvasObject
}

func InitSettingsData(wrapper fyne.CanvasObject, remoteBranches []models.Branch) SettingsData {
	return SettingsData{
		RemoteBranches: remoteBranches,
		Wrapper:        wrapper,
	}
}

func (s SettingsData) Render(data *data.AppData) {
	renderSettingsTab(s, data)
}
