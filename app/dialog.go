package app

import "fyne.io/fyne/v2/dialog"

func ShowErr(e error) {
	dialog.ShowError(e, w)
}

func ShowInfo(info string) {
	dialog.ShowInformation("info", info, w)
}

func Ensure(info string, f func(b bool)) {
	dialog.ShowConfirm("ensure", info, f, w)
}
