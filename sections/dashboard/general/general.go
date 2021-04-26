package general

import (
	"fmt"
	"gitmonitor/services/git"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitGeneralTab() fyne.CanvasObject {
	return widget.NewLabel("General Information")
}

func RenderGeneralTab(wrapper fyne.CanvasObject, data GeneralData) {
	// Preprocess data
	authors := git.GetAuthors(data.Commits)
	startDate := data.Commits[len(data.Commits)-1].Author.When
	endDate := data.Commits[0].Author.When
	dayDiff := utils.GetDayDifference(startDate, endDate)
	avgCommits := float32(len(data.Commits)) / float32(dayDiff)

	go getLinesOfCodeInformation(data.FileInformation, []string{data.ProjectDir})

	// create components to render
	projectNameLabel := widget.NewLabel(data.ProjectName)
	startDateLabel := widget.NewLabel(data.RepoStartDate)
	commitsLabel := widget.NewLabel(fmt.Sprintf("%d (%f average per all days)", len(data.Commits), avgCommits))
	filesLabel := widget.NewLabelWithData(data.TotalFiles)
	linesoCodeLabel := widget.NewLabelWithData(data.TotalCode)
	commentLabel := widget.NewLabelWithData(data.TotalComments)
	blankLabel := widget.NewLabelWithData(data.TotalBlanks)
	authorLabel := widget.NewLabel(fmt.Sprintf("%d contributors", len(authors)))
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Project Name", Widget: projectNameLabel},
			{Text: "Project Start Date", Widget: startDateLabel},
			{Text: "Files", Widget: filesLabel},
			{Text: "Total Lines of Code", Widget: linesoCodeLabel},
			{Text: "Total Comment Lines", Widget: commentLabel},
			{Text: "Total Blank Lines", Widget: blankLabel},
			{Text: "Total Commits", Widget: commitsLabel},
			{Text: "Total Authors", Widget: authorLabel},
		},
	}

	projectInfoWrapper := widget.NewCard("Project Information", "", form)
	weeklyChartWrapper := widget.NewCard(
		"Weekly Chart",
		"Last activities in the last 10 weeks. Week-10 is the current week.",
		getWeeklyChartCanvas(data.Commits),
	)
	monthlyChartWrapper := widget.NewCard(
		"Monthly Chart",
		"Total number of commits every day this month.",
		getMonthlyChartCanvas(data.Commits),
	)
	vBox := container.NewVBox(projectInfoWrapper, weeklyChartWrapper, monthlyChartWrapper)

	generalWrapper := wrapper.(*widget.Card)
	generalWrapper.SetContent(container.NewVScroll(vBox))
}
