package ast

import (
	"testing"
)

func TestCompoundAssignOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x2dc5758 <line:5:2, col:7> 'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'`: &CompoundAssignOperator{
			Addr:                  0x2dc5758,
			Pos:                   NewPositionFromString("line:5:2, col:7"),
			Type:                  "int",
			Opcode:                "+=",
			ComputationLHSType:    "int",
			ComputationResultType: "int",
			ChildNodes:            []Node{},
		},
		`0x2f27888 <line:1975:15, col:21> 'sqlite3_uint64':'unsigned long long' '>>=' ComputeLHSTy='sqlite3_uint64':'unsigned long long' ComputeResultTy='sqlite3_uint64':'unsigned long long'`: &CompoundAssignOperator{
			Addr:                   0x2f27888,
			Pos:                    NewPositionFromString("line:1975:15, col:21"),
			Type:                   "sqlite3_uint64",
			Type2:                  "unsigned long long",
			Opcode:                 ">>=",
			ComputationLHSType:     "sqlite3_uint64",
			ComputationLHSType2:    "unsigned long long",
			ComputationResultType:  "sqlite3_uint64",
			ComputationResultType2: "unsigned long long",
			ChildNodes:             []Node{},
		},
	}

	runNodeTests(t, nodes)
}
