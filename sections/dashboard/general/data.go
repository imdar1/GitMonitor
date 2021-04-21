package general

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/git"
	"path"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hhatto/gocloc"
)

type FileInformation struct {
	TotalFiles    int
	TotalCode     int
	TotalComments int
	TotalBlanks   int
}

type GeneralData struct {
	FileInformation
	OriginUrl     string
	ProjectName   string
	RepoStartDate string
	Commits       []*object.Commit
}

func getLinesOfCodeInformation(paths []string) (FileInformation, error) {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	processor := gocloc.NewProcessor(languages, options)
	result, err := processor.Analyze(paths)
	if err != nil {
		return FileInformation{}, err
	}

	f := FileInformation{
		TotalFiles:    len(result.Files),
		TotalCode:     int(result.Total.Code),
		TotalComments: int(result.Total.Comments),
		TotalBlanks:   int(result.Total.Blanks),
	}
	return f, nil
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

	f, err := getLinesOfCodeInformation([]string{project.ProjectDir})
	if err == nil {
		data.FileInformation = f
	} else {
		fmt.Println(err)
	}

	return data
}
