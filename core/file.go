package core

import (
	"math/rand"
	"sardines/err"
	"sardines/storage"
	"sardines/tool"
	"strconv"
)

func (h *HostNode) UploadFile(file *tool.File) (string, string, error) {

	// update file to ipfs
	cid, err2 := h.api.Upload(file.Raw())
	if err2 != nil {
		return "", "", err2
	}
	file.CID = cid

	// update the keyTable
	intn := rand.Intn(100)
	if b := h.Ktab.Append(strconv.Itoa(intn), []string{file.CID}); !b {
		return "", "", err.KeyTableUpdateErr
	}

	// update the manifest
	err2 = storage.UpdateManifest(file.Entry)
	if err2 != nil {
		return "", "", err2
	}

	return file.ID(), strconv.Itoa(intn), nil
}

func (h *HostNode) SearchFileByCid(cid string) (*tool.File, error) {
	bytes, e := h.api.Download(cid)
	if e != nil {
		return nil, e
	}

	file, e := tool.NewFileFromRaw(bytes)
	if e != nil {
		return nil, e
	}
	file.CID = cid

	e = storage.UpdateManifest(file.Entry)
	if e != nil {
		return file, e
	}
	return file, nil
}

func (h *HostNode) SearchFileByKey(kw string) []*tool.File {
	cids := h.Ktab.Get(kw)
	files := make([]*tool.File, 0, 5)
	for _, cid := range cids {
		f, _ := h.SearchFileByCid(cid)
		files = append(files, f)
	}
	return files
}
