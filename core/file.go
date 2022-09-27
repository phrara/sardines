package core

import (
	"encoding/hex"
	"sardines/err"
	"sardines/storage"
	"sardines/tool"
)

func (h *HostNode) StoreFile(ctnt, path string) (string, error) {

	if ctnt == "" && path == "" {
		return "", err.ErrNothingToStore
	}
	ctntBytes := []byte(ctnt)
	if path != "" {
		b, er := tool.LoadFile(path)
		if er != nil {
			return "", er
		}
		ctntBytes = append(ctntBytes, b...)
	}

	fid, err := tool.HashEncode(ctntBytes)
	if err != nil {
		return "", err
	}
	file := tool.NewFile("txt", "F"+hex.EncodeToString(fid), ctntBytes)

	// * store the file
	err2 := storage.StoreFileData(file)
	if err2 != nil {
		return "", err2
	}

	// send to remote peer

	// TODO: update the keyTable

	h.ipfsDHT.RoutingTable()

	return file.ID(), nil
}
