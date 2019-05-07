package ptnet_test

import (
"testing"

. "github.com/stackdump/gopflow/ptnet"
)

func TestLoadFromFile(t *testing.T) {
	p, err := LoadFile("../examples/octoe.pflow")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Document().SubNets) != 1 {
		t.Fatal("failed to load xml file")
	}
	_, err = p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	//println(string(data))
}
