package pflow

import (
	"encoding/xml"
	"io/ioutil"
)

type Place struct {
	Id       string `xml:"id"`
	X        int    `xml:"x"`
	Y        int    `xml:"y"`
	Label    string `xml:"label"`
	Tokens   int    `xml:"tokens"`
	IsStatic bool   `xml:"isStatic"`
}

type Transition struct {
	Id    string `xml:"id"`
	X     int    `xml:"x"`
	Y     int    `xml:"y"`
	Label string `xml:"label"`
}

type Arc struct {
	Type          string `xml:"type"`
	SourceId      int    `xml:"sourceId"`
	DestinationId int    `xml:"destinationId"`
	Multiplicity  int    `xml:"multiplicity"`
}

type ReferencePlace struct {
	Id               string `xml:"id"`
	X                int    `xml:"x"`
	Y                int    `xml:"y"`
	ConnectedPlaceId int    `xml:"connectedPlaceId"`
}

type ReferenceArc struct {
	PlaceId  int `xml:"placeId"`
	SubnetId int `xml:"subnetId"`
}

type SubNet struct {
	Id              string           `xml:"id"`
	X               int              `xml:"x"`
	Y               int              `xml:"y"`
	Label           string           `xml:"label"`
	Places          []Place          `xml:"place"`
	Transitions     []Transition     `xml:"transition"`
	Arcs            []Arc            `xml:"arc"`
	SubNets         []SubNet         `xml:"subnet"`
	ReferencePlaces []ReferencePlace `xml:"referencePlace"`
	ReferenceArcs   []ReferenceArc   `xml:"referenceArc"`
}

type TransitionId int

type Role struct {
	Id            string         `xml:"id"`
	Name          string         `xml:"name"`
	TransitionIds []TransitionId `xml:"transitionId"`
	CreateCase    bool           `xml:"createCase"`
	DestroyCase   bool           `xml:"destroyCase"`
}

type document struct {
	Id      string   `xml:"id"`
	X       int      `xml:"x"`
	Y       int      `xml:"y"`
	Label   string   `xml:"label"`
	SubNets []SubNet `xml:"subnet"`
	Roles   []Role   `xml:"roles>role"`
}

type PFlow struct {
	document
}

func LoadFile(path string) (*PFlow, error) {
	p := new(PFlow)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return p, err
	}
	_ = p.Unmarshal(data)
	return p, nil
}
func (p *PFlow) Document() document {
	return p.document
}

func (p *PFlow) Marshal() ([]byte, error) {
	return xml.Marshal(p.document)
}

func (p *PFlow) Unmarshal(data []byte) error {
	p0 := new(document)
	err := xml.Unmarshal(data, p0)
	if err != nil {
		return err
	}

	p.document = *p0
	return nil
}
