package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowErr(e error) {
	dialog.ShowError(e, w)
}

func ShowInfo(info string) {
	dialog.ShowInformation("info", info, w)
}

func Ensure(info string, f func(b bool)) {
	dialog.ShowConfirm("ensure", info, f, w)
}

func ShowData(info, data string) {
	e := widget.NewEntry()
	e.SetText(data)
	e.Disable()

	custom := dialog.NewCustom(info, "关闭", e, w)
	custom.Resize(fyne.NewSize(700, 50))
	custom.Show()
}
