package main

import (
	"os"
	"testing"

	"github.com/shotah/forgecli/forgecli"
)

func TestHelp(t *testing.T) {
	expected := 2
	os.Args = []string{"-help"}
	actual := forgecli.CLI(os.Args[1:])
	if actual != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, actual)
	}
}
