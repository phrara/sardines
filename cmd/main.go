package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sardines/app"
	"sardines/config"
	"sardines/core"
	"sardines/err"
	"sardines/tool"
	"time"
)

const LOGO = `
'        ___           ___           ___       ___       ___     
'       /\__\         /\  \         /\__\     /\__\     /\  \    
'      /:/  /        /::\  \       /:/  /    /:/  /    /::\  \   
'     /:/__/        /:/\:\  \     /:/  /    /:/  /    /:/\:\  \  
'    /::\  \ ___   /::\~\:\  \   /:/  /    /:/  /    /:/  \:\  \ 
'   /:/\:\  /\__\ /:/\:\ \:\__\ /:/__/    /:/__/    /:/__/ \:\__\
'   \/__\:\/:/  / \:\~\:\ \/__/ \:\  \    \:\  \    \:\  \ /:/  /
'        \::/  /   \:\ \:\__\    \:\  \    \:\  \    \:\  /:/  / 
'        /:/  /     \:\ \/__/     \:\  \    \:\  \    \:\/:/  /  
'       /:/  /       \:\__\        \:\__\    \:\__\    \::/  /   
'       \/__/         \/__/         \/__/     \/__/     \/__/    `

func main() {

	if b := cliArgsParse(); b == 1 {
		ctx, cancel := context.WithCancel(context.Background())
		go Run(ctx)

		stdReader := bufio.NewReader(os.Stdin)

		time.Sleep(time.Second * 2)
		for {
			fmt.Print("sardines> ")
			signal, err2 := stdReader.ReadString('\n')
			if err2 != nil {
				log.Println(err2)
				return
			}
			switch signal {
			case "quit\r\n":
				cancel()
				time.Sleep(time.Second * 2)
				fmt.Println("bye")
				return
			case "sfa ":
			default:
				continue
			}

		}

	} else if b == 2 {
		app.App()
		return
	}

}

func Run(ctx context.Context) {
	fmt.Printf("\u001B[1;35m%s\u001B[0m\n", LOGO)
	hnode, er := core.GenerateNode()
	if er != nil {
		fmt.Println(er)
		return
	}

	fmt.Printf("\x1b[1;34mHost: %s\x1b[0m\n", hnode.NodeAddr.String())

	msg := hnode.JoinNetwork()
	if f := <-msg; f == 0 {
		fmt.Println(err.ErrJoinNetwork)
		return
	}
	time.Sleep(time.Second * 2)
	//fmt.Println("recent router table:\n", string(hnode.Router.RawData()))

	hnode.RouterDistributeOn(false, 30, ctx)
	time.Sleep(time.Second * 2)
	//hnode.Serv.SendFile(core.BootstrapNodes[0], "A test file content")

	<-ctx.Done()
	hnode.Close()

}

func cliArgsParse() uint {

	var username string
	var password string
	var port string
	var rs int64
	var bsn string
	var ip string

	flag.StringVar(&username, "u", "root", "username, default is root")
	flag.StringVar(&password, "p", "root", "password, default is root")
	flag.StringVar(&bsn, "b", "", "Bootstrap Node")
	flag.StringVar(&port, "P", "8082", "host's port, default is 8082")
	flag.StringVar(&ip, "i", "0.0.0.0", "-i x.x.x.x")
	flag.Int64Var(&rs, "r", 1, "random seed is used for generating your own key-pair, it must be a positive number")
	flag.Parse()

	switch flag.Arg(0) {
	case "conf":
		c, err := config.New(username, password, ip+":"+port, rs, bsn)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		if b := c.Save(); b {
			fmt.Println("configure successfully")
		} else {
			fmt.Println("configure failed")
		}
		return 0
	case "Run":
		return 1
	case "gen-key":
		c := &config.Config{}
		err := c.Load()
		if err != nil {
			fmt.Printf("your node haven't been configure correctly, please use -help for more guidance, err: %v", err)
			return 0
		}
		if err := c.GenKey(); err != nil {
			fmt.Println(err)
			return 0
		}
		if err := c.SaveKey(); err != nil {
			fmt.Println(err)
		}
		return 0
	case "init":
		if err := tool.CreateDir(config.Dir); err != nil {
			fmt.Println(err)
			return 0
		}
		if err := tool.CreateDir(config.FS); err != nil {
			fmt.Println(err)
		}
		return 0
	default:
		return 2
	}

}
