package state

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/sections/dashboard/contribution"
	"gitmonitor/sections/dashboard/general"
	"gitmonitor/sections/dashboard/task"
	"gitmonitor/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TabItemsState struct {
	GeneralContent      fyne.CanvasObject
	TaskContent         fyne.CanvasObject
	ContributionContent fyne.CanvasObject
}

func (tabState *TabItemsState) OnDatabaseLoaded(db *db.DBConfig) {

}

func (tabState *TabItemsState) OnRepositoryLoaded(repo services.GitConfig, db *db.DBConfig, project models.Project) {
	taskContent := tabState.TaskContent.(*widget.Card)

	tasks := db.GetTasksData(project.ProjectId)
	branches := db.GetBranchesData(project.ProjectId)

	taskData := task.TaskData{
		Tasks:    tasks,
		Branches: branches,
	}
	taskContent.SetContent(task.RenderTaskTab(taskData))
	taskContent.Refresh()

}

func InitTabItems() TabItemsState {

	tabItems := TabItemsState{
		GeneralContent:      widget.NewCard("", "", general.InitGeneralTab()),
		TaskContent:         widget.NewCard("", "", task.InitTaskTab()),
		ContributionContent: widget.NewCard("", "", contribution.InitContributionTab()),
	}
	return tabItems
}
