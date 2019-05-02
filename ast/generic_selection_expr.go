package ast

// GenericSelectionExpr is expression.
type GenericSelectionExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	IsLvalue   bool
	ChildNodes []Node
}

func parseGenericSelectionExpr(line string) *GenericSelectionExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<lvalue> lvalue)?`,
		line,
	)

	return &GenericSelectionExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		IsLvalue:   len(groups["lvalue"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *GenericSelectionExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *GenericSelectionExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *GenericSelectionExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *GenericSelectionExpr) Position() Position {
	return n.Pos
}
