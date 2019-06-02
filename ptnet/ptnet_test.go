package ptnet_test

import (
	"testing"

	. "github.com/stackdump/gopflow/ptnet"
)

func TestLoadFromFile(t *testing.T) {
	p := LoadFile("../examples/octoe.pflow")
	println(p.String())
}
