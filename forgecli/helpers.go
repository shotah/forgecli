package forgecli

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

func (app *appEnv) GetTargetDirectory() error {
	if app.destination != "" {
		return nil
	}

	user, err := user.Current()
	if err != nil {
		return err
	}
	os := runtime.GOOS
	switch os {
	case "windows":
		app.destination = fmt.Sprintf("%s\\AppData\\Roaming\\.minecraft\\mods", user.HomeDir)
	case "darwin":
		app.destination = fmt.Sprintf("%s/Library/Application Support/minecraft/mods", user.HomeDir)
	case "linux":
		app.destination = fmt.Sprintf("%s/Library/Application Support/minecraft/mods", user.HomeDir)
	default:
		if err := fmt.Errorf("%s does not have a default directory, please provide target directory", os); err != nil {
			return err
		}
	}
	return nil
}

func (app *appEnv) EnsureDestination() error {
	logrus.Debugf("Making Folder if not exist: %s", app.destination)
	err := os.MkdirAll(app.destination, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (app *appEnv) PrepareDestinationFolder() error {
	app.GetTargetDirectory()
	logrus.Debugf("Mod Destination is set to: %s", app.destination)
	if app.clearMods {
		logrus.Debugf("Removing contents of: %s", app.destination)
		if err := os.RemoveAll(app.destination); err != nil {
			return err
		}
	}
	app.EnsureDestination()
	return nil
}

func (app *appEnv) FetchForgeAPIJSON(url string, data interface{}) error {
	logrus.Debugf("Fetching: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Accept":    []string{"application/json"},
		"x-api-key": []string{app.forgeKey},
	}
	resp, err := app.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) FetchJSON(url string, data interface{}) error {
	logrus.Debugf("Fetching JSON: %s", url)
	resp, err := app.hc.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) FetchXML(url string, data interface{}) error {
	logrus.Debugf("Fetching XML: %s", url)
	resp, err := app.hc.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(respBody, &data)
}

func (app *appEnv) LoadModsFromJSON() error {
	if app.jsonFile == "" {
		return nil
	}
	logrus.Debugf("Fetching JSON from file: %s", app.jsonFile)
	jsonFile, err := os.Open(app.jsonFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result JSONMods
	json.Unmarshal([]byte(byteValue), &result)
	logrus.Debugf("Pulled from json file: %s", result)
	app.modsFromJSON = result
	return nil
}

func (app *appEnv) FetchAndSave(url string, fileName string, destPath string) error {
	logrus.Infof("Downloading: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Accept":    []string{"application/json"},
		"x-api-key": []string{app.forgeKey},
	}
	resp, err := app.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(fmt.Sprintf(destPath + "/" + fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func (app *appEnv) PrintDestinationFiles() {
	files, err := os.ReadDir(app.destination)
	if err != nil {
		log.Fatal(err)
	}

	logrus.Infof("Files in Destination Folder:")
	for _, file := range files {
		logrus.Infof("  %s  ", file.Name())
	}
}
