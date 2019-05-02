package ast

// FunctionNoProtoType is function proto type
type FunctionNoProtoType struct {
	Addr       Address
	Type       string
	Kind       string
	ChildNodes []Node
}

func parseFunctionNoProtoType(line string) *FunctionNoProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)' (?P<kind>.*)",
		line,
	)

	return &FunctionNoProtoType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FunctionNoProtoType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FunctionNoProtoType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FunctionNoProtoType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FunctionNoProtoType) Position() Position {
	return Position{}
}
