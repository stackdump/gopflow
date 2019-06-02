// Packaget ptnet provides a place-transition equivalent of an elementary petri-net
package statemachine

import (
	"bytes"
	"errors"
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
func (s *StateMachine) Transform(action string, multiplier uint64) (vectorOut []int64, role Role, err error) {

	t := s.Transitions[Action(action)]
	for offset, delta := range t.Delta {
		val := int64(s.State[offset]) + delta*int64(multiplier)
		vectorOut = append(vectorOut, val)
		if err == nil && val < 0 {
			err = errors.New("invalid output")
		}
		if err == nil && s.Capacity[offset] != 0 && val > int64(s.Capacity[offset]) {
			err = errors.New("exceeded capacity")
		}
	}
	return vectorOut, t.Role, err
}

func (s *StateMachine) ValidActions(multiplier uint64) (map[string][]uint64, bool) {
	validActions := map[string][]uint64{}

	ok := false
	for a := range s.Transitions {
		action := string(a)
		outState, _, err := s.Transform(action, multiplier)
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
	vectorOut, _, err := s.Transform(action, multiplier)

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
State:   {{ .State }}
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
