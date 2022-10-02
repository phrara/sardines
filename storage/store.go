package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sardines/config"
	"sardines/err"
	"sardines/tool"
)

var manifest map[string]string

func init() {
	manifest = make(map[string]string)
	exist, _ := tool.HasDir(config.Manifest)
	if exist {
		file, e := tool.LoadFile(config.Manifest)
		if e != nil {
			return
		}
		e = json.Unmarshal(file, &manifest)
		if e != nil {
			return
		}
	} else {
		return
	}
}

func StoreFileData(file *tool.File) error {
	subDir := file.ID()[2:5]
	fileDir := filepath.Join(config.FS, "/"+subDir)
	er := tool.CreateDir(fileDir)
	if er != nil && er != err.DirExists {
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

func DeleteFileData(file *tool.File) {
	subDir := file.ID()[2:5]
	filePath := filepath.Join(config.FS, subDir, file.ID())
	os.Remove(filePath)
}

func FileStoreTree() map[string][]string {
	res := make(map[string][]string)
	res[""] = []string{}

	dirs, _ := os.ReadDir(config.FS)
	for _, d := range dirs {
		if d.IsDir() {
			res[""] = append(res[""], d.Name())
			files, er := os.ReadDir(filepath.Join(config.FS, d.Name()))
			if er != nil {
				continue
			}
			for _, f := range files {
				if !f.IsDir() {
					str := fmt.Sprintf("%s| %s ", manifest[f.Name()], f.Name())
					res[d.Name()] = append(res[d.Name()], str)
				}
			}
		}
	}
	return res
}

func UpdateManifest(entry tool.Entry) error {
	exist, err2 := tool.HasDir(config.Manifest)
	if err2 != nil {
		return err2
	}
	if exist {
		file, err2 := tool.LoadFile(config.Manifest)
		if err2 != nil {
			return err2
		}
		if len(file) != 0 {
			mf := make(map[string]string, 0)
			err2 = json.Unmarshal(file, &mf)
			if err2 != nil {
				return err2
			}
			mf[entry.FID] = entry.Origin
			j, err2 := json.Marshal(mf)
			if err2 != nil {
				return err2
			}
			er := tool.WriteFile(j, config.Manifest)
			if er != nil {
				return er
			}
			manifest = mf
			return nil
		} else {
			mf := make(map[string]string)
			mf[entry.FID] = entry.Origin
			j, err2 := json.Marshal(mf)
			if err2 != nil {
				return err2
			}
			er := tool.WriteFile(j, config.Manifest)
			if er != nil {
				return er
			}
			manifest = mf
			return nil
		}
	} else {
		mf := make(map[string]string)
		mf[entry.FID] = entry.Origin
		j, err2 := json.Marshal(mf)
		if err2 != nil {
			return err2
		}
		er := tool.WriteFile(j, config.Manifest)
		if er != nil {
			return er
		}
		manifest = mf
		return nil
	}
}

func OriginalNameFromFID(fid string) string {
	return manifest[fid]
}
