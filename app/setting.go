package app

import (
	"encoding/hex"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sardines/config"
	"sardines/err"
	"strconv"
)

func SettingTab() fyne.CanvasObject {
	c := container.New(layout.NewVBoxLayout())

	lab := widget.NewLabel("Configuration")
	lab.Alignment = fyne.TextAlignCenter

	ipEntry := widget.NewEntry()
	ipEntry.SetText("0.0.0.0")
	portEntry := widget.NewEntry()
	portEntry.SetText("8082")
	rsEntry := widget.NewEntry()
	rsEntry.SetText("1")
	bsEntry := widget.NewEntry()

	conf := &config.Config{}
	er := conf.Load()
	if er != nil {
		ShowErr(er)
	} else {
		ipEntry.SetText(conf.IP)
		portEntry.SetText(conf.Port)
		rsEntry.SetText(strconv.FormatInt(conf.RandomSeed, 10))
		bsEntry.SetText(conf.BootstrapNode)
	}

	cForm := widget.NewForm()
	cForm.Append("IP Address", ipEntry)
	cForm.Append("Listening Port", portEntry)
	cForm.Append("Random Seed", rsEntry)
	cForm.Append("Bootstrap Node", bsEntry)

	btns := container.New(layout.NewHBoxLayout())
	// 保存配置
	submit := widget.NewButton("Submit", func() {
		conf.IP = ipEntry.Text
		conf.Port = portEntry.Text
		conf.BootstrapNode = bsEntry.Text
		conf.RandomSeed, er = strconv.ParseInt(rsEntry.Text, 10, 64)
		if er != nil {
			ShowErr(er)
			return
		}
		if b := conf.Save(); b {
			ShowInfo("configure successfully")
		} else {
			ShowErr(err.ErrConf)
		}

	})

	keyLab := widget.NewLabel("Priv Key")
	keyLab.Alignment = fyne.TextAlignCenter
	keyEntry := widget.NewEntry()
	keyEntry.Disable()
	keyEntry.MultiLine = true
	keyEntry.Wrapping = fyne.TextWrapBreak

	// 生成密钥对
	genKey := widget.NewButton("Gen Key", func() {
		co := &config.Config{}
		er = co.Load()
		if er != nil {
			ShowErr(er)
			return
		}
		err2 := co.GenKey()
		if err2 != nil {
			ShowErr(err2)
			return
		}
		err2 = co.SaveKey()
		if err2 != nil {
			ShowErr(err2)
			return
		}
		ShowInfo("gen key-pair successfully")
		raw, _ := co.PrvKey.Raw()
		keyEntry.SetText(hex.EncodeToString(raw))
	})
	btns.Add(widget.NewLabel("                              "))
	btns.Add(submit)
	btns.Add(widget.NewLabel("                                               "))
	btns.Add(genKey)

	md := widget.NewRichTextFromMarkdown("**please do not lose or share you own private key!**")

	c.Add(lab)
	c.Add(cForm)
	c.Add(btns)
	c.Add(keyLab)
	c.Add(keyEntry)
	c.Add(md)

	return c
}
