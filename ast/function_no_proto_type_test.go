package ast

import (
	"testing"
)

func TestFunctionNoProtoType(t *testing.T) {
	nodes := map[string]Node{
		`0x3e48580 'struct S *()' cdecl`: &FunctionNoProtoType{
			Addr:       0x3e48580,
			Type:       "struct S *()",
			Kind:       "cdecl",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
