package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sardines/tool"
	"strconv"
)

func FileTab(index int, file *tool.File) *container.TabItem {

	c := container.NewWithoutLayout()
	c.Resize(fyne.NewSize(890, 590))

	originEntry := widget.NewEntry()
	originEntry.SetText(file.Origin)
	originEntry.Disable()
	originEntry.Resize(fyne.NewSize(890, btnH))
	originEntry.Move(fyne.NewPos(0, 0))

	cidEntry := widget.NewEntry()
	cidEntry.SetText(file.CID)
	cidEntry.Disable()
	cidEntry.Resize(fyne.NewSize(890, btnH))
	cidEntry.Move(fyne.NewPos(0, 50))

	var contentEntry *widget.Entry
	if file.IsText() {
		contentEntry = widget.NewEntry()
		contentEntry.SetText(string(file.Content))
		contentEntry.Disable()
	} else {
		contentEntry = widget.NewEntry()
		contentEntry.SetText("该文件不可预览")
		contentEntry.Disable()
	}
	contentEntry.Resize(fyne.NewSize(890, 450))
	contentEntry.Move(fyne.NewPos(0, 100))

	c.Add(originEntry)
	c.Add(cidEntry)
	c.Add(contentEntry)

	return container.NewTabItem("file"+strconv.Itoa(index), c)
}
