package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func setupTest(args []string) (*bytes.Buffer, func()) {
	buf := &bytes.Buffer{}
	oldStderr := stderr
	oldArgs := os.Args

	stderr = buf
	os.Args = args

	return buf, func() {
		stderr = oldStderr
		os.Args = oldArgs
	}
}

var cliTests = map[string][]string{
	// Test that help is printed if no files are given
	"TranspileNoFilesHelp": {"test", "transpile"},

	// Test that help is printed if help flag is set, even if file is given
	"TranspileHelpFlag": {"test", "transpile", "-h"},

	// Test that help is printed if no files are given
	"AstNoFilesHelp": {"test", "ast"},

	// Test that help is printed if help flag is set, even if file is given
	"AstHelpFlag": {"test", "ast", "-h"},

	// Test that help is printed if no files are given
	"DebugNoFilesHelp": {"test", "debug"},

	// Test that help is printed if help flag is set, even if file is given
	"DebugHelpFlag": {"test", "debug", "-h"},

	"UnusedNoFilesHelp": {"test", "unused"},
	"UnusedHelpFlag":    {"test", "unused", "-h"},

	// Test that version is printed
	"Version": {"test", "version"},
}

func TestCLI(t *testing.T) {

	snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory("tests"))

	for testName, args := range cliTests {
		t.Run(testName, func(t *testing.T) {
			output, teardown := setupTest(args)
			defer teardown()

			runCommand()

			err := snapshotter.SnapshotMulti(testName, output)
			if err != nil {
				t.Fatalf("error: %s", err)
			}
		})
	}
}
