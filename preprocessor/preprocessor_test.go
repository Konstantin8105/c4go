package preprocessor

import "testing"

func TestNewFilePPFail(t *testing.T) {
	_, err := NewFilePP([]string{""}, []string{""}, false)
	if err == nil {
		t.Fatalf("Haven`t error")
	}
}

func TestGetIncludeListFail(t *testing.T) {
	_, err := getIncludeList([]string{"@sdf s"}, []string{"wqq4 `?p"}, []string{"w3 fdws", "sdfsr 4"}, false)
	if err == nil {
		t.Fatalf("Haven`t error")
	}
}
