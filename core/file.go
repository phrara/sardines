package core

import (
	"errors"
	"sardines/tool"
)

func (h *HostNode) StoreFile(ctnt, path string) (string, error) {
	if ctnt == "" && path == "" {
		return "", errors.New("there is nothing to be stored")
	}
	ctntBytes := []byte(ctnt)
	if path != "" {
		b, err := tool.LoadFile(path)
		if err != nil {
			return "", nil
		}
		ctntBytes = append(ctntBytes, b...)
	}


	fid, err := tool.HashEncode(ctntBytes)
	if err != nil {
		return "", err
	}	
	file := tool.NewFile("txt", "F"+string(fid), ctntBytes)

	dist := tool.GetFileDist(h.NodeInfo.ID.String(), file.FID)
	l := h.Router.GetNodes(dist)
	if l != nil {
		for e := l.Front(); e != nil; e = e.Next() {
			pn := e.Value.(*tool.PeerNode)
			go func() {
				h.Serv.SendFile(pn, file)
			}()
			
		}
	}
	
	// TODO: update the keyTable 


	// TODO: store the file 



	return "", nil
}