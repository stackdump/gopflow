package statemachine_test

import (
	"fmt"
	. "github.com/stackdump/gopflow/statemachine"
	"testing"

	. "github.com/stackdump/gopflow/ptnet"
)

func TestLoadFromFile(t *testing.T) {
	p := LoadFile("../examples/octoe.pflow")
	m := p.StateMachine()
	print(m.String())
}

func TestTransformations(t *testing.T) {
	m := LoadFile("../examples/octoe.pflow").StateMachine()
	s := m.Initial

	xFail := func(action string, multiple int, roleIn string) {
		expectedRole := Role(roleIn)
		stateOut, role, err := m.Transform(s, action, uint64(multiple))
		fmt.Printf("%v, %v, %v\n", stateOut, role, err)

		if role == expectedRole && err == nil {
			t.Fatal("Expected Error")
		}
	}

	xPass := func(action string, multiple int, roleIn string) {
		stateOut, role, err := m.Transform(s, action, uint64(multiple))
		fmt.Printf("%v, %v, %v\n", stateOut, role, err)

		expectedRole := Role(roleIn)
		if role != expectedRole {
			t.Fatalf("expected role %v does not match %v", role, expectedRole)
		}

		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		if err == nil {
			for k, v := range stateOut {
				s[k] = uint64(v)
			}
		}
	}

	// Test guards
	xFail("EXEC", 1, "FATD") // fails guard clause
	xPass("ON", 1, "FATD")   //valid
	xPass("EXEC", 1, "FATD") //valid

	// Test state validation
	xFail("X11", 1, "PlayerO") // bad role
	xFail("X11", 2, "PlayerX") // bad multiple
	xPass("X11", 1, "PlayerX") // valid
	xFail("X11", 1, "PlayerX") // move already taken
	xPass("O01", 1, "PlayerO") // valid

}
