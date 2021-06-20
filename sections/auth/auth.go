package auth

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func AskAuth() transport.AuthMethod {
	var username string
	var password string
	var authMethod http.BasicAuth
	done := make(chan struct{})

	w := fyne.CurrentApp().NewWindow("Authentication page")

	usernameEntry := widget.NewEntry()
	usernameEntry.OnChanged = func(s string) {
		username = s
	}
	passwordEntry := widget.NewEntry()
	passwordEntry.Password = true
	passwordEntry.OnChanged = func(s string) {
		password = s
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
	}
	okButton := widget.NewButton("OK", func() { w.Close() })
	container := container.NewBorder(
		nil, container.NewHBox(layout.NewSpacer(), okButton), nil, nil, form,
	)
	w.SetContent(container)
	w.SetOnClosed(func() {
		authMethod.Username = username
		authMethod.Password = password
		close(done)
	})
	w.Resize(fyne.NewSize(400, w.Content().Size().Height))
	w.CenterOnScreen()
	w.Show()

	<-done
	return &authMethod
}
