package forgecli

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadDotEnv() {
	logrus.SetLevel(logrus.DebugLevel)
	if os.Getenv("FORGEKEY") == "" {
		err := godotenv.Load("../.env")
		check(err)
	}
}

func TestCLIReturnsError(t *testing.T) {
	expected := 2
	cliInput := []string{"-help"}
	actual := CLI(cliInput)
	if actual != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestFabricClientInstallerLatestFabricVersion(t *testing.T) {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	var app appEnv
	app.hc = *http.DefaultClient
	expected := "Fetching XML: https://maven.fabricmc.net"
	app.FabricClientInstallerVersion()
	rawOutput := strings.Trim(buf.String(), "\n")
	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
	if !strings.Contains(output, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
	}
}

func TestValidateJavaInstallation(t *testing.T) {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	var app appEnv
	if err := app.ValidateJavaInstallation(); err != nil {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
	}
	expected := "java version found"
	rawOutput := strings.Trim(buf.String(), "\n")
	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
	if !strings.Contains(output, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
	}
}
