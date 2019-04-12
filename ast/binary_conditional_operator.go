package ast

// BinaryConditionalOperator is type of binary operator
type BinaryConditionalOperator struct {
	Addr       Address
	Pos        Position
	Type       string
	ChildNodes []Node
}

func parseBinaryConditionalOperator(line string) *BinaryConditionalOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type1>.*?)'`,
		line,
	)

	return &BinaryConditionalOperator{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type1"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BinaryConditionalOperator) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *BinaryConditionalOperator) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *BinaryConditionalOperator) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *BinaryConditionalOperator) Position() Position {
	return n.Pos
}
