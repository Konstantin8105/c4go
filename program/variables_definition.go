package program

// CVariables is map of convertion from C var to C4go variable
var CVariables = map[string]string{
	// stdio.h
	"stdin":  "github.com/Konstantin8105/c4go/noarch.Stdin",
	"stdout": "github.com/Konstantin8105/c4go/noarch.Stdout",
	"stderr": "github.com/Konstantin8105/c4go/noarch.Stderr",

	// ctype.h
	"_ISupper":  "github.com/Konstantin8105/c4go/noarch.ISupper",
	"_ISlower":  "github.com/Konstantin8105/c4go/noarch.ISlower",
	"_ISalpha":  "github.com/Konstantin8105/c4go/noarch.ISalpha",
	"_ISdigit":  "github.com/Konstantin8105/c4go/noarch.ISdigit",
	"_ISxdigit": "github.com/Konstantin8105/c4go/noarch.ISxdigit",
	"_ISspace":  "github.com/Konstantin8105/c4go/noarch.ISspace",
	"_ISprint":  "github.com/Konstantin8105/c4go/noarch.ISprint",
	"_ISgraph":  "github.com/Konstantin8105/c4go/noarch.ISgraph",
	"_ISblank":  "github.com/Konstantin8105/c4go/noarch.ISblank",
	"_IScntrl":  "github.com/Konstantin8105/c4go/noarch.IScntrl",
	"_ISpunct":  "github.com/Konstantin8105/c4go/noarch.ISpunct",
	"_ISalnum":  "github.com/Konstantin8105/c4go/noarch.ISalnum",
}
