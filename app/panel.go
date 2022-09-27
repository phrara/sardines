package app

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sardines/core"
	"sardines/err"
)

var (
	ctx    context.Context    = nil
	cancel context.CancelFunc = nil
	hNode  *core.HostNode
)

func PanelTab() fyne.CanvasObject {
	c := container.New(layout.NewVBoxLayout())

	lab := widget.NewLabel("Operation")
	lab.Alignment = fyne.TextAlignCenter

	btnOn := widget.NewButton("RUN", func() {
		if cancel != nil {
			return
		}
		ctx, cancel = context.WithCancel(context.Background())
		msg := make(chan bool, 1)
		go Run(ctx, msg)
		if b := <-msg; b {
			turnOn(hNode.NodeAddr.String())
		}
	})
	btnOn.Alignment = widget.ButtonAlignCenter

	btnOff := widget.NewButton("SHUT", func() {
		if cancel != nil {
			Ensure("are you sure?", func(b bool) {
				if b {
					cancel()
					turnOff()
					cancel = nil
				}
			})
		}
	})
	btnOff.Alignment = widget.ButtonAlignCenter

	c.Add(lab)
	c.Add(btnOn)
	c.Add(btnOff)

	return c
}

func Run(ctx context.Context, msg chan<- bool) {
	h, er := core.GenerateNode()
	if er != nil {
		ShowErr(er)
		hNode = nil
		msg <- false
		return
	}
	hNode = h

	if f := <-h.JoinNetwork(); f == 0 {
		ShowErr(err.ErrJoinNetwork)
		hNode = nil
		msg <- false
		return
	}
	msg <- true

	<-ctx.Done()
	h.Close()

}
