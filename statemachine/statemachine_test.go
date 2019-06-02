package statemachine_test

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

	if p == nil {
		return
	}

	//println(p.String())
	m := p.StateMachine()

	if m == nil {
		t.Fatalf("failed to load state machine")
	}
	print(m.String())
}
