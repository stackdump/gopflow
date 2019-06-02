package ptnet

import (
	"encoding/json"
	pFlow "github.com/stackdump/gopflow/pflow"
	. "github.com/stackdump/gopflow/statemachine"
)

type Place struct {
	Initial  uint64 `json:"initial"`
	Offset   int    `json:"offset"`
	Capacity uint64 `json:"capacity"`
}

type PTNet struct {
	Places      map[string]Place
	Transitions map[Action]Transition
}

func (p PTNet) getOffset(placeName string) (int, bool) {
	for placeID, place := range p.Places {
		if placeID == placeName {
			return place.Offset, true
		}
	}
	return -1, false
}

func (p PTNet) emptyState() []uint64 {
	var emptyState []uint64

	for x := 0; x < len(p.Places); x++ {
		emptyState = append(emptyState, uint64(0))
	}
	return emptyState
}

func (p PTNet) initialState() StateVector {
	if p.Places == nil || len(p.Places) == 0 {
		panic("no places defined")
	}
	init := p.emptyState()
	for _, place := range p.Places {
		init[place.Offset] = place.Initial
	}
	return StateVector(init[:])
}

func (p PTNet) capacityVector() StateVector {
	capacity := p.emptyState()
	for _, place := range p.Places {
		capacity[place.Offset] = place.Capacity
	}
	return StateVector(capacity[:])
}

func (p PTNet) StateMachine() *StateMachine {
	return &StateMachine{
		Initial:     p.initialState(),
		Capacity:    p.capacityVector(),
		Transitions: p.Transitions,
	}
}

func (p PTNet) String() string {
	s, err := json.MarshalIndent(p, "", "    ")

	if err != nil {
		panic("failed to serialize")
	}
	return string(s)
}

func LoadFile(filePath string) *PTNet {
	f, err := pFlow.LoadFile(filePath)
	if err != nil {
		panic("failed to load pflow")
	}
	pp := &pflowLoader{f: f}
	return pp.PTNet()
}
