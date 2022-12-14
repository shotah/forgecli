package main

import (
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
	"github.com/shotah/forgecli/forgecli"
)

func MockCurseForgeVersions(t *testing.T) {
	mockVersionTypes := "./forgecli/mocks/mc_version_types.json"
	body, err := os.ReadFile(mockVersionTypes)
	st.Expect(t, err, nil)
	gock.New("https://api.curseforge.com").
		Get("/v1/games/432/version-types").
		Reply(200).
		JSON(body)
}

func MockCurseForgeModResponse(t *testing.T) {
	mockVersionTypes := "./forgecli/mocks/voice_mod_response.json"
	body, err := os.ReadFile(mockVersionTypes)
	st.Expect(t, err, nil)
	gock.New("https://api.curseforge.com").
		Get("/v1/mods/416089/files?gameVersionTypeID=73250&index=0&pageSize=999").
		Reply(200).
		JSON(body)
}

func TestHelp(t *testing.T) {
	defer gock.Off()
	MockCurseForgeVersions(t)
	MockCurseForgeModResponse(t)

	expected := 2
	os.Args = []string{"-help"}
	actual := forgecli.CLI(os.Args[1:])
	if actual != expected {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestMain(m *testing.M) {
	os.Args = []string{"-help"}
	exitVal := m.Run()
	os.Exit(exitVal)
}
