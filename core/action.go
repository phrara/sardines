package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func (n *HostNode) JoinNetwork() {
	n.Router.Clear()
	for _, bn := range BootstrapNodes {
		node := bn
		if node == nil {
			fmt.Println("Bootstrap Node is empty, this node could be the initial node")
			continue
		}
		go func() {
			res := <-n.Serv.Ping(node)
			if res.Error != nil {
				fmt.Println("can not connect Bootstrap Node")
				return
			} else {
				// ping 通了
				// 发出加入申请
				if b := n.Serv.JoinApply(node); b {
					fmt.Println("Join Network Successfully")
					n.Router.AddNode(node)
				} else {
					fmt.Println("Join Network Failed")
				}
			}
		}()
	}
}

func (n *HostNode) Close() error {
	err := n.Host.Close()
	n.Router.Clear()
	n.Ktab.Close()
	if n.ipfs != nil {
		n.ipfs.Process.Kill()
	}
	return err
}

func (n *HostNode) RunIpfs() error {
	_, err := exec.Command("ipfs", "init").Output()
	if err != nil {
		return err
	}
	cmd := exec.Command("ipfs", "daemon")
	rc, err2 := cmd.StdoutPipe()
	if err2 != nil {
		return err2
	}
	if err2 = cmd.Start(); err2 != nil {
		return err2
	}

	if b, err3 := io.ReadAll(rc); err3 != nil {
		return err3
	} else {
		if strings.Contains(string(b), "ready") {
			n.ipfs = cmd
			cmd.Wait()
		} else {
			cmd.Process.Kill()
			return errors.New("ipfs daemon run failed")
		}
	}
	return nil

}

func (n *HostNode) routerDistribute(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Second * period)
	for  {
		<- ticker.C

		//fmt.Println("Router Distribution Start: ", t.Format("2006-01-02 15:04:05"))
		errNum := n.Serv.RouterDistribute()

		// If the errNum > 33% of the sum of nodes,
		// we regard this as the bad situation of network, then try again after 8 sec;
		// if the errNum > 75% of the sum of nodes,
		// we regard this as the fatal error of network, stop the ticker
		if errNum > n.Router.Sum()/4*3 {
			fmt.Println("fatal Network error")
			fmt.Println("Please restart your server")
			ticker.Stop()
			return
		} else if errNum > n.Router.Sum()/3 {
			fmt.Println("Bad Network Situation")
			time.Sleep(time.Second * (period / 2))
			n.Serv.RouterDistribute()
		}
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		default:
		}
	}
}

func (n *HostNode) autoPeriod(nodeNum int) time.Duration {
	switch {
	case nodeNum >= 100:
		return 60
	case nodeNum >= 50:
		return 30
	case nodeNum >= 25:
		return 15
	default:
		return 10
	}
}

// RouterDistributeOn
// The unit of argument `period` is second.
// The argument `period` will be affective if `auto` is true,
// and it'll be useless if `auto` is false.
func (n *HostNode) RouterDistributeOn(auto bool, period int, ctx context.Context) {
	if auto {
		go n.routerDistribute(ctx, n.autoPeriod(n.Router.Sum()))
	} else {
		go n.routerDistribute(ctx, time.Duration(period))
	}
}
