package program

func (p *Program) initializationStructs() {
	// stdlib.h
	p.Structs["div_t"] = &Struct{
		Name: "div_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "int",
			"rem":  "int",
		},
	}
	p.Structs["c4go_div_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.DivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// stdlib.h
	p.Structs["ldiv_t"] = &Struct{
		Name: "ldiv_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "long int",
			"rem":  "long int",
		},
	}
	p.Structs["c4go_ldiv_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.LdivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// stdlib.h
	p.Structs["lldiv_t"] = &Struct{
		Name: "lldiv_t",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "long long",
			"rem":  "long long",
		},
	}
	p.Structs["c4go_lldiv_t"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.LldivT",
		Type: StructType,
		Fields: map[string]interface{}{
			"quot": "Quot",
			"rem":  "Rem",
		},
	}

	// time.h
	p.Structs["struct tm"] = &Struct{
		Name: "struct tm",
		Type: StructType,
		Fields: map[string]interface{}{
			"tm_sec":   "int",
			"tm_min":   "int",
			"tm_hour":  "int",
			"tm_mday":  "int",
			"tm_mon":   "int",
			"tm_year":  "int",
			"tm_wday":  "int",
			"tm_yday":  "int",
			"tm_isdst": "int",
		},
	}
	p.Structs["c4go_struct tm"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Tm",
		Type: StructType,
		Fields: map[string]interface{}{
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

	// sys/time.h
	p.Structs["struct timeval"] = &Struct{
		Name: "struct timeval",
		Type: StructType,
		Fields: map[string]interface{}{
			"tv_sec":  "long",
			"tv_usec": "long",
		},
	}
	p.Structs["c4go_struct timeval"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Timeval",
		Type: StructType,
		Fields: map[string]interface{}{
			"tv_sec":  "TvSec",
			"tv_usec": "TvUsec",
		},
	}

	// sys/time.h
	p.Structs["struct timezone"] = &Struct{
		Name: "struct timezone",
		Type: StructType,
		Fields: map[string]interface{}{
			"tz_minuteswest": "long",
			"tz_dsttime":     "long",
		},
	}
	p.Structs["c4go_struct timezone"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Timezone",
		Type: StructType,
		Fields: map[string]interface{}{
			"tz_minuteswest": "TzMinuteswest",
			"tz_dsttime":     "TzDsttime",
		},
	}
}
