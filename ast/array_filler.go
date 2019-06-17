package ast

import (
	"strings"
)

// ArrayFiller is type of array filler
type ArrayFiller struct {
	ChildNodes []Node
}

var arrayFillerMarkers = []string{
	"array filler",
	"array_filler",
}

func parseArrayFiller(line string) *ArrayFiller {
	arrfill := &ArrayFiller{
		ChildNodes: []Node{},
	}

	for _, af := range arrayFillerMarkers {
		if strings.HasPrefix(line, af+":") {
			line = line[len(af+":")+1:]
		}
		if strings.HasPrefix(line, af) {
			line = line[len(af):]
		}
	}
	line = strings.TrimSpace(line)
	if line != "" {
		rn, err := Parse(line)
		if err != nil {
			panic(err)
		}
		arrfill.AddChild(rn)
	}

	return arrfill
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ArrayFiller) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. For an ArrayFilter this will
// always be zero. See the documentation for the Address type for more
// information.
func (n *ArrayFiller) Address() Address {
	return 0
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ArrayFiller) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ArrayFiller) Position() Position {
	return Position{}
}
