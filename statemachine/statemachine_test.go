package statemachine_test

import (
	"testing"

	. "github.com/stackdump/gopflow/pflow"
)

func TestCounterMachine(t *testing.T) {
	p, err := LoadFile("../examples/octoe.pflow")
	_ = err
	_ = p
}
