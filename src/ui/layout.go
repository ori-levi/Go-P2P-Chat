package ui

import (
	"github.com/jroimartin/gocui"
)

type App struct {
	screen     *gocui.Gui
	logChannel chan string
}

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func NewApp() (*App, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g.Highlight = true
	g.Cursor = true
	g.InputEsc = true
	g.SelFgColor = gocui.ColorGreen

	return &App{
		screen:     g,
		logChannel: make(chan string),
	}, nil
}

func (a *App) Close() {
	a.screen.Close()
}

func (a *App) Run(managers ...gocui.Manager) error {
	g := a.screen

	g.SetManager(managers...)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
