package program

// DefinitionVariable is map of convertion from C var to C4go variable
var DefinitionVariable = map[string]string{
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

	// termios.h
	"TCSANOW":   "github.com/pkg/term/termios.TCSANOW",
	"TCSADRAIN": "github.com/pkg/term/termios.TCSADRAIN",
	"TCSAFLUSH": "github.com/pkg/term/termios.TCSAFLUSH",
	//
	"TCIFLUSH":  "github.com/pkg/term/termios.TCIFLUSH",
	"TCOFLUSH":  "github.com/pkg/term/termios.TCOFLUSH",
	"TCIOFLUSH": "github.com/pkg/term/termios.TCIOFLUSH",
	//
	"TCSETS":  "github.com/pkg/term/termios.TCSETS",
	"TCSETSW": "github.com/pkg/term/termios.TCSETSW",
	"TCSETSF": "github.com/pkg/term/termios.TCSETSF",
	"TCFLSH":  "github.com/pkg/term/termios.TCFLSH",
	"TCSBRK":  "github.com/pkg/term/termios.TCSBRK",
	"TCSBRKP": "github.com/pkg/term/termios.TCSBRKP",
	//
	"IXON":    "github.com/pkg/term/termios.IXON",
	"IXANY":   "github.com/pkg/term/termios.IXANY",
	"IXOFF":   "github.com/pkg/term/termios.IXOFF",
	"CRTSCTS": "github.com/pkg/term/termios.CRTSCTS",

	// sys/ioctl.h
	"TIOCGWINSZ": "golang.org/x/sys/unix.TIOCGWINSZ",
}
