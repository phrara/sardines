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

func (h *HostNode) JoinNetwork() <-chan int {
	msg := make(chan int, 1)
	for _, bn := range BootstrapNodes {
		node := bn
		if node == nil {
			//fmt.Println("Bootstrap Node is empty, this node could be the initial node")
			continue
		}
		go func() {
			res := <-h.Serv.Ping(node)
			if res.Error != nil {
				msg <- 0
				return
			} else {
				// ping 通了
				// 发出加入申请
				if b := h.Serv.JoinApply(node); b {
					msg <- 1
					h.Router.AddNode(node)
				} else {
					msg <- 0
				}
			}
		}()
	}
	msg <- 1
	return msg
}

func (h *HostNode) Close() error {
	err := h.Host.Close()
	h.Ktab.Close()
	if h.ipfs != nil {
		h.ipfs.Process.Kill()
	}
	return err
}

func (h *HostNode) RunIpfs() error {
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
			h.ipfs = cmd
			cmd.Wait()
		} else {
			cmd.Process.Kill()
			return errors.New("ipfs daemon run failed")
		}
	}
	return nil

}

func (h *HostNode) routerDistribute(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Second * period)
	for {
		<-ticker.C

		// before routerdistributing
		errNum := h.Serv.RouterDistribute()

		if len(BootstrapNodes) != 0 {

			// If the errNum > 33% of the sum of nodes,
			// we regard this as the bad situation of network, then try again after 8 sec;
			// if the errNum > 75% of the sum of nodes,
			// we regard this as the fatal error of network, stop the ticker
			if errNum > h.Router.Sum()/4*3 {
				fmt.Println("fatal Network error")
				fmt.Println("Please restart your server")
				ticker.Stop()
				return
			} else if errNum > h.Router.Sum()/3 {
				fmt.Println("Bad Network Situation")
				time.Sleep(time.Second * (period / 2))
				h.Serv.RouterDistribute()
			}
		}

		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		default:
		}
	}
}

func (h *HostNode) autoPeriod(nodeNum int) time.Duration {
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
func (h *HostNode) RouterDistributeOn(auto bool, period int, ctx context.Context) {
	if auto {
		go h.routerDistribute(ctx, h.autoPeriod(h.Router.Sum()))
	} else {
		go h.routerDistribute(ctx, time.Duration(period))
	}
}
