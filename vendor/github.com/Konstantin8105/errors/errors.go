package errors

import "github.com/disiqueira/gotree"

// Tree is struct of error tree
type Tree struct {
	Name string
	errs []error
}

// New create a new tree error
func New(name string) *Tree {
	tr := new(Tree)
	tr.Name = name
	return tr
}

// Add error in tree node
func (e *Tree) Add(err error) *Tree {
	e.errs = append(e.errs, err)
	return e
}

// Error is typical function for interface error
func (e Tree) Error() (s string) {
	return e.getTree().Print()
}

// IsError check have errors in tree
func (e Tree) IsError() bool {
	return len(e.errs) > 0
}

func (e Tree) getTree() gotree.Tree {
	name := "+"
	if e.Name != "" {
		name = e.Name
	}
	t := gotree.New(name)
	for _, err := range e.errs {
		if et, ok := err.(Tree); ok {
			t.AddTree(et.getTree())
			continue
		}
		t.Add(err.Error())
	}
	return t
}

// Reset errors in tree
func (e *Tree) Reset() {
	e.errs = nil
}
