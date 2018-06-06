package ast

// CXXRecord struct
type CXXRecord struct {
	Addr       Address
	Type       string
	ChildNodes []Node
}

func parseCXXRecord(line string) *CXXRecord {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &CXXRecord{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXRecord) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXRecord) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXRecord) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXRecord) Position() Position {
	return Position{}
}
