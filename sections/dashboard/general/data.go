package general

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/git"
	"gitmonitor/services/utils"
	"path"
	"regexp"

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
	ProjectDir    string
	OriginUrl     string
	ProjectName   string
	RepoStartDate string
	Commits       []*object.Commit
}

func getLinesOfCodeInformation(fileInformation FileInformation, paths []string) {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	processor := gocloc.NewProcessor(languages, options)
	result, err := processor.Analyze(paths)
	utils.CheckErr(err)

	fileInformation.TotalFiles.Set(fmt.Sprintf("%d files", len(result.Files)))
	fileInformation.TotalCode.Set(fmt.Sprintf("%d lines", int(result.Total.Code)))
	fileInformation.TotalComments.Set(fmt.Sprintf("%d lines", int(result.Total.Comments)))
	fileInformation.TotalBlanks.Set(fmt.Sprintf("%d lines", int(result.Total.Blanks)))
}

func InitGeneralData(project models.Project, r git.GitConfig) GeneralData {
	var data GeneralData
	data.OriginUrl = r.GetOriginUrl()
	baseName := path.Base(data.OriginUrl)
	re := regexp.MustCompile(`^(.+)\.git$`)
	match := re.FindStringSubmatch(baseName)
	if match != nil {
		data.ProjectName = match[1]
	} else {
		data.ProjectName = baseName
	}

	commits, err := r.GetCommitObjects()
	if err == nil {
		data.Commits = commits
	} else {
		fmt.Println(err)
	}

	if len(commits) > 0 {
		data.RepoStartDate = commits[len(commits)-1].Author.When.Format("2 Jan 2006 15:04:05")
	} else {
		data.RepoStartDate = "No date"
	}

	data.ProjectDir = project.ProjectDir
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

	return data
}
