package statemachine

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type StateVector []uint64
type Action string
type Role string
type Delta []int64
type Condition string

type Transition struct {
	Delta  Delta
	Role   Role
	Guards map[Condition]Delta
}

type StateMachine struct {
	Initial     StateVector
	Capacity    StateVector
	Transitions map[Action]Transition
}

func (s StateMachine) guardTest(state []uint64, tx Transition) error {
	TESTING: for label, g := range tx.Guards {
		for offset, delta := range g {
			val := int64(state[offset]) + delta
			if val < 0 {
				continue TESTING
			}
		}
		return errors.New(fmt.Sprintf("guard failure: %v", label))
	}
	return nil
}

func (s *StateMachine) Transform(state []uint64, action string, multiplier uint64) (vectorOut []int64, role Role, err error) {

	t := s.Transitions[Action(action)]

	for offset, delta := range t.Delta {
		val := int64(state[offset]) + delta*int64(multiplier)
		vectorOut = append(vectorOut, val)

		if err == nil && val < 0 {
			err = errors.New(fmt.Sprintf("underflow offset: %v => %v ", offset, val))
		}

		if err == nil && s.Capacity[offset] != 0 && val > int64(s.Capacity[offset]) {
			err = errors.New(fmt.Sprintf("overflow offset: %v => %v ", offset, val))
		}
	}

	if err == nil {
		err = s.guardTest(state, t)
	}

	return vectorOut, t.Role, err
}

func (s *StateMachine) ValidActions(state []uint64, roles []Role, multiplier uint64) (map[string][]uint64, bool) {
	validActions := map[string][]uint64{}

	ok := false
	for a := range s.Transitions {
		action := string(a)
		outState, _, err := s.Transform(state, action, multiplier)
		if nil == err {
			ok = true
			var newState []uint64
			for _, val := range outState {
				newState = append(newState, uint64(val))
			}
			validActions[action] = newState
		}
	}

	return validActions, ok
}

var stateFormat = `
Initial:   {{ .Initial }}
Capacity:   {{ .Capacity }}
Transitions: {{ range $action, $txn := .Transitions }}
	{{ $action }} {{ printf "%v" $txn }}{{ end }}
`
var stateTemplate = template.Must(
	template.New("").Parse(stateFormat),
)

func (s StateMachine) String() string {
	b := &bytes.Buffer{}
	_ = stateTemplate.Execute(b, s)
	return b.String()
}

var vectorFormat = `
        Role:   {{.Role}}
        Delta:  {{.Delta}}
        Guards: {{ range $label, $g := .Guards }}
	            {{ $label }} {{ printf "%v" $g }} {{ end }}`

var vectorTemplate = template.Must(
	template.New("").Parse(vectorFormat),
)

func (t Transition) String() string {
	b := &bytes.Buffer{}
	_ = vectorTemplate.Execute(b, t)
	return b.String()
}
