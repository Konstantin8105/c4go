package ast

import (
	"testing"
)

func TestDecayedType(t *testing.T) {
	nodes := map[string]Node{
		`0x1dfca42 'struct __va_list_tag *' sugar`: &DecayedType{
			Addr:       0x1dfca42,
			Type:       "struct __va_list_tag *",
			Tags:       "sugar",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
