package program

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
)

// Struct represents the definition for a C struct.
type Struct struct {
	// The name of the struct.
	Name string

	// True if the struct kind is an union.
	// This field is used to avoid to dupplicate code for union case the type is the same.
	// Plus, this field is used in collaboration with the method "c4go/program".*Program.GetStruct()
	IsUnion bool

	// Each of the fields and their C type. The field may be a string or an
	// instance of Struct for nested structures.
	Fields map[string]interface{}
}

// NewStruct creates a new Struct definition from an ast.RecordDecl.
func NewStruct(n *ast.RecordDecl) (_ *Struct, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot create new structure : %T. Error : %v",
				n, err)
		}
	}()
	fields := make(map[string]interface{})

	for _, field := range n.Children() {
		switch f := field.(type) {
		case *ast.FieldDecl:
			fields[f.Name] = f.Type

		case *ast.IndirectFieldDecl:
			fields[f.Name] = f.Type

		case *ast.RecordDecl:
			fields[f.Name], err = NewStruct(f)
			if err != nil {
				return
			}

		case *ast.MaxFieldAlignmentAttr,
			*ast.AlignedAttr,
			*ast.PackedAttr,
			*ast.EnumDecl,
			*ast.TransparentUnionAttr,
			*ast.FullComment:
			// FIXME: Should these really be ignored?

		default:
			err = fmt.Errorf("cannot decode: %#v", f)
			return
		}
	}

	return &Struct{
		Name:    n.Name,
		IsUnion: n.Kind == "union",
		Fields:  fields,
	}, nil
}

// IsUnion - return true if the cType is 'union' or
// typedef of union
func (p *Program) IsUnion(cType string) bool {
	if strings.HasPrefix(cType, "union ") {
		return true
	}
	if _, ok := p.Unions[cType]; ok {
		return true
	}
	if _, ok := p.Unions["union "+cType]; ok {
		return true
	}
	if _, ok := p.GetBaseTypeOfTypedef("union " + cType); ok {
		return true
	}
	if t, ok := p.GetBaseTypeOfTypedef(cType); ok {
		if t == cType {
			panic(fmt.Errorf("Cannot be same name: %s", t))
		}
		if strings.HasPrefix(t, "struct ") {
			return false
		}
		if t == "" {
			panic(fmt.Errorf("Type cannot be empty"))
		}
		return p.IsUnion(t)
	}
	return false
}

// GetBaseTypeOfTypedef - return typedef type
func (p *Program) GetBaseTypeOfTypedef(cTypedef string) (
	cBase string, ok bool) {

	cBase, ok = p.TypedefType[cTypedef]
	if cBase == "" && ok {
		panic(fmt.Errorf("Type cannot be empty"))
	}

	return
}
