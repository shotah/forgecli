package forgecli

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// CLI Main Module Entrypoint
func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

type appEnv struct {
	hc                   http.Client
	jsonFile             string
	forgeKey             string
	version              string
	clientInstaller      bool
	projectIDs           string
	downloadDependencies bool
	clearMods            bool
	modfamily            FamilyType
	modReleaseType       ReleaseType
	destination          string
	modsFromJSON         JSONMods
	forgeGameVersionType int
	modsToDownload       map[int]ForgeMod
	isDebug              bool
}

func (app *appEnv) fromArgs(args []string) error {
	// Shallow copy of default client
	app.hc = *http.DefaultClient
	app.modsToDownload = make(map[int]ForgeMod)

	fl := flag.NewFlagSet("forgecli", flag.ContinueOnError)
	fl.StringVar(&app.jsonFile, "file", "", "file json with required mods and settings")
	fl.StringVar(&app.forgeKey, "forgekey", "", "ForgeAPIKey used in Authentication with the Forge API")
	fl.StringVar(&app.destination, "destination", "", "destination directory for mods")
	fl.StringVar(&app.version, "version", "", "Minecraft version you are installing")
	fl.BoolVar(&app.clientInstaller, "client", false, "Downloads and installs Client based on Family (if no family, no client install will be done)")
	fl.BoolVar(&app.downloadDependencies, "dependencies", true, "Download Mods Dependencies")
	fl.BoolVar(&app.clearMods, "clear", false, "Clear Mods from destination (mods folder)")
	fl.BoolVar(&app.isDebug, "debug", false, "enable debug logging")
	fl.StringVar(&app.projectIDs, "projects", "", "Forge Project IDs separated by commas 12345,67890")
	inputReleaseType := fl.String("release", "release", "Mods release type, release, beta, alpha")
	inputFamily := fl.String("family", "", "Minecraft type: Fabric, Forge, Bukkit")

	// Parsing the Args before they can be used
	if err := fl.Parse(args); err != nil {
		return err
	}

	// Type Conversions
	app.modReleaseType = releaseLookup[*inputReleaseType]
	if len(*inputFamily) > 0 {
		app.modfamily = familyTypeLookup[*inputFamily]
	}

	// Setting up logrus
	if app.isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Checking for Required Fields
	if app.projectIDs == "" && app.jsonFile == "" {
		fmt.Fprintf(os.Stderr, "Did not receive Project IDs to Download.\n")
		fl.Usage()
		return flag.ErrHelp
	}

	return nil
}

func (app *appEnv) run() error {
	logrus.Info("Starting Forge Mod Lookup")
	app.SetForgeAPIKey()

	app.GetMCVersion()
	logrus.Debugf("Using Minecraft Version: %s", app.version)

	app.GetVersionTypeNumber()
	logrus.Debugf("Using Forge VersionType: %d", app.forgeGameVersionType)

	app.GetModsByProjectIDs()
	app.LoadModsFromJSON()
	app.GetModsByJSONFile()

	app.GetModsDependencies()

	app.PrepareDestinationFolder()
	app.DownloadMods()
	app.PrintDestinationFiles()
	logrus.Info("Download Complete.")
	return nil
}

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

func (app *appEnv) ModFilter(currMod ForgeMod, modToGet JSONMod) bool {
	var results []bool
	// Filtering on Filename Only
	if modToGet.Filename != "" {
		return strings.EqualFold(currMod.Filename, modToGet.Filename)
	}

	// Apply mod family filter
	if app.modfamily != "" {
		result := contains(currMod.GameVersions, string(app.modfamily))
		results = append(results, result)
	}

	// Apply Version filter
	if modToGet.Version != "" {
		result := containsPrefix(currMod.GameVersions, string(modToGet.Version))
		results = append(results, result)
	} else {
		result := containsPrefix(currMod.GameVersions, string(app.version))
		results = append(results, result)
	}
	return allTrue(results)
}

func (app *appEnv) GetMCVersion() error {
	var resp MCVersionResponse
	if err := app.FetchJSON(MinecraftVersionURL, &resp); err != nil {
		return fmt.Errorf("could not get minecraft version from:\n%s", MinecraftVersionURL)
	}
	if app.version == "" {
		version := resp.Latest.Release
		app.version = version
		return nil
	}
	inputVersion := app.version
	for _, v := range resp.Versions {
		if v.ID == inputVersion {
			returnVersion := v.ID
			app.version = returnVersion
			return nil
		}
	}
	app.version = ""
	return fmt.Errorf("could not find minecraft version: %s", inputVersion)
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
		if err := app.FetchAndSave(mod.DownloadURL, mod.Filename); err != nil {
			return err
		}
	}
	return nil
}

func (app *appEnv) FabricClientInstaller() error {
	var fabricXMLResponse XMLFabric
	if err := app.FetchXML(FabricMetadataURL, &fabricXMLResponse); err != nil {
		return fmt.Errorf("could not get fabric version from:\n%s", MinecraftVersionURL)
	}
	if app.version == "" {
		// Install command
		// java -jar './fabric-installer-0.11.1.jar' client
		logrus.Debugf("java -jar './fabric-installer-%s.jar' client", fabricXMLResponse.Versioning.Latest)
	} else {
		// Install command.. etc...
		// java -jar './fabric-installer-0.11.1.jar' client -mcversion 1.18.1
		// java -jar './fabric-installer-0.11.1.jar' client -mcversion app.version
		logrus.Debugf("java -jar './fabric-installer-%s.jar' client -mcversion %s", fabricXMLResponse.Versioning.Latest, app.version)
	}
	return nil
}

func (app *appEnv) ClientInstaller() error {
	if app.modfamily == Fabric {
		if err := app.FabricClientInstaller(); err != nil {
			return err
		}
	}
	return nil
}
