// Package forgecli is main cli
package forgecli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// FabricAPIBaseURL Base Forge API URL
// Example: `https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.11.1/fabric-installer-0.11.1.jar`
const FabricAPIBaseURL = "https://maven.fabricmc.net/net/fabricmc/fabric-installer/"

func (app *appEnv) FabricClientInstallerVersion() error {
	if app.clientInstallerVersion != "" {
		return nil
	}
	var fabricXMLResponse XMLFabric
	if err := app.FetchXML(FabricMetadataURL, &fabricXMLResponse); err != nil {
		return fmt.Errorf("could not get fabric version from:\n%s", FabricMetadataURL)
	}
	app.clientInstallerVersion = fabricXMLResponse.Versioning.Latest
	return nil
}

func (app *appEnv) FabricClientDownload() error {
	logrus.Debug("FabricClientDownload")
	app.FabricClientInstallerVersion()
	app.clientInstallerFileName = fmt.Sprintf("fabric-installer-%s.jar", app.clientInstallerVersion)
	clientDownloadURL := FabricAPIBaseURL + app.clientInstallerVersion + "/" + app.clientInstallerFileName
	logrus.Debugf("FabricClientDownload: URL: %s", clientDownloadURL)
	// download the client where you are running the code from:
	err := app.FetchAndSave(clientDownloadURL, app.clientInstallerFileName, ".")
	return err
}

func (app *appEnv) FabricClientRemoval() error {
	// Validates file is downloaded
	filePath := "./" + app.clientInstallerFileName
	_, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return err
	}
	logrus.Debugf("Removing file: %s", filePath)
	return os.Remove(filePath)
}

func (app *appEnv) FabricClientInstaller() error {
	app.FabricClientDownload()
	installCommands := []string{"-jar", app.clientInstallerFileName, "client"}
	if app.version != "" {
		installCommands = []string{"-jar", app.clientInstallerFileName, "client", "-mcversion", app.version}
	}
	logrus.Debugf("java %v", installCommands)
	clientInstall, err := exec.Command("java", installCommands...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("fabric client install failed: %s", err)
	}
	logrus.Debugf("Install Output: %s", clientInstall)
	app.FabricClientRemoval()
	return nil
}
