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

	// termios.h
	p.Structs["struct termios"] = &Struct{
		Name: "struct termios",
		Type: StructType,
		Fields: map[string]interface{}{
			"c_iflag": "unsigned int",
			"c_oflag": "unsigned int",
			"c_cflag": "unsigned int",
			"c_lflag": "unsigned int",
			"c_cc":    "unsigned char[32]",
		},
	}
	p.Structs["c4go_struct termios"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Termios",
		Type: StructType,
		Fields: map[string]interface{}{
			"c_iflag": "Iflag",
			"c_oflag": "Oflag",
			"c_cflag": "Cflag",
			"c_lflag": "Lflag",
			"c_cc":    "Cc",
		},
	}

	// sys/resource.h
	p.Structs["struct rusage"] = &Struct{
		Name: "struct rusage",
		Type: StructType,
		Fields: map[string]interface{}{
			"ru_utime":    "struct timeval",
			"ru_stime":    "struct timeval",
			"ru_maxrss":   "long",
			"ru_ixrss":    "long",
			"ru_idrss":    "long",
			"ru_isrss":    "long",
			"ru_minflt":   "long",
			"ru_majflt":   "long",
			"ru_nswap":    "long",
			"ru_inblock":  "long",
			"ru_oublock":  "long",
			"ru_msgsnd":   "long",
			"ru_msgrcv":   "long",
			"ru_nsignals": "long",
			"ru_nvcsw":    "long",
			"ru_nivcsw":   "long",
		},
	}
	p.Structs["c4go_struct rusage"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Rusage",
		Type: StructType,
		Fields: map[string]interface{}{
			"ru_utime":    "Utime",
			"ru_stime":    "Stime",
			"ru_maxrss":   "Maxrss",
			"ru_ixrss":    "Ixrss",
			"ru_idrss":    "Idrss",
			"ru_isrss":    "Isrss",
			"ru_minflt":   "Minflt",
			"ru_majflt":   "Majflt",
			"ru_nswap":    "Nswap",
			"ru_inblock":  "Inblock",
			"ru_oublock":  "Oublock",
			"ru_msgsnd":   "Msgsnd",
			"ru_msgrcv":   "Msgrcv",
			"ru_nsignals": "Nsignals",
			"ru_nvcsw":    "Nvcsw",
			"ru_nivcsw":   "Nivcsw",
		},
	}

	// sys/ioctl.h
	p.Structs["struct winsize"] = &Struct{
		Name: "struct winsize",
		Type: StructType,
		Fields: map[string]interface{}{
			"ws_row":    "unsigned short",
			"ws_col":    "unsigned short",
			"ws_xpixel": "unsigned short",
			"ws_ypixel": "unsigned short",
		},
	}
	p.Structs["c4go_struct winsize"] = &Struct{
		Name: "github.com/Konstantin8105/c4go/noarch.Winsize",
		Type: StructType,
		Fields: map[string]interface{}{
			"ws_row":    "Row",
			"ws_col":    "Col",
			"ws_xpixel": "Xpixel",
			"ws_ypixel": "Ypixel",
		},
	}
}
