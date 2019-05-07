package ptnet

import (
	stateMachine "github.com/stackdump/goflow/statemachine"
	pFlow "github.com/stackdump/gopflow/pflow"
)

type Place struct {
	Initial  uint64 `json:"initial"`
	Offset   int    `json:"offset"`
	Capacity uint64 `json:"capacity"`
}

type PetriNet struct {
	Places      map[string]Place
	Transitions map[stateMachine.Action]stateMachine.Transition
}

func (p PetriNet) getOffset(placeName string) (int, bool) {
	for placeID, place := range p.Places {
		if placeID == placeName {
			return place.Offset, true
		}
	}
	return -1, false
}

func getWeight(a pFlow.Arc) int64 {
	return int64(a.Multiplicity)
}

/* FIXME
func GetInitialValue(m pFlow.InitialMarking) uint64 {
	tokenVals := strings.Split(m.TokenValueCsv, ",")
	val, err := strconv.ParseInt(tokenVals[1], 10, 64)

	if err != nil || tokenVals[0] != "Default" {
		panic("Error Parsing InitialMarking")
	}
	return uint64(val)
}
*/

func (p PetriNet) GetEmptyVector() []int64 {
	var emptyVector []int64

	for x := 0; x < len(p.Places); x++ {
		emptyVector = append(emptyVector, int64(0))
	}
	return emptyVector
}

func (p PetriNet) GetEmptyState() []uint64 {
	var emptyState []uint64

	for x := 0; x < len(p.Places); x++ {
		emptyState = append(emptyState, uint64(0))
	}
	return emptyState
}

func (p PetriNet) GetInitialState() stateMachine.StateVector {
	if p.Places == nil || len(p.Places) == 0 {
		panic("no places defined")
	}
	init := p.GetEmptyState()
	for _, place := range p.Places {
		init[place.Offset] = place.Initial
	}
	return stateMachine.StateVector(init[:])
}

func (p PetriNet) GetCapacityVector() stateMachine.StateVector {
	capacity := p.GetEmptyState()
	for _, place := range p.Places {
		capacity[place.Offset] = place.Capacity
	}
	return stateMachine.StateVector(capacity[:])
}

func (p PetriNet) StateMachine() stateMachine.StateMachine {
	return stateMachine.StateMachine{
		Initial:     p.GetInitialState(),
		Capacity:    p.GetCapacityVector(),
		Transitions: p.Transitions,
		State:       stateMachine.StateVector{},
	}
}

func LoadFile(pflowPath string) (*pFlow.PFlow, error) {
	// FIXME: return PetriNet instead of PFlow
	return pFlow.LoadFile(pflowPath)
}

/*
func LoadPnmlFromFile(path string) (PetriNet, stateMachine.StateMachine) {
	pFlowDef, _ := pFlow.LoadFile(path)

	petriNet := PetriNet{
		map[string]Place{},
		map[stateMachine.Action]stateMachine.Transition{},
		pFlowDef,
	}

	if len(pFlowDef.Nets) != 1 {
		panic("Expect a single petri-net definition per file")
	}

	net := pFlowDef.Nets[0]

	for offset, p := range net.Places {
		petriNet.Places[p.Id] =
			Place{
				Initial:  GetInitialValue(p.InitialMarking),
				Offset:   offset,
				Capacity: p.Capacity.Value,
			}
	}

	for _, txn := range net.Transitions {
		petriNet.Transitions[stateMachine.Action(txn.Id)] = petriNet.GetEmptyVector()
	}

	for _, arc := range net.Arcs {
		var action string
		var sign int64
		var offset int

		targetOffset, targetIsPlace := petriNet.getOffset(arc.Target)
		sourceOffset, sourceIsPlace := petriNet.getOffset(arc.Source)

		if sourceIsPlace {
			action = arc.Target
			offset = sourceOffset
			sign = -1
		}

		if targetIsPlace {
			action = arc.Source
			offset = targetOffset
			sign = 1
		}

		petriNet.Transitions[stateMachine.Action(action)][offset] += sign * getWeight(arc)
	}

	return petriNet, petriNet.StateMachine()
}
*/
