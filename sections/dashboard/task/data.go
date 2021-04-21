package task

import (
	"gitmonitor/db"
	"gitmonitor/models"
	"gitmonitor/services/git"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type TaskData struct {
	Project  models.Project
	Tasks    []models.Task
	Branches []models.Branch
}

func askAuth() transport.AuthMethod {
	var username string
	var password string
	var authMethod http.BasicAuth
	done := make(chan struct{})

	w := fyne.CurrentApp().NewWindow("Authentication page")

	usernameEntry := widget.NewEntry()
	usernameEntry.OnChanged = func(s string) {
		username = s
	}
	passwordEntry := widget.NewEntry()
	passwordEntry.Password = true
	passwordEntry.OnChanged = func(s string) {
		password = s
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
	}
	okButton := widget.NewButton("OK", func() { w.Close() })
	container := container.NewBorder(
		nil, container.NewHBox(layout.NewSpacer(), okButton), nil, nil, form,
	)
	w.SetContent(container)
	w.SetOnClosed(func() {
		authMethod.Username = username
		authMethod.Password = password
		close(done)
	})
	w.Resize(fyne.NewSize(400, w.Content().Size().Height))
	w.CenterOnScreen()
	w.Show()

	<-done
	return &authMethod
}

func (t *TaskData) ReadTaskData(gitConfig git.GitConfig, db db.DBConfig) {
	// Initialize tasks list
	tasks := db.GetTasksData(t.Project.ProjectId)

	// get remote branches and sync branches with db
	branches, err := gitConfig.GetRemoteBranches(askAuth)
	utils.CheckErr(err)
	err = db.SyncBranches(t.Project.ProjectId, branches)
	utils.CheckErr(err)

	// get branches stored in db
	branchModels, err := db.GetBranchesData(t.Project.ProjectId)
	utils.CheckErr(err)

	// sync task and branch in db
	err = db.SyncTask(tasks, branchModels)
	utils.CheckErr(err)

	// get updated tasks
	tasks = db.GetTasksData(t.Project.ProjectId)

	t.Tasks = tasks
	t.Branches = branchModels
}
