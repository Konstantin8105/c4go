package noarch

// Timeval - struct from <sys/time.h>
type Timeval struct {
	TvSec  int32
	TvUsec int32
}

// Timezone - struct from <sys/time.h>
type Timezone struct {
	TzMinuteswest int // minutes west of Greenwich
	TzDsttime     int // type of DST correction
}
