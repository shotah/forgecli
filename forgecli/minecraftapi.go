package forgecli

import "fmt"

// MinecraftVersionURL URL to grab the json manifest
const MinecraftVersionURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

// MCLatest Latest release
type MCLatest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

// MCVersion Specific Mojang versions
type MCVersion struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

// MCVersionResponse Default Mojang Version Response
type MCVersionResponse struct {
	Latest   MCLatest    `json:"latest"`
	Versions []MCVersion `json:"versions"`
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
