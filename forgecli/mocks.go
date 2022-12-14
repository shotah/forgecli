package forgecli

import (
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

// MockFabricXML used to return mock xml data for fabric calls
func MockFabricXML(t *testing.T) {
	mockFile := "./mocks/fabric.xml"
	gock.New("https://maven.fabricmc.net").
		Get("/net/fabricmc/fabric-installer/maven-metadata.xml").
		Reply(200).File(mockFile)
}

// MockFabricJAR returns a mock jar in place of downloading the real jar file.
func MockFabricJAR(t *testing.T) {
	mockFile := "./mocks/fake.jar"
	gock.New("https://maven.fabricmc.net").
		Get("/net/fabricmc/fabric-installer/0.11.1/fabric-installer-0.11.1.jar").
		Reply(200).File(mockFile)
}

// MockMCVersions returns a manifest instead of calling the MC api for MC versions
func MockMCVersions(t *testing.T) {
	mockFile := "./mocks/mc_version_manifest.json"
	body, err := os.ReadFile(mockFile)
	st.Expect(t, err, nil)
	gock.New("https://launchermeta.mojang.com").
		Get("/mc/game/version_manifest.json").
		Reply(200).
		JSON(body)
}

// MockCurseForgeVersions returns versions types instead of calling the curse forge for the version types
func MockCurseForgeVersions(t *testing.T) {
	mockFile := "./mocks/mc_version_types.json"
	body, err := os.ReadFile(mockFile)
	st.Expect(t, err, nil)
	gock.New("https://api.curseforge.com").
		Get("/v1/games/432/version-types").
		Reply(200).
		JSON(body)
}

// MockCurseForgeModResponse returns a mock module json, instead of calling the real api for a mod json
func MockCurseForgeModResponse(t *testing.T) {
	mockFile := "./mocks/voice_mod_response.json"
	body, err := os.ReadFile(mockFile)
	st.Expect(t, err, nil)
	gock.New("https://api.curseforge.com").
		Get("/v1/mods/416089/files").
		Reply(200).
		JSON(body)
}
