package core

import (
	"sardines/storage"
	"sardines/tool"
)

func (h *HostNode) UploadFile(file *tool.File) (string, error) {

	// store the file locally
	err2 := storage.StoreFileData(file)
	if err2 != nil {
		return "", err2
	}

	// update the manifest
	err2 = storage.UpdateManifest(file.Entry)
	if err2 != nil {
		storage.DeleteFileData(file)
		return "", err2
	}

	// send to remote peer

	// TODO: update the keyTable

	//h.DHT.RoutingTable()

	return file.ID(), nil
}
