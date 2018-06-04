package ast

// CXXConstructExpr is node constructor.
type CXXConstructExpr struct {
	Addr       Address
	Pos        Position
	Position2  string
	Type       string
	Type2      string
	ChildNodes []Node
}

func parseCXXConstructExpr(line string) *CXXConstructExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> 
		'(?P<type>.*?)'
		 '(?P<type2>.*?)'
		`,
		line,
	)

	return &CXXConstructExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type2:      groups["type2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXConstructExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXConstructExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXConstructExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXConstructExpr) Position() Position {
	return n.Pos
}
