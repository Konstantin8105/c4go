// Package tree create and print tree.
package tree

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	middleItem   = "├──"
	continueItem = "│  "
	emptyItem    = "   "
	lastItem     = "└──"
	nullNode     = "<< NULL >>"
)

// Tree struct of tree
type Tree struct {
	Name  string
	nodes []*Tree
}

// New returns a new tree
func New(name string) (tr *Tree) {
	tr = new(Tree)
	tr.Name = name
	return tr
}

// Add node in tree
func (t *Tree) Add(node interface{}) {
	if tr, ok := node.(Tree); ok {
		node = &tr
	}
	if tr, ok := node.(*Tree); ok {
		if tr != (*Tree)(nil) {
			t.nodes = append(t.nodes, tr)
			return
		}
		node = nullNode
	}
	n := new(Tree)
	n.Name = toString(node)
	t.nodes = append(t.nodes, n)
	return
}

// String return string with tree view
func (t Tree) String() (out string) {
	return t.printNode(false, []string{})
}

func toString(i interface{}) (out string) {
	out = nullNode
	if i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()) {
		return
	}

	switch v := i.(type) {
	case interface {
		String() string
	}:
		if v != nil {
			out = v.String()
		}

	case string:
		if v != "" {
			out = v
		}

	case interface {
		Error() string
	}:
		if v != nil {
			out = v.Error()
		}

	default:
		if i != nil {
			out = fmt.Sprintf("%v", i)
		}
	}

	return
}

func (t Tree) printNode(isLast bool, spaces []string) (out string) {
	// clean name from spaces at begin and end of string
	var name string
	if t.Name == "" {
		name = nullNode
	} else {
		name = strings.TrimSpace(t.Name)
	}

	// split name into strings lines
	lines := strings.Split(name, "\n")

	var tab [2]string
	for i, level := 0, len(spaces); i < level; i++ {
		if i > 0 {
			tab[0] += spaces[i]
			tab[1] += spaces[i]
		}
		if i == level-1 {
			if isLast {
				tab[0] += lastItem
				tab[1] += emptyItem
			} else {
				tab[0] += middleItem
				tab[1] += continueItem
			}
		}
	}

	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
		if i == 0 {
			out += tab[0] + lines[i]
		} else {
			out += tab[1] + lines[i]
		}
		out += "\n"
	}

	size := len(spaces)
	if isLast {
		spaces = append(spaces, emptyItem)
	} else {
		spaces = append(spaces, continueItem)
	}
	defer func() {
		spaces = spaces[:size]
	}()

	for i := 0; i < len(t.nodes); i++ {
		out += t.nodes[i].printNode(i == len(t.nodes)-1, spaces)
	}

	return
}
