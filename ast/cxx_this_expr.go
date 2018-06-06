package ast

// CXXThisExpr
type CXXThisExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	IsThis     bool
	ChildNodes []Node
}

func parseCXXThisExpr(line string) *CXXThisExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> 
		'(?P<type>.*?)'
		(?P<this> this)?
		`,
		line,
	)

	return &CXXThisExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		IsThis:     len(groups["this"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXThisExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXThisExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXThisExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXThisExpr) Position() Position {
	return n.Pos
}
