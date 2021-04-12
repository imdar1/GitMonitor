package state

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type TabItemsState struct {
	GeneralContent      *fyne.Container
	TaskContent         *fyne.Container
	ContributionContent *fyne.Container
}

func (tabState *TabItemsState) OnDatabaseLoaded(db *db.DBConfig) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(repo services.GitConfig, db *db.DBConfig, project models.Project) {
	tabState.GeneralContent.Objects = nil
	tabState.TaskContent.Objects = nil

	tasks := db.GetTasksData(project.ProjectId)
	branches := db.GetBranchesData(project.ProjectId)

	TaskData := task.TaskData{
		Tasks:    tasks,
		Branches: branches,
	}
	tabState.GeneralContent.Add(task.RenderTaskTab(TaskData))
	tabState.GeneralContent.Refresh()

}

func InitTabItems() TabItemsState {
	tabItems := TabItemsState{
		GeneralContent:      container.NewVBox(general.InitGeneralTab()),
		TaskContent:         container.NewVBox(task.InitTaskTab()),
		ContributionContent: container.NewVBox(contribution.InitContributionTab()),
	}
	return tabItems
}
