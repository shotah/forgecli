package forgecli

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type appEnv struct {
	hc                      http.Client
	jsonFile                string
	forgeKey                string
	version                 string
	clientInstaller         bool
	projectIDs              string
	downloadDependencies    bool
	clearMods               bool
	modfamily               FamilyType
	modReleaseType          ReleaseType
	destination             string
	modsFromJSON            JSONMods
	forgeGameVersionType    int
	modsToDownload          map[int]ForgeMod
	clientInstallerVersion  string
	clientInstallerFileName string
	isDebug                 bool
}

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
	logrus.Info("Mods Download Complete.")

	app.ClientInstaller()
	return nil
}

var javaVersion, javaErr = exec.Command("java", "-version").CombinedOutput()

func (app *appEnv) ValidateJavaInstallation() error {
	if javaErr != nil {
		logrus.Debug("java version not found")
		return fmt.Errorf("unable to find java, please install and try again")
	}
	logrus.Debugf("java version found: %s", javaVersion)
	return nil
}

// ClientInstaller main handler for the client installations
func (app *appEnv) ClientInstaller() error {
	if !app.clientInstaller {
		return nil
	}
	logrus.Info("Starting Client Installer")
	// Makes sure their is an instance of Java installed to use for the installation.
	if err := app.ValidateJavaInstallation(); err != nil {
		return err
	}
	// Starts the fabric installation process.
	switch app.modfamily {
	case Fabric:
		if err := app.FabricClientInstaller(); err != nil {
			return err
		}
		// case Forge:
		// 	if err := app.FabricClientInstaller(); err != nil {
		// 		return err
		// 	}
	}

	logrus.Info("Finishing Client Installer")
	return nil
}
