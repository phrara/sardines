package tool

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"sardines/err"
)

type Entry struct {
	Origin string `json:"origin"`
	FID    string `json:"fid"`
}

type File struct {
	Entry
	Content []byte `json:"content"`
}

func NewFile(origin, fid string, content []byte) *File {
	return &File{
		Content: content,
		Entry: Entry{
			Origin: origin,
			FID:    fid,
		},
	}
}

func NewFileFromContent(origin string, content []byte) *File {
	fid, e := HashEncode(content)
	if e != nil {
		return nil
	}
	return NewFile(origin, "SF"+hex.EncodeToString(fid), content)
}

func NewFileFromRaw(raw []byte) (*File, error) {
	f := new(File)
	err2 := json.Unmarshal(raw, f)
	if err2 != nil {
		return nil, err2
	}
	return f, nil
}

func (f *File) ID() string {
	return f.FID
}

func (f *File) Raw() []byte {
	wrap, err2 := json.Marshal(*f)
	if err2 != nil {
		return nil
	}
	return wrap
}

func (f *File) Size() int {
	return len(f.Content)
}

func (f *File) Info() string {
	return f.Origin + " | " + f.FID
}

func WriteFile(b []byte, path string) error {
	err2 := os.WriteFile(path, b, os.ModePerm)
	if err2 != nil {
		return err2
	}
	return nil
}

func LoadFile(path string) ([]byte, error) {
	file, err2 := os.ReadFile(path)
	if err2 != nil {
		return nil, err2
	}
	return file, nil
}

// HasDir 判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, err2 := os.Stat(path)
	if err2 == nil {
		return true, nil
	}
	if os.IsNotExist(err2) {
		return false, nil
	}
	return false, err2
}

// CreateDir 创建文件夹
func CreateDir(path string) error {
	exist, er := HasDir(path)
	if er != nil {
		return er
	}
	if exist {
		return err.DirExists
	} else {
		er = os.Mkdir(path, os.ModePerm)
		if er != nil {
			return er
		} else {
			return nil
		}
	}
}
