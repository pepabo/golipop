package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/pepabo/golipop"
)

func TestVersion(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)

	cli := &CLI{outStream: out, errStream: err}
	args := strings.Split("lolp --version", " ")

	if status := cli.run(args); status != ExitOK {
		t.Fatalf("expected: \"%d\", actual: \"%d\"", ExitErr, status)
	}

	expected := fmt.Sprintf("lolp version %s", lolp.Version)
	if !strings.Contains(err.String(), expected) {
		t.Fatalf("expected: \"%s\", actual: \"%s\"", expected, err.String())
	}
}
