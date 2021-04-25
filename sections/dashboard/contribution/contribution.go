package contribution

import (
	"fmt"
	"gitmonitor/services/git"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type authorTable struct {
	git.Author
	git.AuthorInfo
}

func InitContributionTab() fyne.CanvasObject {
	return widget.NewLabel("Contribution Information")
}

func RenderContributorTab(wrapper fyne.CanvasObject, data ContributorData) {
	authorList := []authorTable{}
	for key, value := range data.authorMap {
		authorList = append(authorList, authorTable{key, value})
	}

	minLabelWidth := []float32{0, 0, 0, 0, 0, 0, 0}

	var table *widget.Table
	table = widget.NewTable(
		func() (int, int) {
			return len(data.authorMap), 7
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(" ")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// set header
				switch id.Col {
				case 0:
					label.SetText("Author Name")
				case 1:
					label.SetText("Author Mail")
				case 2:
					label.SetText("Commits")
				case 3:
					label.SetText("Added lines")
				case 4:
					label.SetText("Deleted lines")
				case 5:
					label.SetText("First commit")
				case 6:
					label.SetText("Last commit")
				default:
					label.SetText(" ")
				}
			} else {
				// set content
				switch id.Col {
				case 0:
					label.SetText(authorList[id.Row-1].Name)
				case 1:
					label.SetText(authorList[id.Row-1].Email)
				case 2:
					label.SetText(fmt.Sprintf("%d", authorList[id.Row-1].TotalCommit))
				case 3:
					label.SetText(fmt.Sprintf("%d", authorList[id.Row-1].TotalAddLines))
				case 4:
					label.SetText(fmt.Sprintf("%d", authorList[id.Row-1].TotalDelLines))
				case 5:
					label.SetText(string(authorList[id.Row-1].FirstCommit.Format("2 Jan 2006 15:04:05")))
				case 6:
					label.SetText(string(authorList[id.Row-1].LastCommit.Format("2 Jan 2006 15:04:05")))
				default:
					label.SetText(" ")
				}
			}
			if minLabelWidth[id.Col] < label.MinSize().Width {
				minLabelWidth[id.Col] = label.MinSize().Width
				go func() {
					table.SetColumnWidth(id.Col, minLabelWidth[id.Col])
					table.Refresh()
				}()
			}
		},
	)

	for i := 0; i <= 6; i++ {
		table.SetColumnWidth(i, 200)
	}
	contributorWrapper := wrapper.(*widget.Card)
	contributorWrapper.SetContent(table)
}
