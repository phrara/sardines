package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sardines/storage"
)

func FilesTab() fyne.CanvasObject {
	c := container.NewWithoutLayout()

	tree := widget.NewTree(childUIDs, isBranch, createNode, updateNode)
	tree.Move(fyne.NewPos(20, 30))
	tree.Resize(fyne.NewSize(750, 500))
	tree.OpenAllBranches()

	c.Add(tree)

	return c
}

func childUIDs(uid string) (c []string) {
	data := storage.FileStoreTree()
	c = data[uid]
	return
}
func isBranch(uid string) (b bool) {
	data := storage.FileStoreTree()
	_, b = data[uid]
	return
}
func createNode(branch bool) fyne.CanvasObject {
	return widget.NewLabel("Template Object")
}
func updateNode(uid string, branch bool, node fyne.CanvasObject) {
	node.(*widget.Label).SetText(uid)
}
