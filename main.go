package main

import (
	"gitmonitor/sections"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	mainApp := app.New()
	content := sections.GetContent()
	window := mainApp.NewWindow("Git Monitor")
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
