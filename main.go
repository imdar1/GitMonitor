package main

import (
	"gitmonitor/db"
	"gitmonitor/sections"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}

	mainApp := app.New()
	window := mainApp.NewWindow("Git Monitor")
	content := sections.GetContent(window, db)
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
