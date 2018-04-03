package ast

// C4goErrorNode is error node type
type C4goErrorNode struct {
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n C4goErrorNode) AddChild(_ Node) {
	panic("Not acceptable to use for that node")
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n C4goErrorNode) Address() (a Address) {
	panic("Not acceptable to use for that node")
	return
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n C4goErrorNode) Children() (node []Node) {
	panic("Not acceptable to use for that node")
	return
}

// Position returns the position in the original source code.
func (n C4goErrorNode) Position() (p Position) {
	panic("Not acceptable to use for that node")
	return
}
