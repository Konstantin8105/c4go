package ast

import (
	"reflect"
	"testing"

	"github.com/Konstantin8105/c4go/util"
)

func TestArrayFiller(t *testing.T) {
	expected := &ArrayFiller{
		ChildNodes: []Node{},
	}
	actual, err := Parse(`array filler`)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s", util.ShowDiff(formatMultiLine(expected),
			formatMultiLine(actual)))
	}
	if err != nil {
		t.Errorf("Error parsing")
	}

	if uint64(actual.Address()) != 0 {
		t.Fatal("Address is not zero")
	}
	var v ArrayFiller
	actual.AddChild(&v)
	if len(actual.Children()) == 0 {
		t.Fatal("Childrens is not correct")
	}

	_ = actual.Position()
}
