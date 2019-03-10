package noarch

type Lconv struct {
	Currency_symbol []byte
	IntCurrSymbol   []byte
}

func Setlocale(category int32, locale []byte) []byte {
	return []byte("dd")
}

func Localeconv() []Lconv {
	var l Lconv
	l.Currency_symbol = []byte("fake data")
	l.IntCurrSymbol = []byte("fake data")
	return append([]Lconv{}, l)
}
