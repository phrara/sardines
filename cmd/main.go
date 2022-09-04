package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sardines/config"
	"sardines/core"
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

	if b := cliArgsParse(); b {
		ctx, cancel := context.WithCancel(context.Background())
		go run(ctx)

		stdReader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("dctrl> ")
			signal, err := stdReader.ReadString('\n')
			if err != nil {
				log.Println(err)
				return
			}
			switch signal {
			case "quit\r\n":
				cancel()
				time.Sleep(time.Second * 2)
				fmt.Println("bye")
				return
			default:
				continue
			}

		}

	} else {
		return
	}

}

func run(ctx context.Context) {
	fmt.Printf("\u001B[1;35m%s\u001B[0m\n", LOGO)
	hnode, err := core.GenerateNode()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\x1b[1;34mHost: %s\x1b[0m\n", hnode.NodeAddr.String())

	hnode.JoinNetwork()
	time.Sleep(time.Second * 2)
	//fmt.Println("recent router table:\n", string(hnode.Router.RawData()))

	hnode.RouterDistributeOn(false, 30, ctx)
	time.Sleep(time.Second * 2)
	//hnode.Serv.SendFile(core.BootstrapNodes[0], "A test file content")

	<-ctx.Done()
	hnode.Close()

}

func cliArgsParse() bool {

	var username string
	var password string
	var port string
	var rs int64
	var bsn string

	flag.StringVar(&username, "u", "root", "username, default is root")
	flag.StringVar(&password, "p", "root", "password, default is root")
	flag.StringVar(&bsn, "b", "", "Bootstrap Node")
	flag.StringVar(&port, "P", "8082", "host's port, default is 8082")
	flag.Int64Var(&rs, "r", 1, "random seed is used for generating your own key-pair, it must be a positive number")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("please use -help for more guidance!")
		return false
	}

	switch flag.Arg(0) {
	case "conf":
		c, err := config.New(username, password, "127.0.0.1:"+port, rs, bsn)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if b := c.Save(); b {
			fmt.Println("configure successfully")
		} else {
			fmt.Println("configure failed")
		}
		return false
	case "run": 
		return true
	case "gen-key":
		c := (&config.Config{}).Load()
		if c == nil {
			fmt.Println("your node haven't been configure correctly, please use -help for more guidance")
			return false
		}
		if err := c.GenKey(); err != nil {
			fmt.Println(err)
			return false
		}
		if err := c.SaveKey(); err != nil {
			fmt.Println(err)
		}
		return false
	case "init":
		if err := tool.CreateDir(config.Dir); err != nil {
			fmt.Println(err)
		}
		return false
	default:
		fmt.Println("please use -help for more guidance!")
		return false
	}

}

