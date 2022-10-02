package app

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sardines/core"
	"sardines/err"
	"sardines/tool"
)

var (
	ctx    context.Context    = nil
	cancel context.CancelFunc = nil
	btnW   float32            = 500
	btnH   float32            = 40
	btnX   float32            = 350
	// 节点句柄
	hNode *core.HostNode
)

// 启动按钮
func btnOn() *widget.Button {
	btn := widget.NewButton("启动", func() {
		if cancel != nil || hNode != nil {
			return
		}
		ctx, cancel = context.WithCancel(context.Background())
		msg := make(chan bool, 1)
		go Run(ctx, msg)
		if b := <-msg; b {
			turnOn(hNode.NodeAddr.String())
		}
	})
	btn.Alignment = widget.ButtonAlignCenter
	btn.Move(fyne.NewPos(btnX, 50))
	btn.Resize(fyne.NewSize(btnW, btnH))
	return btn
}

// 关闭按钮
func btnOff() *widget.Button {
	btn := widget.NewButton("关闭", func() {
		if cancel != nil && hNode != nil {
			Ensure("are you sure?", func(b bool) {
				if b {
					cancel()
					turnOff()
					cancel = nil
				}
			})
		}
	})
	btn.Alignment = widget.ButtonAlignCenter
	btn.Move(fyne.NewPos(btnX, 100))
	btn.Resize(fyne.NewSize(btnW, btnH))
	return btn
}

func btnUpload() *widget.Button {
	btn := widget.NewButton("上传文件", func() {

		if cancel != nil && hNode != nil {
			fileDialog := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err2 error) {
				if closer != nil {
					p := closer.URI().Path()
					o := closer.URI().Name()
					content, e := tool.LoadFile(p)
					if e != nil {
						ShowErr(e)
						return
					}
					file := tool.NewFileFromContent(o, content)

					fid, e := hNode.UploadFile(file)
					if e != nil {
						ShowErr(e)
						return
					}
					ShowData("上传成功", fid)
					fileTree.Refresh()
				}
			}, w)
			fileDialog.Resize(fyne.NewSize(900, 600))
			fileDialog.Show()

		} else {
			ShowErr(err.NodeNotStarted)
		}

	})
	btn.Alignment = widget.ButtonAlignCenter
	btn.Move(fyne.NewPos(btnX, 150))
	btn.Resize(fyne.NewSize(btnW, btnH))
	return btn
}

func PanelTab() fyne.CanvasObject {
	c := container.NewWithoutLayout()

	lab := widget.NewLabel("节点控制")
	lab.Alignment = fyne.TextAlignCenter
	lab.Move(fyne.NewPos(title, 0))

	c.Add(lab)
	c.Add(btnOn())
	c.Add(btnOff())
	c.Add(btnUpload())

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
		ShowErr(err.JoinNetworkFailed)
		hNode = nil
		msg <- false
		return
	}
	hNode = h
	msg <- true

	<-ctx.Done()
	hNode.Close()
	hNode = nil
}
