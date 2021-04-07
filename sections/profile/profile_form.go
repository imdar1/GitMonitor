package profile

import (
	"gitmonitor/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetProfileWindow(w fyne.Window, gitConfig *services.GitConfig) fyne.CanvasObject {
	title := widget.NewLabel("Select a Git repository below")
	selectEntry := widget.NewSelectEntry([]string{"Dir A", "Dir B", "Dir C"})
	selectEntry.PlaceHolder = "Type or select project directory"
	loadButton := widget.NewButton("Load", func() {
		_, err := services.InitGit(selectEntry.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})
	loadButton.Disable()
	exploreButton := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(
			func(uri fyne.ListableURI, err error) {
				if err == nil {
					selectEntry.SetText(uri.Path())
					loadButton.Enable()
				}
			},
			w,
		)
	})
	buttonWrapper := container.NewHBox(exploreButton, loadButton)

	dirExplorer := container.New(
		layout.NewBorderLayout(
			nil,
			nil,
			nil,
			buttonWrapper,
		),
		buttonWrapper,
		selectEntry,
	)
	dirWrapper := container.New(layout.NewVBoxLayout(), title, dirExplorer, widget.NewSeparator())

	return dirWrapper
}
