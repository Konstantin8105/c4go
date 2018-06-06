package preprocessor

import "testing"

func TestNewFilePPFail(t *testing.T) {
	_, err := NewFilePP([]string{""}, []string{""}, false)
	if err == nil {
		t.Fatalf("Haven`t error")
	}
}

func TestgetIncludeListFail(t *testing.T) {
	_, err := getIncludeList([]string{"@sdf s"}, []string{"wqq4 `?p"}, "w3 fdws", false)
	if err == nil {
		t.Fatalf("Haven`t error")
	}
}
