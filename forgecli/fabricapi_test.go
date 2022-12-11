package forgecli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// TODO: MOCK THE CALLS!
func TestFabricClientDownload(t *testing.T) {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	var app appEnv
	if err := app.FabricClientDownload(); err != nil {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
	}

	// Validates logging during the client download
	expected := "Downloading: https://maven.fabricmc.net"
	rawOutput := strings.Trim(buf.String(), "\n")
	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
	if !strings.Contains(output, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
	}

	// Removes downloaded file:
	if err := app.FabricClientRemoval(); err != nil {
		t.Errorf("Test failed, could not remove client jar, error:  '%s'", err)
	}
}

// TODO: MOCK THE CALLS!
// BREAKING: CI/CD because of java call in github
// TESTS the full Version/Download/Install - Need to figure out mocks!
// func TestFabricClientInstaller(t *testing.T) {
// 	var buf bytes.Buffer
// 	logrus.SetOutput(&buf)
// 	var app appEnv
// 	if err := app.FabricClientInstaller(); err != nil {
// 		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
// 	}

// Validates logging during the client download
// 	expected := "Removing test file: ./fabric-installer"
// 	rawOutput := strings.Trim(buf.String(), "\n")
// 	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
// 	if !strings.Contains(output, expected) {
// 		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
// 	}
// }
