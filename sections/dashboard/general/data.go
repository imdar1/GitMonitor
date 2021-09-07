package general

import (
	"fmt"
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"
	"path"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hhatto/gocloc"
)

type FileInformation struct {
	TotalFiles    binding.String
	TotalCode     binding.String
	TotalComments binding.String
	TotalBlanks   binding.String
}

type GeneralData struct {
	FileInformation
	ProjectDir        string
	OriginUrl         string
	ProjectName       string
	ProjectStartDate  string
	ProjectEndDate    string
	ProjectTaskStatus string
	Commits           []*object.Commit
	Wrapper           fyne.CanvasObject

	tasks []models.Task
}

func (data GeneralData) Render(appData *data.AppData) {
	renderGeneralTab(data)
}

func (data *GeneralData) UpdateProjectTaskStatus(appData *data.AppData) {
	data.ProjectTaskStatus = "Project schedule has not been set"

	if appData.SelectedProject.ProjectEndDate > 0 {
		mapTaskStatusToCounter := map[int]int{}
		taskCounter := 0
		for _, task := range data.tasks {
			mapTaskStatusToCounter[task.TaskStatus]++
			taskCounter++
		}

		if time.Now().Unix() > appData.SelectedProject.ProjectEndDate {
			data.ProjectTaskStatus = "Late from the planned project end time"
			waitingAndInProgressCount :=
				mapTaskStatusToCounter[int(constants.Waiting)] + mapTaskStatusToCounter[int(constants.InProgress)]
			if waitingAndInProgressCount > 0 {
				data.ProjectTaskStatus = fmt.Sprintf(
					"%s (%d remaining tasks)",
					data.ProjectTaskStatus,
					waitingAndInProgressCount,
				)
			} else {
				data.ProjectTaskStatus = "Project has finished"
			}
		} else if mapTaskStatusToCounter[int(constants.DoneLate)] > 0 {
			data.ProjectTaskStatus = fmt.Sprintf(
				"%d of %d tasks finished late",
				mapTaskStatusToCounter[int(constants.DoneLate)],
				taskCounter,
			)
		} else {
			data.ProjectTaskStatus = "On-track"
		}
	}
}

func (data *GeneralData) UpdateProjectStartDate(appData *data.AppData) {
	if appData.SelectedProject.ProjectStartDate != 0 {
		data.ProjectStartDate = time.Unix(appData.SelectedProject.ProjectStartDate, 0).Format("2 Jan 2006 ")
	} else {
		data.ProjectStartDate = "No date"
	}
}

func (data *GeneralData) UpdateProjectEndDate(appData *data.AppData) {
	if appData.SelectedProject.ProjectEndDate != 0 {
		data.ProjectEndDate = time.Unix(appData.SelectedProject.ProjectEndDate, 0).Format("2 Jan 2006")
	} else {
		data.ProjectEndDate = "No date"
	}
}

func getLinesOfCodeInformation(fileInformation FileInformation, paths []string) {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	processor := gocloc.NewProcessor(languages, options)
	result, err := processor.Analyze(paths)
	utils.CheckErr("getLinesOfCodeInformation", err)

	fileInformation.TotalFiles.Set(fmt.Sprintf("%d files", len(result.Files)))
	fileInformation.TotalCode.Set(fmt.Sprintf("%d lines", int(result.Total.Code)))
	fileInformation.TotalComments.Set(fmt.Sprintf("%d lines", int(result.Total.Comments)))
	fileInformation.TotalBlanks.Set(fmt.Sprintf("%d lines", int(result.Total.Blanks)))
}

func InitGeneralData(wrapper fyne.CanvasObject, tasks []models.Task, appData *data.AppData) GeneralData {
	var data GeneralData
	data.OriginUrl = appData.Repo.GetOriginUrl()
	baseName := path.Base(data.OriginUrl)
	re := regexp.MustCompile(`^(.+)\.git$`)
	match := re.FindStringSubmatch(baseName)
	if match != nil {
		data.ProjectName = match[1]
	} else {
		data.ProjectName = baseName
	}

	commits, err := appData.Repo.GetCommitObjects(
		appData.SelectedProject.DefaultRemoteName,
		appData.SelectedProject.DefaultBranchName,
	)
	if err == nil {
		data.Commits = commits
	} else {
		utils.CheckErr("InitGeneralData", err)
	}

	data.tasks = tasks
	data.UpdateProjectStartDate(appData)
	data.UpdateProjectEndDate(appData)
	data.UpdateProjectTaskStatus(appData)

	data.ProjectDir = appData.SelectedProject.ProjectDir
	data.FileInformation = FileInformation{
		TotalFiles:    binding.NewString(),
		TotalCode:     binding.NewString(),
		TotalComments: binding.NewString(),
		TotalBlanks:   binding.NewString(),
	}

	data.TotalFiles.Set("Analyzing...")
	data.TotalCode.Set("Analyzing...")
	data.TotalComments.Set("Analyzing...")
	data.TotalBlanks.Set("Analyzing...")
	data.Wrapper = wrapper

	return data
}
