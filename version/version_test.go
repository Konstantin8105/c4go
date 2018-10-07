package version

import (
	"reflect"
	"testing"

	"github.com/Konstantin8105/c4go/version"
)

func TestVersion(t *testing.T) {
	tempVersion := version.Version()

	versionDataType := reflect.TypeOf(tempVersion).Kind()

	if versionDataType != reflect.String {
		t.Errorf("Version computed is incorrect")
	}
}
