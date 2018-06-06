package ast

import (
	"strings"
)

// CXXConstructorDecl is node constructor.
type CXXConstructorDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	IsImplicit bool
	IsUsed     bool
	Type       string
	Type2      string
	IsInline   bool
	Other      string
	ChildNodes []Node
}

func parseCXXConstructorDecl(line string) *CXXConstructorDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<used> used)?
		( (?P<type>\w+)?)?
		( '(?P<type2>.*?)')?
		(?P<inline> inline)?
		(?P<other>.*)`,
		line,
	)

	return &CXXConstructorDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		IsImplicit: len(groups["implicit"]) > 0,
		IsUsed:     len(groups["used"]) > 0,
		Type:       groups["type"],
		Type2:      groups["type2"],
		IsInline:   len(groups["inline"]) > 0,
		Other:      groups["other"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXConstructorDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXConstructorDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXConstructorDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXConstructorDecl) Position() Position {
	return n.Pos
}
