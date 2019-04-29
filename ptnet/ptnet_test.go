package ptnet_test

import (
	"testing"

	. "github.com/stackdump/gopflow/pflow"
)

func TestCounterMachine(t *testing.T) {
	p, err := LoadFile("../examples/counter.pflow")
	_ = err
	_ = p
}
