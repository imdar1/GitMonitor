package profile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetProfileWindow() fyne.CanvasObject {
	title := widget.NewLabel("Select a Git repository below")
	selectEntry := widget.NewSelectEntry([]string{"Dir A", "Dir B", "Dir C"})
	selectEntry.PlaceHolder = "Type or select project directory"
	exploreButton := widget.NewButton("...", func() {})

	dirExplorer := container.New(layout.NewBorderLayout(nil, nil, nil, exploreButton), exploreButton, selectEntry)
	dirWrapper := container.New(layout.NewVBoxLayout(), title, dirExplorer, widget.NewSeparator())

	return dirWrapper
}
