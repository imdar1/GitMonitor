package main

import (
	"gitmonitor/sections"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	mainApp := app.New()
	window := mainApp.NewWindow("Git Monitor")
	content := sections.GetContent(window)
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
