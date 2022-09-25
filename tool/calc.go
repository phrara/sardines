package tool

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// GetPeerDist
// To get distance between peers
func GetPeerDist(dst, src string) int {
	dst = dst[2:]
	src = src[2:]
	d, _ := HashEncode([]byte(dst))
	s, _ := HashEncode([]byte(src))
	return calcDist(d, s)
}

// GetFileDist
// To get distance between peer and file
func GetFileDist(pid, fid string) int {
	p, _ := HashEncode([]byte(pid))
	fid = fid[1:]
	b, _ := hex.DecodeString(fid)
	return calcDist(p, b)
}

func calcDist(d, s []byte) int {
	for i, v := range d {
		res := v ^ s[i]
		if res != 0 {
			switch {
			case res > 127:
				return 256 - i*8
			case res > 63:
				return 256 - (i*8 + 1)
			case res > 31:
				return 256 - (i*8 + 2)
			case res > 15:
				return 256 - (i*8 + 3)
			case res > 7:
				return 256 - (i*8 + 4)
			case res > 3:
				return 256 - (i*8 + 5)
			case res > 1:
				return 256 - (i*8 + 6)
			case res == 1:
				return 256 - (i*8 + 7)
			default:
				return 256 - ((i + 1) * 8)
			}
		}
	}
	return 0
}

func HashEncode(ctnt []byte, readers ...io.Reader) ([]byte, error) {
	h := sha256.New()
	if len(ctnt) != 0 {
		h.Write(ctnt)
	}
	if len(readers) != 0 {
		for _, v := range readers {
			if _, err := io.Copy(h, v); err != nil {
				return nil, err
			}
		}
	}
	return h.Sum(nil), nil
}