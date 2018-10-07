package version

import (
	"testing"

	"github.com/Konstantin8105/c4go/version"
)

func TestVersion(t *testing.T) {

	tempVersion := version.Version()

	version.GitSHA = "test"
	tempVersion2 := version.Version()

	t.Log(tempVersion)
	t.Log(tempVersion2)

	if tempVersion == tempVersion2 {
		t.Errorf("Version has not changed with different Git hash.")
	}
}
