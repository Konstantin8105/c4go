package ast

// GNUInlineAttr
type GNUInlineAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseGNUInlineAttr(line string) *GNUInlineAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>`,
		line,
	)

	return &GNUInlineAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *GNUInlineAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *GNUInlineAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *GNUInlineAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *GNUInlineAttr) Position() Position {
	return n.Pos
}
