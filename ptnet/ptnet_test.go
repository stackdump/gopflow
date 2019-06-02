package ptnet_test

import (
	"testing"

	. "github.com/stackdump/gopflow/ptnet"
)

func TestLoadFromFile(t *testing.T) {
	p, err := LoadFile("../examples/octoe.pflow")
	if err != nil {
		t.Fatal(err)
		return
	}

	if p != nil {
		println(p.String())
	}
}
