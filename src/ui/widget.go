package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Handler func(*gocui.Gui, *gocui.View) error

type KeyHandler struct {
	gocui.Key
	Handler
}
type Handlers []KeyHandler
type PointCalculator func(int) int

type Widget struct {
	Name          string
	title         string
	editable      bool
	autoscroll    bool
	wrap          bool
	x0, y0        PointCalculator
	x1, y1        PointCalculator
	data          []string
	handlers      Handlers
	isCurrentView bool

	// events
	OnValueChange chan string
}

func (w *Widget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(w.Name, w.x0(maxX), w.y0(maxY), w.x1(maxX), w.y1(maxY)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = w.title
		v.Editable = w.editable
		v.Autoscroll = w.autoscroll
		v.Wrap = w.wrap

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

	if w.isCurrentView {
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	return nil
}

func (w *Widget) AddHandler(handler KeyHandler) {
	w.handlers = append(w.handlers, handler)
}

func NewWidget(
	name string,
	title string,
	autoscroll bool,
	editable bool,
	wrap bool,
	x0, y0 PointCalculator,
	x1, y1 PointCalculator,
	isCurrentView bool,
	data []string,
) *Widget {
	return &Widget{
		Name:          name,
		title:         title,
		autoscroll:    autoscroll,
		editable:      editable,
		wrap:          wrap,
		x0:            x0,
		y0:            y0,
		x1:            x1,
		y1:            y1,
		data:          data,
		isCurrentView: isCurrentView,
		OnValueChange: make(chan string),
	}
}
