package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sardines/config"
	"sardines/err"
	"sardines/tool"
)

func StoreFileData(file *tool.File) error {
	subDir := file.ID()[2:5]
	fileDir := filepath.Join(config.FS, "/"+subDir)
	er := tool.CreateDir(fileDir)
	if er != nil && er != err.ErrDirExists {
		return er
	}
	filePath := filepath.Join(fileDir, "/"+file.ID())
	err2 := tool.WriteFile(file.Raw(), filePath)
	if err2 != nil {
		return err2
	}
	return nil
}

func LoadFileData(fid string) (*tool.File, error) {
	subDir := fid[2:5]
	filePath := filepath.Join(config.FS, "/"+subDir, "/"+fid)
	b, err2 := tool.LoadFile(filePath)
	if err2 != nil {
		return nil, err2
	}
	f := tool.NewFile("txt", fid, b)
	return f, nil
}

func FileStoreTree() map[string][]string {
	res := make(map[string][]string)
	res[""] = []string{"FileStore"}

	dirs, _ := os.ReadDir(config.FS)
	for _, d := range dirs {
		if d.IsDir() {
			res["FileStore"] = append(res["FileStore"], d.Name())
			files, er := os.ReadDir(filepath.Join(config.FS, d.Name()))
			if er != nil {
				fmt.Println(er)
			}
			for _, f := range files {
				if !f.IsDir() {
					res[d.Name()] = append(res[d.Name()], f.Name())
				}
			}
		}
	}
	return res
}
