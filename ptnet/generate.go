package ptnet

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var sourceFileFormat = `package gen

import (
	. "github.com/stackdump/gopflow/statemachine"
)

var {{ .Name }} PetriNet = PetriNet{
	Places: map[string]Place { {{ range $key, $place := .Places}}
		"{{ $key }}": Place{
				Initial: {{ $place.Initial }},
				Offset: {{ $place.Offset }},
				Capacity: {{ $place.Capacity }},
		},{{ end }}
	},
	Transitions: map[Action]Transition { {{ range $action, $transition := .TransitionLiterals}}
		"{{ $action }}": Transition{ {{ $transition }} },{{ end }}
	},
}
`

func (s SourceFile) TransitionLiterals() map[string]string {
	out := map[string]string{}
	for action, t := range s.Transitions {
		var convert []string
		for _, v := range t {
			convert = append(convert, fmt.Sprintf("%v", v))
		}
		out[string(action)] = strings.Join(convert, ",")
	}
	return out
}

type SourceFile struct {
	PetriNet
	Name string
}

var sourceTemplate = template.Must(
	template.New("").Parse(sourceFileFormat),
)

func Generate(net PetriNet, name string) *bytes.Buffer {
	f := SourceFile{net, name}
	b := &bytes.Buffer{}
	err := sourceTemplate.Execute(b, f)
	if nil != err {
		panic(err)
	}
	return b
}
