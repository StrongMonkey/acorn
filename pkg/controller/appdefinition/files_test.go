package appdefinition

import (
	"testing"

	"github.com/acorn-io/acorn/pkg/scheme"
	"github.com/acorn-io/baaah/pkg/router/tester"
)

func TestFileModes(t *testing.T) {
	tester.DefaultTest(t, scheme.Scheme, "testdata/files", DeploySpec)
}

func TestFileModesBug(t *testing.T) {
	tester.DefaultTest(t, scheme.Scheme, "testdata/files-bug", DeploySpec)
}

func TestInterpolation(t *testing.T) {
	tester.DefaultTest(t, scheme.Scheme, "testdata/interpolation", DeploySpec)
}
