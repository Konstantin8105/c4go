package ast

// CompoundLiteralExpr C99 6.5.2.5
type CompoundLiteralExpr struct {
	Addr       Address
	Pos        Position
	Type1      string
	Type2      string
	IsLvalue   bool
	ChildNodes []Node
}

func parseCompoundLiteralExpr(line string) *CompoundLiteralExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type1>.*?)'(:'(?P<type2>.*?)')?
		(?P<lvalue> lvalue)?
		`,
		line,
	)

	return &CompoundLiteralExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type1:      groups["type1"],
		Type2:      groups["type2"],
		IsLvalue:   len(groups["lvalue"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundLiteralExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CompoundLiteralExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CompoundLiteralExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CompoundLiteralExpr) Position() Position {
	return n.Pos
}
