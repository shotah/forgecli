package forgecli

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
