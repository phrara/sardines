package tool

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

func GetDistByXor(dst, src peer.ID) int {
	dst = dst[2:]
	src = src[2:]
	for i, v := range dst {
		res := v ^ int32(src[i])
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
