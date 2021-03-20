package forms

import "fyne.io/fyne/v2"

type StateWindow struct {
	Title string
	View  func(a fyne.App) fyne.CanvasObject
}

var (
	StateMap = map[string]StateWindow{
		"profile":   {"Select Profile", nil},
		"dashboard": {"Git Monitor", nil},
	}
)
