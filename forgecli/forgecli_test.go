package forgecli

import (
	"fmt"
	"net/http"
	"os"
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

// Test will fail as MC receives version updates
func TestGetVersionTypeNumber(t *testing.T) {
	LoadDotEnv()
	var app appEnv
	app.version = "1.18.1"
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	expected := 73250
	app.GetVersionTypeNumber()
	if app.forgeGameVersionType != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, app.forgeGameVersionType)
	}
}

func TestGetModsByProjectIDsWithFamily(t *testing.T) {
	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	app.modsToDownload = make(map[int]ForgeMod)
	app.modReleaseType = 2 // Beta
	app.version = "1.18.1"
	// With only fabric mods:
	app.projectIDs = "416089,391366,552655"
	app.modfamily = "fabric"
	expected := 3
	app.GetVersionTypeNumber()
	app.GetModsByProjectIDs()
	if len(app.modsToDownload) != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, app.forgeGameVersionType)
	}
}

func TestGetModsByProjectIDsWithOuFamily(t *testing.T) {
	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	app.modsToDownload = make(map[int]ForgeMod)
	app.modReleaseType = 2 // Beta
	app.version = "1.18.1"

	// With non-fabric mods:
	app.projectIDs = "306612,416089,220318"
	expected := 3
	app.GetVersionTypeNumber()
	app.GetModsByProjectIDs()
	if len(app.modsToDownload) != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, app.forgeGameVersionType)
	}
}

func TestFetchforgeAPIJSON(t *testing.T) {
	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	var resp ForgeMods
	url := "https://api.curseforge.com/v1/mods/306612/files?gameVersionTypeID=73250&index=0&pageSize=3"
	if err := app.FetchForgeAPIJSON(url, &resp); err != nil {
		t.Errorf("Test failed, by throwing error")
	}
	fmt.Println(resp.Data[0].DisplayName)
}

func TestGetMCVersionNoInput(t *testing.T) {
	var app appEnv
	app.hc = *http.DefaultClient
	app.version = ""
	expected := "1.18.2"
	app.GetMCVersion()
	if app.version != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.version)
	}
}

func TestGetMCVersionWithInput(t *testing.T) {
	var app appEnv
	app.hc = *http.DefaultClient
	app.version = "1.17.1"
	expected := "1.17.1"
	app.GetMCVersion()
	if app.version != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.version)
	}
}

func TestGetMCVersionWithBadInput(t *testing.T) {
	var app appEnv
	app.hc = *http.DefaultClient
	app.version = "1.17.0"
	expected := "nil"
	app.GetMCVersion()
	if app.version != "" {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.version)
	}
}
