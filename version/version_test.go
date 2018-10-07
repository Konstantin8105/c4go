package version

import (
	"reflect"
	"testing"

	"github.com/UtkarshGupta-CS/c4go/version"
)

func TestVersion(t *testing.T) {
	tempVersion := version.Version()

	versionDataType := reflect.TypeOf(tempVersion).Kind()

	if versionDataType != reflect.String {
		t.Errorf("Version computed is incorrect")
	}
}
