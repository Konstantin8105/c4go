package ast

import (
	"github.com/Konstantin8105/c4go/util"
)

// CharacterLiteral is type of character literal
type CharacterLiteral struct {
	Addr       Address
	Pos        Position
	Type       string
	Value      int
	ChildNodes []Node
}

func parseCharacterLiteral(line string) *CharacterLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &CharacterLiteral{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Value:      util.Atoi(groups["value"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CharacterLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CharacterLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CharacterLiteral) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CharacterLiteral) Position() Position {
	return n.Pos
}

// CharacterLiteralError represents one instance of an error where the exact
// character value of a CharacterLiteral could not be determined from the
// original source. See RepairCharacterLiteralsFromSource for a full explanation.
type CharacterLiteralError struct {
	Node *CharacterLiteral
	Err  error
}
