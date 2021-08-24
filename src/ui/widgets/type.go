package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Widget struct {
	Name       string
	Title      string
	Editable   bool
	Autoscroll bool
	Wrap       bool
	x0, y0     func(int) int
	x1, y1     func(int) int
	data       []string
	keys       []gocui.Key
	handler    func(g *gocui.Gui, v *gocui.View) error
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

		if w.handler != nil {
			for _, key := range w.keys {
				if err := g.SetKeybinding(w.Name, key, gocui.ModNone, w.handler); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
