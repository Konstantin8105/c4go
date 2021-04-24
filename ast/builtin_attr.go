package ast

import "strings"

// BuiltinAttr is a type of attribute ...
type BuiltinAttr struct {
	Addr        Address
	Pos         Position
	IsImplicit  bool
	IsInherited bool
	Name        string
	IsAligned   bool
	ChildNodes  []Node
}

func parseBuiltinAttr(line string) *BuiltinAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<implicit> Implicit)?
		(?P<name> .*)?
		`,
		line,
	)

	return &BuiltinAttr{
		Addr:        ParseAddress(groups["address"]),
		Pos:         NewPositionFromString(groups["position"]),
		IsImplicit:  len(groups["implicit"]) > 0,
		IsInherited: len(groups["inherited"]) > 0,
		Name:        strings.TrimSpace(groups["name"]),
		ChildNodes:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BuiltinAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *BuiltinAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *BuiltinAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *BuiltinAttr) Position() Position {
	return n.Pos
}
