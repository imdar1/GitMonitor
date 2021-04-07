package observer

import (
	"gitmonitor/config"

	"fyne.io/fyne/v2/container"
)

type DashboardObserver struct {
	GitRepo        config.GitConfig
	DashboardItems *container.TabItem
}

func (do *DashboardObserver) Update() {

}
