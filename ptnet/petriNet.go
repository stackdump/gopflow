package ptnet

import (
	"encoding/json"
	"fmt"
	pFlow "github.com/stackdump/gopflow/pflow"
	stateMachine "github.com/stackdump/gopflow/statemachine"
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

func (p PetriNet) StateMachine() *stateMachine.StateMachine {
	return &stateMachine.StateMachine{
		Initial:     p.GetInitialState(),
		Capacity:    p.GetCapacityVector(),
		Transitions: p.Transitions,
		State:       stateMachine.StateVector{},
	}
}

func (p PetriNet) String() string {
	s, err := json.MarshalIndent(p, "", "    ")

	if err != nil {
		panic("failed to serialize")
	}
	return string(s)
}

type pflowLoader struct {
	f                  *pFlow.PFlow
	placeElements      []pFlow.Place
	transitionElements []pFlow.Transition
	refPlaceElements   []pFlow.ReferencePlace
	arcElements        []pFlow.Arc
	placeIdOffset      map[int]int
	placeIdLabel       map[int]string
}

func (pfl *pflowLoader) emptyVector() stateMachine.Delta {
	d := stateMachine.Delta{}
	for range pfl.placeElements {
		d = append(d, 0)
	}
	return d
}

func (pfl *pflowLoader) places() map[string]Place {
	pfl.placeIdOffset = make(map[int]int)
	pfl.placeIdLabel = make(map[int]string)

	l := make(map[string]Place)
	for i, p := range pfl.placeElements {
		l[p.Label] = Place{
			p.Tokens,
			i,
			p.Capacity,
		}
		pfl.placeIdOffset[p.Id] = i
		pfl.placeIdLabel[p.Id] = p.Label
	}

	for _, rp := range pfl.refPlaceElements {
		// index connected places
		pfl.placeIdOffset[rp.Id] = pfl.placeIdOffset[rp.ConnectedPlaceId]
		pfl.placeIdLabel[rp.Id] = pfl.placeIdLabel[rp.ConnectedPlaceId]
	}

	return l
}

func (pfl *pflowLoader) getOffset(placeId int) int {
	o, found := pfl.placeIdOffset[placeId]
	if !found {
		panic(fmt.Sprintf("unknownPlace %v", placeId))
	}
	return o
}

func (pfl *pflowLoader) getPlaceLabel(placeId int) string {
	label, found := pfl.placeIdLabel[placeId]
	if !found {
		panic(fmt.Sprintf("unknownPlace %v", placeId))
	}
	return label
}

func (pfl *pflowLoader) isPlace(placeId int) bool {
	_, found := pfl.placeIdOffset[placeId]
	return found
}

func (pfl *pflowLoader) transitions() map[stateMachine.Action]stateMachine.Transition {
	roleMap := make(map[pFlow.TransitionId]stateMachine.Role)

	for _, r := range pfl.f.Roles {
		for _, txId := range r.TransitionIds {
			roleMap[txId] = stateMachine.Role(r.Name)
		}
	}

	type txBuilder struct {
		unit      int64
		inhibitor bool
		place     string
		offset    int
	}

	txMap := make(map[int][]txBuilder)

	for _, a := range pfl.arcElements {
		t := txBuilder{inhibitor: a.Type == "inhibitor"}

		if pfl.isPlace(a.SourceId) {
			t.unit = int64(-1 * a.Multiplicity)
			t.place = pfl.getPlaceLabel(a.SourceId)
			t.offset = pfl.getOffset(a.SourceId)
			txMap[a.DestinationId] = append(txMap[a.DestinationId], t)
		} else {
			t.unit = int64(a.Multiplicity)
			t.place = pfl.getPlaceLabel(a.DestinationId)
			t.offset = pfl.getOffset(a.DestinationId)
			txMap[a.SourceId] = append(txMap[a.SourceId], t)
		}
	}

	l := make(map[stateMachine.Action]stateMachine.Transition)
	for _, t := range pfl.transitionElements {
		delta := pfl.emptyVector()
		guards := make(map[stateMachine.Condition]stateMachine.Delta, 0)

		for _, txn := range txMap[t.Id] {
			if txn.inhibitor {
				g := pfl.emptyVector()
				g[txn.offset] = txn.unit
				guards[stateMachine.Condition(txn.place)] = g
			} else {
				delta[txn.offset] = txn.unit
			}
		}
		action := stateMachine.Action(t.Label)
		l[action] = stateMachine.Transition{
			Delta:  delta,
			Role:   roleMap[pFlow.TransitionId(t.Id)],
			Guards: guards,
		}
	}

	return l
}

func (pfl *pflowLoader) loadSubnets(sub pFlow.SubNet) {
	for _, s := range sub.SubNets {
		pfl.loadSubnets(s)
	}

	for _, p := range sub.Places {
		pfl.placeElements = append(pfl.placeElements, p)
	}

	for _, t := range sub.Transitions {
		pfl.transitionElements = append(pfl.transitionElements, t)
	}

	for _, r := range sub.ReferencePlaces {
		pfl.refPlaceElements = append(pfl.refPlaceElements, r)
	}

	for _, a := range sub.Arcs {
		pfl.arcElements = append(pfl.arcElements, a)
	}
}

func (pfl *pflowLoader) toNet() (p *PetriNet, err error) {
	for _, s := range pfl.f.SubNets {
		pfl.loadSubnets(s)
	}

	n := PetriNet{
		pfl.places(),
		pfl.transitions(),
	}

	return &n, err
}

func LoadFile(pflowPath string) (p *PetriNet, err error) {
	f, err := pFlow.LoadFile(pflowPath)
	if err != nil {
		panic("failed to load pflow")
	}
	pp := &pflowLoader{f: f}
	return pp.toNet()
}
