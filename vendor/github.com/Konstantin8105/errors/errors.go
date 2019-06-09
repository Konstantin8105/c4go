package errors

import (
	"github.com/Konstantin8105/tree"
)

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
	if err == (error)(nil) {
		return nil
	}
	if e == (*Tree)(nil) {
		return nil
	}
	if et, ok := err.(Tree); ok {
		err = &et
	}
	e.errs = append(e.errs, err)
	return e
}

// Error is typical function for interface error
func (e Tree) Error() (s string) {
	return e.getTree().String()
}

// IsError check have errors in tree
func (e Tree) IsError() bool {
	return len(e.errs) > 0
}

func (e Tree) getTree() *tree.Tree {
	name := "+"
	if e.Name != "" {
		name = e.Name
	}
	t := tree.New(name)
	for _, err := range e.errs {
		if et, ok := err.(*Tree); ok {
			t.Add(et.getTree())
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

// Walk walking by error tree
func Walk(t *Tree, f func(error)) {
	if t == (*Tree)(nil) {
		return
	}
	for i := range t.errs {
		if et, ok := t.errs[i].(*Tree); ok {
			Walk(et, f)
			continue
		}
		f(t.errs[i])
	}
}
