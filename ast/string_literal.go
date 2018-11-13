package ast

// StringLiteral is type of string literal
type StringLiteral struct {
	Addr       Address
	Pos        Position
	Type       string
	Value      string
	IsLvalue   bool
	ChildNodes []Node
}

func parseStringLiteral(line string) *StringLiteral {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<value>".*")`,
		line,
	)

	return &StringLiteral{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Value:      unquote(groups["value"]),
		IsLvalue:   true,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *StringLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *StringLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *StringLiteral) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *StringLiteral) Position() Position {
	return n.Pos
}
