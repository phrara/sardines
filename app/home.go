package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func HomeTab() fyne.CanvasObject {
	c := container.New(layout.NewVBoxLayout())

	btn := widget.NewButton("RUN", func() {
		println("sss")
	})
	btn.Resize(fyne.Size{
		Width:  50,
		Height: 150,
	})
	btn.Alignment = widget.ButtonAlignCenter

	lab := widget.NewLabel("Sardines")
	lab.Alignment = fyne.TextAlignCenter
	c.Add(lab)
	c.Add(btn)

	return c
}
