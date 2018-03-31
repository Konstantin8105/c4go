package ast

// CompoundAssignOperator is type of compound assign operator
type CompoundAssignOperator struct {
	Addr                   Address
	Pos                    Position
	Type                   string
	Type2                  string
	Opcode                 string
	ComputationLHSType     string
	ComputationLHSType2    string
	ComputationResultType  string
	ComputationResultType2 string
	ChildNodes             []Node
}

func parseCompoundAssignOperator(line string) *CompoundAssignOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.+?)'(:'(?P<type2>.*)')?
		 '(?P<opcode>.+?)'
		 ComputeLHSTy='(?P<clhstype>.+?)'(:'(?P<clhstype2>.*)')?
		 ComputeResultTy='(?P<crestype>.+?)'(:'(?P<crestype2>.*)')?`,
		line,
	)

	return &CompoundAssignOperator{
		Addr:                   ParseAddress(groups["address"]),
		Pos:                    NewPositionFromString(groups["position"]),
		Type:                   groups["type"],
		Type2:                  groups["type2"],
		Opcode:                 groups["opcode"],
		ComputationLHSType:     groups["clhstype"],
		ComputationLHSType2:    groups["clhstype2"],
		ComputationResultType:  groups["crestype"],
		ComputationResultType2: groups["crestype2"],
		ChildNodes:             []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundAssignOperator) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CompoundAssignOperator) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CompoundAssignOperator) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CompoundAssignOperator) Position() Position {
	return n.Pos
}
