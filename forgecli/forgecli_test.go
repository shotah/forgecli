package forgecli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadDotEnv() error {
	logrus.SetLevel(logrus.DebugLevel)
	if os.Getenv("FORGEKEY") == "" {
		if err := godotenv.Load("../.env"); err != nil {
			return err
		}
	}
	return nil
}

func TestCLIFromArgsReturnsError(t *testing.T) {
	expected := 2
	cliInput := []string{"-help"}
	actual := CLI(cliInput)
	if actual != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestCLIReturnsZero(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeModResponse(t)

	expected := 0
	cliInput := []string{"-projects", "416089"}
	actual := CLI(cliInput)
	if actual != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestValidateJavaInstallation(t *testing.T) {
	version := "Java Version Mock"
	javaVersion = []byte(version)
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	var app appEnv
	if err := app.ValidateJavaInstallation(); err != nil {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
	}

	expected := "java version found: Java Version Mock"
	rawOutput := strings.Trim(buf.String(), "\n")
	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
	if !strings.Contains(output, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
	}
}
