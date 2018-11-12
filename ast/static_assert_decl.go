package ast

import "strings"

// StaticAssertDecl
type StaticAssertDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	ChildNodes []Node
}

func parseStaticAssertDecl(line string) *StaticAssertDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		`,
		line,
	)

	return &StaticAssertDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *StaticAssertDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *StaticAssertDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *StaticAssertDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *StaticAssertDecl) Position() Position {
	return n.Pos
}
