package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var titleStyle = fyne.TextStyle{
	Bold:   true,
	Italic: true,
}

func newUi(input chan string) *fyne.Window {
	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(1024, 768))

	chatPanel, _ := newChatPanel()
	usersPanel, _ := newUsersPanel()
	notificationPanel, _ := newNotificationPanel()

	w.SetContent(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				nil,
				container.NewGridWithRows(2,
					usersPanel,
					newHelpPanel(),
				),
				chatPanel,
			),
			newInputPanel(input),
			notificationPanel,
		),
	)

	return &w
}

func newInputPanel(input chan string) *fyne.Container {
	inputEntry := widget.NewEntry()
	inputEntry.PlaceHolder = "What's On Your Mind?"

	return container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("Send", func() {
			if len(inputEntry.Text) == 0 {
				return
			}

			input <- inputEntry.Text
			inputEntry.SetText("")
		}),
		inputEntry,
	)
}

func newChatPanel() (*widget.Label, *string) {
	chat := "chat"
	l := widget.NewLabelWithData(binding.BindString(&chat))
	return l, &chat
}

func newUsersPanel() (*fyne.Container, *string) {
	users := "users"

	w := container.NewVBox(
		widget.NewLabelWithStyle("Users", fyne.TextAlignLeading, titleStyle),
		widget.NewLabelWithData(binding.BindString(&users)),
	)

	return w, &users
}

func newHelpPanel() *fyne.Container {
	help := []string{
		fmt.Sprintf("%-9v <ip> <port>", "/connect"),
		fmt.Sprintf("%-10v <name> <message...>", "/pm"),
		fmt.Sprintf("%-10v <name> <command...>", "/shell"),
		"/exit",
	}

	return container.NewVBox(
		widget.NewLabelWithStyle("Help", fyne.TextAlignLeading, titleStyle),
		widget.NewLabel(strings.Join(help, "\n")),
	)
}

func newNotificationPanel() (*fyne.Container, *string) {
	log := "logging"
	w := container.NewVBox(
		widget.NewLabelWithStyle("Notifications", fyne.TextAlignLeading, titleStyle),
		widget.NewLabelWithData(binding.BindString(&log)),
	)

	return w, &log
}
