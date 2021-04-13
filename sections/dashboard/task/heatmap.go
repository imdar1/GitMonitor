package task

import (
	"gitmonitor/services"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type heatmap struct {
	colors  []string
	days    []string
	month   string
	tasks   []string
	authors map[string]string
}

func initData(taskData TaskData) heatmap {
	firstDate, lastDate := services.GetFirstAndLastDayOfMonth()
	_, month, firstDay := firstDate.Date()
	_, _, lastDay := lastDate.Date()
	var days []string
	for i := firstDay; i <= lastDay; i++ {
		days = append(days, strconv.Itoa(i))
	}

	// var tasks []string
	// authors := make(map[string]string)
	// for _, value := range taskData.Tasks {
	// 	tasks = append(tasks, value.Name)
	// 	authors
	// }

	h := heatmap{
		days:  days,
		month: month.String(),
		tasks: []string{
			"abc", "def", "ghi",
		},
		colors: []string{"#50a3ba", "#eac736", "#d94e5d"},
	}
	return h

	// TODO: using models.Task, extract all information
}

func (h *heatmap) getTable() fyne.CanvasObject {
	t := widget.NewTable(
		func() (int, int) {
			rows, cols := len(h.tasks)+1, len(h.days)+1
			return rows, cols
		},
		func() fyne.CanvasObject {
			// return container.NewCenter(
			// 	container.NewWithoutLayout(canvas.NewRectangle(color.NRGBA{0x80, 0, 0, 0xff})),
			// )
			return widget.NewLabel("cek")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// outer := cell.(*fyne.Container)
			// container := outer.Objects[0].(*fyne.Container)
			label := cell.(*widget.Label)
			if id.Col == 0 && id.Row >= 1 {
				label.SetText(h.tasks[id.Row-1])
				// container.Add(widget.NewLabel(h.tasks[id.Row-1]))
			}
			if id.Row == 0 && id.Col >= 1 {
				label.SetText(h.days[id.Col-1])
				// container.Add(widget.NewLabel(h.days[id.Col-1]))
			}
		},
	)
	t.SetColumnWidth(0, 102)
	return t
}
