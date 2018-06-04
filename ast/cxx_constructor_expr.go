package ast

// CXXConstructorExpr is node constructor.
type CXXConstructorExpr struct {
	Addr       Address
	Pos        Position
	Position2  string
	Type       string
	Type2      string
	ChildNodes []Node
}

func parseCXXConstructorExpr(line string) *CXXConstructorExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> 
		'(?P<type>.*?)'
		( '(?P<type2>.*?)')?
		`,
		line,
	)

	return &CXXConstructorExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type2:      groups["type2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXConstructorExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXConstructorExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXConstructorExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXConstructorExpr) Position() Position {
	return n.Pos
}
