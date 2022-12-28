package forgecli

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
	"github.com/sirupsen/logrus"
)

func TestFabricClientInstallerLatestFabricVersion(t *testing.T) {
	defer gock.Off()
	MockFabricXML(t)

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
	logrus.SetOutput(os.Stdout)
}

func TestFabricClientDownload(t *testing.T) {
	defer gock.Off()
	MockFabricXML(t)
	MockFabricJAR(t)

	var app appEnv
	app.hc = *http.DefaultClient
	app.version = "1.18.1"
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	logrus.Debug("TestFabricClientDownload")
	err := app.FabricClientDownload()
	st.Expect(t, err, nil)

	// Validates logging during the client download
	expected := "Downloading: https://maven.fabricmc.net"
	rawOutput := strings.Trim(buf.String(), "\n")
	logrus.Debug(rawOutput)
	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
	if !strings.Contains(output, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
	}

	// Removes downloaded file:
	if err := app.FabricClientRemoval(); err != nil {
		t.Errorf("Test failed, could not remove client jar, error:  '%s'", err)
	}
	logrus.SetOutput(os.Stdout)
}

func TestFabricClientInstaller(t *testing.T) {
	javaInstallFabric = func(installCommands []string) ([]byte, error) {
		return []byte("mock jar install"), nil
	}
	defer gock.Off()
	MockFabricXML(t)
	MockFabricJAR(t)

	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	var app appEnv
	if err := app.FabricClientInstaller(); err != nil {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
	}
	rawOutput := strings.Trim(buf.String(), "\n")
	t.Log(rawOutput)

	// Validates Installation of Jar
	expectedJarInstall := "Install Output: mock jar install"
	if !strings.Contains(rawOutput, expectedJarInstall) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedJarInstall, rawOutput)
	}

	// Validates remove of Jar
	expectedRemovalOutput := "Removing file: ./fabric-installer-0.11.1.jar"
	if !strings.Contains(rawOutput, expectedRemovalOutput) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedRemovalOutput, rawOutput)
	}
	logrus.SetOutput(os.Stdout)
}
