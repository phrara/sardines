package tool

import (
	"bytes"
	"encoding/binary"
)

const HEADER = 8

// Packet
// "TLV" style binary packets: Tag | Len | Value
type Packet struct {
	Tag   uint32
	Len   uint32
	Value []byte
}

func (p *Packet) Wrap() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, p.Tag)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, p.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, p.Value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Packet) ParseHeader(header []byte) error {
	buf := bytes.NewBuffer(header)
	err := binary.Read(buf, binary.LittleEndian, &(p.Tag))
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &(p.Len))
	if err != nil {
		return err
	}

	return nil
}

func (p *Packet) ValString() string {
	return string(p.Value)
}
