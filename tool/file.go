package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type File struct {
	Type    string `json:"type"`
	Content []byte `json:"content"`
	FID string `json:"fid"`
}

func NewFile(typ, fid string, content []byte) *File {
	return &File{
		Type:    typ,
		Content: content,
		FID: fid,
	}
}

func (f *File) ID() string {
	return f.FID
}

func (f *File) Wrap() []byte {
	wrap, err := json.Marshal(*f)
	if err != nil {
		return nil
	}
	return wrap
}

func (f *File) Unwrap(wrap []byte) *File {
	err := json.Unmarshal(wrap, f)
	if err != nil {
		fmt.Println(err)
	}
	return f
}

func (f *File) Size() int {
	return len(f.Content)
}

func WriteFile(b []byte, path string) error {
	err := os.WriteFile(path, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}


//判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
 
//创建文件夹
func CreateDir(path string) error {
	exist, err := HasDir(path)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("dir: `values` already exists")
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}
