package tool

import (
	"fmt"
	"testing"
)

func TestCalc(t *testing.T) {
	d := GetPeerDist("QmUsMCHhYvx2LgQcv7KN7QrbpmPdB97D3JsZgGMdQqQ6CF", "QmUsMCHhYvx2LgQcv7KN7QrbpmPdB97D3JsZgGMdQqQ6CF")
	fmt.Println(d)
}
