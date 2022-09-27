package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func PeersTab() fyne.CanvasObject {
	c := container.NewWithoutLayout()

	lab := widget.NewLabel("Peer List")
	lab.Alignment = fyne.TextAlignCenter
	lab.Move(fyne.NewPos(375, 0))

	list := widget.NewList(
		length,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		updateItem)
	list.Move(fyne.NewPos(20, 30))
	list.Resize(fyne.NewSize(750, 500))

	btnUP := widget.NewButton(`^`, func() {
		list.ScrollToTop()
	})
	btnUP.Alignment = widget.ButtonAlignCenter
	btnUP.Move(fyne.NewPos(800, 480))
	btnUP.Resize(fyne.NewSize(30, 20))

	c.Add(lab)
	c.Add(list)
	c.Add(btnUP)

	return c
}

func length() int {
	if hNode != nil {
		return hNode.Router.Sum()
	} else {
		return 10
	}
}

func updateItem(i widget.ListItemID, o fyne.CanvasObject) {
	str := ""
	if hNode != nil {
		nodes := hNode.Router.AllNodes()
		if nodes[i] != nil {
			str = fmt.Sprintf("%d: %s", i, nodes[i].String())
		}
	}
	o.(*widget.Label).SetText(str)
}
