package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	ON  = true
	OFF = false
)

var (
	CX        float32 = 810
	CY        float32 = 550
	CD        float32 = 27
	w         fyne.Window
	circleOn  *canvas.Circle
	circleOff *canvas.Circle
	hostEntry *widget.Entry
)

func App() {

	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	w = a.NewWindow("sardines")
	w.Resize(fyne.Size{Width: 930, Height: 620})

	c := container.NewWithoutLayout()

	hostEntry = widget.NewEntry()
	hostEntry.Disable()
	hostEntry.SetPlaceHolder("haven't run yet")
	hostInfo := widget.NewForm()
	hostInfo.Resize(fyne.NewSize(750, 30))
	hostInfo.Move(fyne.NewPos(20, CY-5))
	hostInfo.Append("Host", hostEntry)
	initCircle()

	tabs := container.NewAppTabs(
		container.NewTabItem("Panel", PanelTab()),
		container.NewTabItem("Setting", SettingTab()),
		container.NewTabItem("Peers", PeersTab()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.Resize(fyne.Size{
		Width:  890,
		Height: 450,
	})

	c.Add(circleOn)
	c.Add(circleOff)
	c.Add(tabs)
	c.Add(hostInfo)

	w.SetContent(c)
	w.ShowAndRun()
}

func initCircle() {
	circleOn = canvas.NewCircle(color.RGBA{
		R: 69,
		G: 196,
		B: 190,
		A: 1,
	})
	circleOn.StrokeColor = color.Gray{Y: 0x99}
	circleOn.StrokeWidth = 2
	circleOn.Position1 = fyne.NewPos(CX, CY)
	circleOn.Position2 = fyne.NewPos(CX+CD, CY+CD)

	circleOff = canvas.NewCircle(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 1,
	})
	circleOff.StrokeColor = color.Gray{Y: 0x99}
	circleOff.StrokeWidth = 2
	circleOff.Position1 = fyne.NewPos(CX, CY)
	circleOff.Position2 = fyne.NewPos(CX+CD, CY+CD)

	turnOff()
}

func turnOn(host string) {
	circleOn.Show()
	circleOff.Hide()
	hostEntry.SetText(host)
}

func turnOff() {
	circleOff.Show()
	circleOn.Hide()
	hostEntry.SetText("")
}
