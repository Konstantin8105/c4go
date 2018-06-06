package ast

import (
	"testing"
)

func TestCXXRecord(t *testing.T) {
	nodes := map[string]Node{
		`0x34caec8 '__locale_struct'`: &CXXRecord{
			Addr:       0x34caec8,
			Type:       "__locale_struct",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
