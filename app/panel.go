package app

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sardines/core"
	"sardines/err"
	"sardines/storage"
	"sardines/tool"
)

var (
	ctx    context.Context    = nil
	cancel context.CancelFunc = nil
	hNode  *core.HostNode
)

func PanelTab() fyne.CanvasObject {
	c := container.New(layout.NewVBoxLayout())

	lab := widget.NewLabel("节点控制")
	lab.Alignment = fyne.TextAlignCenter

	btnOn := widget.NewButton("启动", func() {
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

	btnOff := widget.NewButton("关闭", func() {
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

	btnUpload := widget.NewButton("上传文件", func() {

		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err2 error) {

			p := closer.URI().Path()
			t := closer.URI().Extension()
			content, e := tool.LoadFile(p)
			if e != nil {
				ShowErr(e)
				return
			}
			file := tool.NewFileFromContent(t, content)
			e = storage.StoreFileData(file)
			if e != nil {
				ShowErr(e)
				return
			}
			ShowInfo("文件上传成功")

		}, w)
	})
	btnUpload.Alignment = widget.ButtonAlignCenter

	c.Add(lab)
	c.Add(btnOn)
	c.Add(btnOff)
	c.Add(btnUpload)

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
	if f := <-h.JoinNetwork(); f == 0 {
		ShowErr(err.ErrJoinNetwork)
		hNode = nil
		msg <- false
		return
	}
	hNode = h
	msg <- true

	<-ctx.Done()
	h.Close()

}
