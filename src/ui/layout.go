package ui

import (
	"github.com/jroimartin/gocui"
	"sync"
)

type LogHandler func(*gocui.Gui, string)

type App struct {
	screen      *gocui.Gui
	logChannel  chan string
	lock        sync.RWMutex
	logHandlers []LogHandler
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
	go a.handleLogChannel()

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

func (a *App) AddLogHandler(f LogHandler) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.logHandlers = append(a.logHandlers, f)
}

func (a *App) handleLogChannel() {
	for {
		a.runAllHandlers(<-a.logChannel)
	}
}

func (a *App) runAllHandlers(s string) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, handler := range a.logHandlers {
		go handler(a.screen, s)
	}
}
