package observer

import "gitmonitor/config"

type ProfileObserver struct {
	GitRepo config.GitConfig
}

func (po *ProfileObserver) Update() {

}
