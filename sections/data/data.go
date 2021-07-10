package data

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
