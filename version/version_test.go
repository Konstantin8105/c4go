package version

import "testing"

func TestVersion(t *testing.T) {

	tempVersion := Version()

	GitSHA = "test"
	tempVersion2 := Version()

	t.Log(tempVersion)
	t.Log(tempVersion2)

	if tempVersion == tempVersion2 {
		t.Errorf("Version has not changed with different Git hash.")
	}
}
