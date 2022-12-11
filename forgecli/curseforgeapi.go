// Package forgecli is main cli
package forgecli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func (app *appEnv) SetForgeAPIKey() error {
	if app.forgeKey != "" {
		logrus.Debug("forgeKey Provided")
		return nil
	}
	if os.Getenv("FORGEKEY") != "" {
		app.forgeKey = os.Getenv("FORGEKEY")
		logrus.Debug("forgeKey found in FORGEKEY")
		return nil
	}
	if os.Getenv("MODS_FORGEAPI_KEY") != "" {
		app.forgeKey = os.Getenv("MODS_FORGEAPI_KEY")
		logrus.Debug("forgeKey found in MODS_FORGEAPI_KEY")
		return nil
	}
	return fmt.Errorf("MISSING a required field: forgeKey")
}

func (app *appEnv) GetModsByProjectIDs() {
	if app.projectIDs == "" {
		return
	}
	projectIDs := strings.Split(app.projectIDs, ",")
	for _, projectID := range projectIDs {
		var convertedJSONMod JSONMod
		convertedJSONMod.ProjectID = projectID
		logrus.Debugf("Getting Mod: %s", convertedJSONMod.ProjectID)
		err := app.GetModsFromForge(convertedJSONMod, app.modReleaseType)
		check(err)
	}
}

func (app *appEnv) GetModsByJSONFile() {
	if app.jsonFile == "" {
		return
	}
	logrus.Debugf("Getting Mod: %s", app.modsFromJSON)
	var releaseType ReleaseType
	for _, fileMod := range app.modsFromJSON {
		logrus.Debugf("Getting Mod: %s", fileMod.ProjectID)
		if fileMod.ReleaseType != "" {
			releaseType = releaseLookup[fileMod.ReleaseType]
		} else {
			releaseType = app.modReleaseType
		}
		err := app.GetModsFromForge(fileMod, releaseType)
		check(err)
	}
}

func (app *appEnv) GetModsDependencies() error {
	if !app.downloadDependencies {
		return nil
	}
	for _, mod := range app.modsToDownload {
		for _, modDep := range mod.Dependencies {
			if modDep.RelationType == 3 {
				var convertedJSONMod JSONMod
				convertedJSONMod.ProjectID = strconv.Itoa(modDep.ModID)
				logrus.Debugf("Getting Mod dependency: %s", convertedJSONMod.ProjectID)
				err := app.GetModsFromForge(convertedJSONMod, releaseLookup["release"])
				check(err)
			}
		}
	}
	return nil
}

func (app *appEnv) GetModsFromForge(modToGet JSONMod, _ ReleaseType) error {
	var resp ForgeMods
	pageIndex := 0
	pageSize := 999
	url := fmt.Sprintf(
		"https://api.curseforge.com/v1/mods/%s/files?gameVersionTypeID=%d&index=%d&pageSize=%d",
		modToGet.ProjectID, app.forgeGameVersionType, pageIndex, pageSize,
	)
	err := app.FetchForgeAPIJSON(url, &resp)
	check(err)

	foundID := 0
	var foundMod ForgeMod
	for _, currMod := range resp.Data {
		if currMod.ID > foundID && app.ModFilter(currMod, modToGet) {
			foundID = currMod.ID
			foundMod = currMod
		}
	}
	if foundID == 0 {
		return fmt.Errorf("could not find %s for minecraft version: %s or family: %s", modToGet.ProjectID, app.version, app.modfamily)
	}
	app.modsToDownload[foundID] = foundMod
	logrus.Infof("Found Latest FileID: %d for Mod: %s", foundID, modToGet.ProjectID)
	return nil
}

func (app *appEnv) ModFamilyFilter(currMod ForgeMod) bool {
	// Apply mod family filter
	if app.modfamily != "" {
		result := contains(currMod.GameVersions, string(app.modfamily))
		logrus.Debugf("Mod's family Filter Result: %t", result)
		return result
	}
	return true
}

func (app *appEnv) ModVersionFilter(currMod ForgeMod, modToGet JSONMod) bool {
	// Apply Version filter
	modVersion := app.version
	if modToGet.Version != "" {
		modVersion = modToGet.Version
	}
	logrus.Debugf("Mod's Game Version Filters: \n%s \n%s", currMod.GameVersions, string(modVersion))
	result := contains(currMod.GameVersions, string(modVersion))
	logrus.Debugf("Mod's Game Version Filter Result: %t", result)
	return result
}

func (app *appEnv) ModFilter(currMod ForgeMod, modToGet JSONMod) bool {
	// Filtering on Filename Only
	if modToGet.Filename != "" {
		return false
	}
	return app.ModFamilyFilter(currMod) && app.ModVersionFilter(currMod, modToGet)
}

func (app *appEnv) GetVersionTypeNumber() error {
	// Forge has a specific format to validate Minecraft 1.17
	shortNumber := strings.Join(strings.Split(app.version, ".")[:2], ".")
	forgeVersionName := "Minecraft " + shortNumber
	logrus.Debugf("Fetching VersionType for MC version: %s, %s", forgeVersionName, ForgeVersionTypeURL)

	var resp ForgeVersions
	if err := app.FetchForgeAPIJSON(ForgeVersionTypeURL, &resp); err != nil {
		return err
	}
	for _, v := range resp.Data {
		if v.Name == forgeVersionName {
			app.forgeGameVersionType = v.ID
			logrus.Debugf("found VersionType: %d", v.ID)
			return nil
		}
	}
	app.forgeGameVersionType = 0
	return fmt.Errorf("could not find forge version for:%s", forgeVersionName)
}

func (app *appEnv) DownloadMods() error {
	for _, mod := range app.modsToDownload {
		if err := app.FetchAndSave(mod.DownloadURL, mod.Filename, app.destination); err != nil {
			return err
		}
	}
	return nil
}
