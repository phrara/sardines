package tool

import (
	"encoding/json"
	"fmt"
	"os"
)

type File struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewFile(typ, content string) *File {
	return &File{
		Type:    typ,
		Content: content,
	}
}

func (f *File) Wrap() []byte {
	wrap, err := json.Marshal(*f)
	if err != nil {
		return nil
	}
	wrap = append(wrap, '\n')
	return wrap
}

func (f *File) Unwrap(wrap []byte) *File {
	wrap = wrap[:len(wrap)-1]
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
	err := os.WriteFile(path, b, 02)
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
