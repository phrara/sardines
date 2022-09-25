package storage

import (
	"path/filepath"
	"sardines/config"
	"sardines/err"
	"sardines/tool"
)

func StoreFileData(file *tool.File) error {
	subDir := file.ID()[24:26]
	fileDir := filepath.Join(config.FS, "/"+subDir)
	er := tool.CreateDir(fileDir)
	if er != nil && er != err.ErrDirExists {
		return er
	}
	filePath := filepath.Join(fileDir, "/"+file.ID())
	err2 := tool.WriteFile(file.Content, filePath)
	if err2 != nil {
		return err2
	}
	return nil
}

func LoadFileData(fid string) (*tool.File, error) {
	subDir := fid[24:26]
	filePath := filepath.Join(config.FS, "/"+subDir, "/"+fid)
	b, err2 := tool.LoadFile(filePath)
	if err2 != nil {
		return nil, err2
	}
	f := tool.NewFile("txt", fid, b)
	return f, nil
} 