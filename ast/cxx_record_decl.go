package ast

// CXXRecordDecl is node represents a record declaration.
type CXXRecordDecl struct {
	*RecordDecl
}

func parseCXXRecordDecl(line string) (res *CXXRecordDecl) {
	res = &CXXRecordDecl{parseRecordDecl(line)}
	return
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CXXRecordDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CXXRecordDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CXXRecordDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CXXRecordDecl) Position() Position {
	return n.Pos
}
