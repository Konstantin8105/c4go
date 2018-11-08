package ast

import (
	"strings"
)

// CXXMethodDecl method of class
type CXXMethodDecl struct {
	Addr       Address
	Parent     string
	Prev       string
	Pos        Position
	Position2  string
	IsImplicit bool
	IsUsed     bool
	MethodName string
	Type       string
	IsInline   bool
	Other      string
	ChildNodes []Node
}

func parseCXXMethodDecl(line string) *CXXMethodDecl {
	groups := groupsFromRegex(
		`(?:parent (?P<parent>0x[0-9a-f]+) )?
		(?:prev (?P<prev>0x[0-9a-f]+) )?
		<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<used> used)?
		( (?P<name>\w+)?)?
		( '(?P<type>.*?)')?
		(?P<inline> inline)?
		(?P<other>.*)`,
		line,
	)

	return &CXXMethodDecl{
		Addr:       ParseAddress(groups["address"]),
		Prev:       groups["prev"],
		Parent:     groups["parent"],
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		IsImplicit: len(groups["implicit"]) > 0,
		IsUsed:     len(groups["used"]) > 0,
		Type:       groups["type"],
		MethodName: groups["name"],
		IsInline:   len(groups["inline"]) > 0,
		Other:      groups["other"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXMethodDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXMethodDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXMethodDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXMethodDecl) Position() Position {
	return n.Pos
}
