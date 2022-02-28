package forgecli

// JSONMod is a simple mod definition structure
type JSONMod struct {
	Name        string `json:"name"`
	ProjectID   string `json:"projectID"`
	ReleaseType string `json:"releaseType"`
	Filename    string `json:"filename"`
}

// JSONMods is just a default array of mods
type JSONMods []JSONMod
