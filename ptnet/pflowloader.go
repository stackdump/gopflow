package ptnet

import (
	"fmt"

	pFlow "github.com/stackdump/gopflow/pflow"
	. "github.com/stackdump/gopflow/statemachine"
)

type pflowLoader struct {
	f                  *pFlow.PFlow
	placeElements      []pFlow.Place
	transitionElements []pFlow.Transition
	refPlaceElements   []pFlow.ReferencePlace
	arcElements        []pFlow.Arc
	placeIdOffset      map[int]int
	placeIdLabel       map[int]string
}

func (pfl *pflowLoader) emptyVector() Delta {
	d := Delta{}
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

func (pfl *pflowLoader) offset(placeId int) int {
	o, found := pfl.placeIdOffset[placeId]
	if !found {
		panic(fmt.Sprintf("unknownPlace %v", placeId))
	}
	return o
}

func (pfl *pflowLoader) placeLabel(placeId int) string {
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

func (pfl *pflowLoader) transitions() map[Action]Transition {
	roleMap := make(map[pFlow.TransitionId]Role)

	for _, r := range pfl.f.Roles {
		for _, txId := range r.TransitionIds {
			roleMap[txId] = Role(r.Name)
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
			t.place = pfl.placeLabel(a.SourceId)
			t.offset = pfl.offset(a.SourceId)
			txMap[a.DestinationId] = append(txMap[a.DestinationId], t)
		} else {
			t.unit = int64(a.Multiplicity)
			t.place = pfl.placeLabel(a.DestinationId)
			t.offset = pfl.offset(a.DestinationId)
			txMap[a.SourceId] = append(txMap[a.SourceId], t)
		}
	}

	l := make(map[Action]Transition)
	for _, t := range pfl.transitionElements {
		delta := pfl.emptyVector()
		guards := make(map[Condition]Delta, 0)

		for _, txn := range txMap[t.Id] {
			if txn.inhibitor {
				g := pfl.emptyVector()
				g[txn.offset] = txn.unit
				guards[Condition(txn.place)] = g
			} else {
				delta[txn.offset] = txn.unit
			}
		}
		action := Action(t.Label)
		l[action] = Transition{
			Delta:  delta,
			Role:   roleMap[pFlow.TransitionId(t.Id)],
			Guards: guards,
		}
	}

	return l
}

func (pfl *pflowLoader) loadSubnet(sub pFlow.SubNet) {
	for _, s := range sub.SubNets {
		pfl.loadSubnet(s)
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

func (pfl *pflowLoader) PTNet() *PTNet {
	for _, s := range pfl.f.SubNets {
		pfl.loadSubnet(s)
	}

	return &PTNet{
		pfl.places(),
		pfl.transitions(),
	}
}
