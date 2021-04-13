package main

import (
	"gitmonitor/db"
	"gitmonitor/sections"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	dbConfig, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer dbConfig.Close()

	mainApp := app.New()
	window := mainApp.NewWindow("Git Monitor")
	content := sections.GetContent(window, &dbConfig)
	window.SetContent(content)
	window.Resize(fyne.NewSize(1024, 720))
	window.ShowAndRun()
}
