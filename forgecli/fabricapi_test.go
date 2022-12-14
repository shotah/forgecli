package forgecli

import (
	"bytes"
	"net/http"
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
}

// TODO: Figure out how to mock 'java' call
// BREAKS: on CI/CD in Github
// func TestFabricClientInstaller(t *testing.T) {
// 	defer gock.Off()
// 	MockFabricXML(t)
// 	MockFabricJAR(t)

// 	var buf bytes.Buffer
// 	logrus.SetOutput(&buf)
// 	var app appEnv
// 	if err := app.FabricClientInstaller(); err != nil {
// 		t.Errorf("Test failed, expected: '%s', got:  '%s'", "nil", err)
// 	}

// 	// Validates logging during the client download
// 	expected := "Removing file: ./fabric-installer-0.11.1.jar"
// 	rawOutput := strings.Trim(buf.String(), "\n")
// 	output := rawOutput[strings.LastIndex(rawOutput, "=")+1:]
// 	if !strings.Contains(output, expected) {
// 		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, output)
// 	}
// }
