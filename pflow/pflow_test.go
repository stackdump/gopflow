package pflow_test

import (
	"testing"

	. "github.com/stackdump/gopflow/pflow"
)

func TestCounterMachine(t *testing.T) {
	_, err := LoadFile("../examples/octoe.pflow")

	if err != nil {
		t.Fatal("failed to unmarshal")
	}
}
