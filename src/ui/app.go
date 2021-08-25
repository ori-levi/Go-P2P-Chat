package ui

import "github.com/jroimartin/gocui"

type LogConsumer func(*gocui.Gui, interface{})

type App struct {
	screen      *gocui.Gui
	logHandlers *Handlers
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
		screen:      g,
		logHandlers: NewHandlers(),
	}, nil
}

func (a *App) Close() {
	a.screen.Close()
}

func (a *App) Run(managers ...gocui.Manager) error {
	g := a.screen

	g.SetManager(managers...)

	a.logHandlers.Start()

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

func (a *App) AddLogConsumer(f LogConsumer) {
	a.logHandlers.AddConsumer(func(value interface{}) {
		f(a.screen, value)
	})
}

// todo fix this!
//func (a *App) AddLogProvider(f func(string, ...interface{})) {
//	a.logHandlers.AddProvider(func() interface{} {
//		return f()
//	})
//}
