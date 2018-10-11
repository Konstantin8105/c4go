package ast

import (
	"strings"
)

// LinkageSpecDecl
type LinkageSpecDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	IsImplicit bool
	Name       string
	ChildNodes []Node
}

func parseLinkageSpecDecl(line string) *LinkageSpecDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<name> \w+?)?
		`,
		line,
	)

	return &LinkageSpecDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		IsImplicit: len(groups["implicit"]) > 0,
		Name:       strings.TrimSpace(groups["name"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *LinkageSpecDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *LinkageSpecDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *LinkageSpecDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *LinkageSpecDecl) Position() Position {
	return n.Pos
}
