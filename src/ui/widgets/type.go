package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type KeyHandler struct {
	Key     gocui.Key
	Handler func(*gocui.Gui, *gocui.View) error
}
type Handlers []KeyHandler
type PointCalculator func(int) int

type Widget struct {
	Name       string
	Title      string
	Editable   bool
	Autoscroll bool
	Wrap       bool
	x0, y0     PointCalculator
	x1, y1     PointCalculator
	data       []string
	handlers   Handlers

	// events
	OnValueChange chan string
}

func (w *Widget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(w.Name, w.x0(maxX), w.y0(maxY), w.x1(maxX), w.y1(maxY)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = w.Title
		v.Editable = w.Editable
		v.Autoscroll = w.Autoscroll
		v.Wrap = w.Wrap

		if len(w.data) > 0 {
			v.Clear()
			for _, row := range w.data {
				if _, err := fmt.Fprintln(v, row); err != nil {
					return err
				}
			}
		}

		for _, keyHandler := range w.handlers {
			if err := g.SetKeybinding(w.Name, keyHandler.Key, gocui.ModNone, keyHandler.Handler); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Widget) AddHandler(handler KeyHandler) {
	w.handlers = append(w.handlers, handler)
}
