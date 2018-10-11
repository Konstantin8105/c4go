package program

// CStdStructType - conversion map from C standart library structures to
// c4go structures
var CStdStructType = map[string]string{
	// stdlib.h
	"div_t":   "github.com/Konstantin8105/c4go/noarch.DivT",
	"ldiv_t":  "github.com/Konstantin8105/c4go/noarch.LdivT",
	"lldiv_t": "github.com/Konstantin8105/c4go/noarch.LldivT",

	// time.h
	"tm":        "github.com/Konstantin8105/c4go/noarch.Tm",
	"struct tm": "github.com/Konstantin8105/c4go/noarch.Tm",
	"time_t":    "github.com/Konstantin8105/c4go/noarch.TimeT",

	"fpos_t": "int",
}

// This map is used to rename struct member names.
var StructFieldTranslations = map[string]map[string]string{
	// stdlib.h
	"div_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},
	"ldiv_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},
	"lldiv_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},

	// time.h
	"struct tm": {
		"tm_sec":   "TmSec",
		"tm_min":   "TmMin",
		"tm_hour":  "TmHour",
		"tm_mday":  "TmMday",
		"tm_mon":   "TmMon",
		"tm_year":  "TmYear",
		"tm_wday":  "TmWday",
		"tm_yday":  "TmYday",
		"tm_isdst": "TmIsdst",
	},
}
