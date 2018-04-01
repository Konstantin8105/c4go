package ast

import (
	"testing"
)

func TestVAArgExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff7d314bca8 <col:6, col:31> 'int *'`: &VAArgExpr{
			Addr:       0x7ff7d314bca8,
			Pos:        NewPositionFromString("col:6, col:31"),
			Type:       "int *",
			ChildNodes: []Node{},
		},

		`0x46f93e8 <col:28, col:58> 'LOGFUNC_t':'void (*)(void *, int, const char *)'`: &VAArgExpr{
			Addr:       0x46f93e8,
			Pos:        NewPositionFromString("col:28, col:58"),
			Type:       "LOGFUNC_t",
			Type2:      "void (*)(void *, int, const char *)",
			ChildNodes: []Node{},
		},
		`0x46f99b8 <col:30, col:64> 'sqlite3_int64':'long long'`: &VAArgExpr{
			Addr:       0x46f99b8,
			Pos:        NewPositionFromString("col:30, col:64"),
			Type:       "sqlite3_int64",
			Type2:      "long long",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
