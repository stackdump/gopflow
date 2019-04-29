// Packaget ptnet provides a place-transition equivalent of an elementary petri-net
package statemachine

import (
	"bytes"
	"errors"
	"text/template"
)

type StateVector []uint64
type Transition []int64
type Action string

type StateMachine struct {
	Initial     StateVector
	Capacity    StateVector
	Transitions map[Action]Transition
	State       StateVector
}

func (s *StateMachine) Init() {
	for _, val := range s.Initial {
		s.State = append(s.State, val)
	}
}

func (s *StateMachine) Clone(state StateVector) StateMachine {
	return StateMachine{
		Initial:     s.Initial,
		Capacity:    s.Capacity,
		Transitions: s.Transitions,
		State:       state,
	}
}

// apply the transformation without overwriting state
func (s StateMachine) Transform(action string, multiplier uint64) ([]int64, error) {
	var vectorOut []int64
	var err error = nil

	for offset, delta := range s.Transitions[Action(action)] {
		val := int64(s.State[offset]) + delta*int64(multiplier)
		vectorOut = append(vectorOut, val)
		if err == nil && val < 0 {
			err = errors.New("invalid output")
		}
		if err == nil && s.Capacity[offset] != 0 && val > int64(s.Capacity[offset]) {
			err = errors.New("exceeded capacity")
		}
	}
	return vectorOut, err
}

func (s StateMachine) ValidActions(multiplier uint64) (map[string][]uint64, bool) {
	validActions := map[string][]uint64{}

	ok := false
	for a := range s.Transitions {
		action := string(a)
		outState, err := s.Transform(action, multiplier)
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

// apply the transformation and overwrite state
func (s *StateMachine) Commit(action string, multiplier uint64) ([]int64, error) {
	vectorOut, err := s.Transform(action, multiplier)

	if err == nil {
		for offset, val := range vectorOut {
			s.State[offset] = uint64(val)
		}
	}

	return vectorOut, err
}

var stateFormat = `
Initial:   {{ .Initial }}
Capacity:   {{ .Capacity }}
Transitions: {{ range $action, $txn := .Transitions }}
	{{ $action }} {{ printf "%v" $txn }}{{ end }}
State:   {{ .State }}
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
Vector:   {{ .Initial }}
`

var vectorTemplate = template.Must(
	template.New("").Parse(vectorFormat),
)

func (sv StateVector) String() string {
	b := &bytes.Buffer{}
	_ = vectorTemplate.Execute(b, sv)
	return b.String()
}
