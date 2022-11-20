package main

import "testing"

func TestBinding(t *testing.T) {
	args := DefaultProgramArgs()
	args.inputFiles = []string{"tests/raylib/raygui.h"}
	args.outputFile = "./tests/bind.result.go"
	args.state = StateBinding
	args.verbose = false

	if err := Start(args); err != nil {
		t.Fatalf("Cannot transpile `%v`: %v", args, err)
	}
}
