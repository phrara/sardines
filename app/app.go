package app

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func App() {
	a := app.New()
	w := a.NewWindow("sardines")

	tabs := container.NewAppTabs(
		container.NewTabItem("Home", HomeTab()),
		container.NewTabItem("Setting", widget.NewLabel("World!")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.Size{Width: 720,Height: 480})
	w.ShowAndRun()
}