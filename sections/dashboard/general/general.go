package general

import (
	"fmt"
	"gitmonitor/constants"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getAuthors(commits []*object.Commit) []string {
	set := make(map[string]bool)
	var signatures []string
	for _, v := range commits {
		authorFormat := fmt.Sprintf("%s%s%s", v.Author.Email, constants.Separator, v.Author.Name)
		set[authorFormat] = true
	}

	for k := range set {
		signatures = append(signatures, k)
	}

	return signatures
}

func renderGeneralTab(data GeneralData) {
	// Preprocess data
	authors := getAuthors(data.Commits)
	startDate := data.Commits[len(data.Commits)-1].Author.When
	endDate := data.Commits[0].Author.When
	dayDiff := utils.GetDayDifference(startDate, endDate)
	avgCommits := float32(len(data.Commits)) / float32(dayDiff)

	go getLinesOfCodeInformation(data.fileInformation, []string{data.ProjectDir})

	// create components to render
	projectNameLabel := widget.NewLabel(data.ProjectName)
	startDateLabel := widget.NewLabel(data.ProjectStartDate)
	endDateLabel := widget.NewLabel(data.ProjectEndDate)
	projectStatusLabel := widget.NewLabel(data.ProjectTaskStatus)
	lastCommitLabel := widget.NewLabel(data.Commits[0].Committer.When.Format("2 Jan 2006 15:04:05"))
	commitsLabel := widget.NewLabel(fmt.Sprintf("%d (%f average per all days)", len(data.Commits), avgCommits))
	filesLabel := widget.NewLabelWithData(data.TotalFiles)
	linesOfCodeLabel := widget.NewLabelWithData(data.TotalCode)
	commentLabel := widget.NewLabelWithData(data.TotalComments)
	blankLabel := widget.NewLabelWithData(data.TotalBlanks)
	authorLabel := widget.NewLabel(fmt.Sprintf("%d contributors", len(authors)))
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Project Name", Widget: projectNameLabel},
			{Text: "Project Start Date", Widget: startDateLabel},
			{Text: "Project End Date", Widget: endDateLabel},
			{Text: "Project Status", Widget: projectStatusLabel},
			{Text: "Last Commit", Widget: lastCommitLabel},
			{Text: "Files", Widget: filesLabel},
			{Text: "Total Lines of Code", Widget: linesOfCodeLabel},
			{Text: "Total Comment Lines", Widget: commentLabel},
			{Text: "Total Blank Lines", Widget: blankLabel},
			{Text: "Total Commits", Widget: commitsLabel},
			{Text: "Total Authors", Widget: authorLabel},
		},
	}

	projectInfoWrapper := widget.NewCard("Project Information", "", form)
	weeklyChartWrapper := widget.NewCard(
		"Weekly Chart",
		"Last activities in the last 10 weeks.",
		getWeeklyChartCanvas(data.Commits),
	)
	monthlyChartWrapper := widget.NewCard(
		"Monthly Chart",
		"Total number of commits in the last 30-day.",
		getMonthlyChartCanvas(data.Commits),
	)
	vBox := container.NewVBox(projectInfoWrapper, weeklyChartWrapper, monthlyChartWrapper)

	generalWrapper := data.Wrapper.(*widget.Card)
	generalWrapper.SetContent(container.NewVScroll(vBox))
}
