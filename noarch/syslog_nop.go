// +build windows nacl plan9
// A nop syslog logger for platforms with no syslog support

package noarch

// void    closelog(void);
func Closelog() {
}

// void    openlog(const char *, int, int);
func Openlog(ident *byte, logopt int, facility int) {
}

// int     setlogmask(int);
func Setlogmask(mask int) int {
	return 0
}

// void    syslog(int, const char *, ...);
func Syslog(priority int, format *byte, args ...interface{}) {
}

// void    vsyslog(int, const char *, struct __va_list_tag *);
func Vsyslog(priority int, format *byte, args VaList) {
}
