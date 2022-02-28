package forgecli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		panic(e)
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

func (app *appEnv) osTargetDirectory() {
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
		log.Fatalf(fmt.Sprintf("%s does not have a default directory, please provide target directory", os))
	}
}

func (app *appEnv) ensureDestination() {
	logrus.Debugf("Making Folder if not exist: %s", app.destination)
	err := os.MkdirAll(app.destination, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		check(err)
	}
}

func (app *appEnv) PrepareDestinationFolder() {
	app.osTargetDirectory()
	logrus.Debugf("Mod Destination is set to: %s", app.destination)
	if app.clearMods {
		logrus.Debugf("Removing contents of: %s", app.destination)
		err := os.RemoveAll(app.destination)
		check(err)
	}
	app.ensureDestination()
}

func (app *appEnv) fetchforgeAPIJSON(url string, data interface{}) error {
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

func (app *appEnv) fetchJSON(url string, data interface{}) error {
	logrus.Debugf("Fetching JSON: %s", url)
	resp, err := app.hc.Get(url)
	check(err)
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (app *appEnv) LoadModsFromJSON() {
	if app.jsonFile == "" {
		return
	}
	logrus.Debugf("Fetching JSON from file: %s", app.jsonFile)
	jsonFile, err := os.Open(app.jsonFile)
	check(err)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result JSONMods
	json.Unmarshal([]byte(byteValue), &result)
	logrus.Debugf("Pulled from json file: %s", result)
	app.modsFromJSON = result
}

func (app *appEnv) fetchAndSave(url, destPath string) error {
	logrus.Debugf("Fetching and Saving: %s", url)
	resp, err := app.hc.Get(url)
	check(err)
	defer resp.Body.Close()

	f, err := os.Create(fmt.Sprintf(app.destination + "/" + destPath))
	check(err)

	_, err = io.Copy(f, resp.Body)
	return err
}
