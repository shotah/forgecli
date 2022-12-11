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

type forgecliError interface {
	Error() string
}

func check(e forgecliError) {
	if e != nil {
		logrus.Error(e.Error())
		logrus.Error("Exiting...")
		os.Exit(1)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

func (app *appEnv) GetTargetDirectory() {
	if app.destination != "" {
		return
	}

	user, err := user.Current()
	check(err)
	os := runtime.GOOS
	switch os {
	case "windows":
		app.destination = fmt.Sprintf("%s\\AppData\\Roaming\\.minecraft\\mods", user.HomeDir)
	case "darwin":
		app.destination = fmt.Sprintf("%s/Library/Application Support/minecraft/mods", user.HomeDir)
	case "linux":
		app.destination = fmt.Sprintf("%s/Library/Application Support/minecraft/mods", user.HomeDir)
	default:
		err := fmt.Errorf("%s does not have a default directory, please provide target directory", os)
		check(err)
	}
}

func (app *appEnv) EnsureDestination() {
	logrus.Debugf("Making Folder if not exist: %s", app.destination)
	err := os.MkdirAll(app.destination, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		check(err)
	}
}

func (app *appEnv) PrepareDestinationFolder() {
	app.GetTargetDirectory()
	logrus.Debugf("Mod Destination is set to: %s", app.destination)
	if app.clearMods {
		logrus.Debugf("Removing contents of: %s", app.destination)
		err := os.RemoveAll(app.destination)
		check(err)
	}
	app.EnsureDestination()
}

func (app *appEnv) FetchForgeAPIJSON(url string, data interface{}) forgecliError {
	logrus.Debugf("Fetching: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header = http.Header{
		"Accept":    []string{"application/json"},
		"x-api-key": []string{app.forgeKey},
	}
	resp, err := app.hc.Do(req)
	check(err)
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) FetchJSON(url string, data interface{}) forgecliError {
	logrus.Debugf("Fetching JSON: %s", url)
	resp, err := app.hc.Get(url)
	check(err)
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) FetchXML(url string, data interface{}) forgecliError {
	logrus.Debugf("Fetching XML: %s", url)
	resp, err := app.hc.Get(url)
	check(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	check(err)
	return xml.Unmarshal(respBody, &data)
}

func (app *appEnv) LoadModsFromJSON() {
	if app.jsonFile == "" {
		return
	}
	logrus.Debugf("Fetching JSON from file: %s", app.jsonFile)
	jsonFile, err := os.Open(app.jsonFile)
	check(err)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result JSONMods
	json.Unmarshal([]byte(byteValue), &result)
	logrus.Debugf("Pulled from json file: %s", result)
	app.modsFromJSON = result
}

func (app *appEnv) FetchAndSave(url string, fileName string, destPath string) forgecliError {
	logrus.Infof("Downloading: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header = http.Header{
		"Accept":    []string{"application/json"},
		"x-api-key": []string{app.forgeKey},
	}
	resp, err := app.hc.Do(req)
	check(err)
	defer resp.Body.Close()

	f, err := os.Create(fmt.Sprintf(destPath + "/" + fileName))
	check(err)
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
