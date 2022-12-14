package forgecli

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
	"github.com/sirupsen/logrus"
)

func TestSetForgeAPIKeyNoKey(t *testing.T) {
	os.Setenv("FORGEKEY", "")
	os.Setenv("MODS_FORGEAPI_KEY", "")

	var app appEnv
	app.forgeKey = ""
	err := app.SetForgeAPIKey()
	expectedErrorMessage := fmt.Errorf("MISSING a required field: forgeKey")
	st.Expect(t, err, expectedErrorMessage)
	expected := ""
	received := app.forgeKey
	if received != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, received)
	}
}

func TestSetForgeAPIKeyProvidedKey(t *testing.T) {
	os.Setenv("FORGEKEY", "")
	os.Setenv("MODS_FORGEAPI_KEY", "")

	var app appEnv
	app.forgeKey = "testMockKey"
	err := app.SetForgeAPIKey()
	st.Expect(t, err, nil)
	expected := "testMockKey"
	received := app.forgeKey
	if received != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, received)
	}
}

func TestSetForgeAPIKeyEnvForgeKey(t *testing.T) {
	expected := "testForgeMockKey"
	os.Setenv("FORGEKEY", expected)
	os.Setenv("MODS_FORGEAPI_KEY", "")

	var app appEnv
	err := app.SetForgeAPIKey()
	st.Expect(t, err, nil)
	received := app.forgeKey
	if received != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, received)
	}
}

func TestSetForgeAPIKeyEnvModsForgeKey(t *testing.T) {
	expected := "testModsMockKey"
	os.Setenv("FORGEKEY", "")
	os.Setenv("MODS_FORGEAPI_KEY", expected)

	var app appEnv
	err := app.SetForgeAPIKey()
	st.Expect(t, err, nil)
	received := app.forgeKey
	if received != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, received)
	}
}

func TestGetModsByJSONFileNoFileProvided(t *testing.T) {
	var app appEnv
	err := app.GetModsByJSONFile()
	st.Expect(t, err, nil)
}

func TestGetModsByJSONFileFileProvided(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeVersions(t)
	MockCurseForgeModResponse(t)

	var app appEnv
	app.modsToDownload = make(map[int]ForgeMod)
	app.version = "1.19.2"
	app.modfamily = Fabric
	app.jsonFile = "./mocks/modsfile_by_json.json"
	err := app.LoadModsFromJSON()
	st.Expect(t, err, nil)
	err = app.GetModsByJSONFile()
	st.Expect(t, err, nil)
}

// Test will fail as MC receives version updates
func TestGetVersionTypeNumber(t *testing.T) {
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeVersions(t)

	LoadDotEnv()
	var app appEnv
	app.version = "1.18.1"
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	expected := 73250
	err := app.GetVersionTypeNumber()
	st.Expect(t, err, nil)
	if app.forgeGameVersionType != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, app.forgeGameVersionType)
	}
}

func TestGetModsByProjectIDsWithFamily(t *testing.T) {
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeVersions(t)
	MockCurseForgeModResponse(t)

	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	// app.forgeKey = os.Getenv("FORGEKEY")
	app.modsToDownload = make(map[int]ForgeMod)
	app.modReleaseType = 2 // Beta
	app.version = "1.19.3"

	// setup mc versions with version type checker
	err := app.GetVersionTypeNumber()
	st.Expect(t, err, nil)

	// With only fabric voice mod
	app.projectIDs = "416089"
	app.modfamily = "fabric"
	err = app.GetModsByProjectIDs()
	st.Expect(t, err, nil)

	// Validate we got the one mod we expected
	received := len(app.modsToDownload)
	expected := 1
	if received != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, received)
	}
}

func TestGetModsByProjectIDsWithOuFamily(t *testing.T) {
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeVersions(t)
	MockCurseForgeModResponse(t)

	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	app.modsToDownload = make(map[int]ForgeMod)
	app.modReleaseType = 2 // Beta
	app.version = "1.19.1"

	// With non-fabric mods
	app.projectIDs = "416089"
	expected := 1
	app.GetVersionTypeNumber()
	app.GetModsByProjectIDs()
	received := len(app.modsToDownload)
	if received != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, received)
	}
}

func TestFetchforgeAPIJSON(t *testing.T) {
	defer gock.Off()
	MockMCVersions(t)
	MockCurseForgeModResponse(t)

	LoadDotEnv()
	var app appEnv
	app.hc = *http.DefaultClient
	app.forgeKey = os.Getenv("FORGEKEY")
	var resp ForgeMods
	url := "https://api.curseforge.com/v1/mods/416089/files?gameVersionTypeID=73250&index=0&pageSize=999"
	if err := app.FetchForgeAPIJSON(url, &resp); err != nil {
		t.Errorf("Test failed, by throwing error")
	}
	fmt.Println(resp.Data[0].DisplayName)
}

func TestGetMCVersionNoInput(t *testing.T) {
	MockMCVersions(t)

	var app appEnv
	app.hc = *http.DefaultClient
	app.version = ""
	expected := "1.19.3"
	app.GetMCVersion()
	if app.version != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.version)
	}
}

func TestGetMCVersionWithInput(t *testing.T) {
	MockMCVersions(t)

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
	MockMCVersions(t)

	var app appEnv
	app.hc = *http.DefaultClient
	app.version = "1.17.0"
	expected := "nil"
	app.GetMCVersion()
	if app.version != "" {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.version)
	}
}
