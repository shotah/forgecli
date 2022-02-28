package forgecli

// MinecraftGameID Default Minecraft Game ID for Forge
const MinecraftGameID = "432"

// ForgeAPIBaseURL Base Forge API URL
const ForgeAPIBaseURL = "https://api.curseforge.com/v1"

// ForgeVersionTypeURL Default version types URL
const ForgeVersionTypeURL = ForgeAPIBaseURL + "/games/" + MinecraftGameID + "/version-types"

// ReleaseType int declaration
type ReleaseType int

const (
	// Release used to fetch mods that are marked as release
	Release ReleaseType = 1
	// Beta used to fetch mods that are marked as Beta
	Beta ReleaseType = 2
	// Alpha used to fetch mods that are marked as Alpha
	Alpha ReleaseType = 3
)

// releaseLookup used when accepting strings and converting to ints
var releaseLookup = map[string]ReleaseType{
	"release": Release,
	"beta":    Beta,
	"alpha":   Alpha,
}

// ForgePagination was to be used with Pagination, currently not used
type ForgePagination struct {
	Index       int `json:"index"`
	PageSize    int `json:"pageSize"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}

// ForgeVersion Individual version Struct
type ForgeVersion struct {
	ID     int    `json:"id"`
	GameID int    `json:"gameID"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

// ForgeVersions main Data loop for ForgeVersion
type ForgeVersions struct {
	Data []ForgeVersion `json:"data"`
}

// ForgeDependency part of the forge mod declaration, used when we download dependencies
type ForgeDependency struct {
	ModID        int `json:"modID"`
	FileID       int `json:"fileID"`
	RelationType int `json:"relationType"`
}

// ForgeMod All relative define fields that are supplied by Forge per mod
type ForgeMod struct {
	ID                   int               `json:"id"`
	GameID               int               `json:"gameID"`
	ModID                int               `json:"modID"`
	IsAvailable          bool              `json:"isAvailable"`
	DisplayName          string            `json:"displayName"`
	Filename             string            `json:"filename"`
	ReleaseType          ReleaseType       `json:"releaseType"`
	FileStatus           int               `json:"fileStatus"`
	Hashes               []interface{}     `json:"Hashes"`
	FileDate             string            `json:"fileDate"`
	FileLength           int               `json:"fileLength"`
	DownloadCount        int               `json:"downloadCount"`
	DownloadURL          string            `json:"downloadURL"`
	GameVersions         []string          `json:"gameVersions"`
	SortableGameVersions []interface{}     `json:"sortableGameVersions"`
	Dependencies         []ForgeDependency `json:"dependencies"`
	AlternateFileID      int               `json:"alternateFileID"`
	IsServerPack         bool              `json:"isServerPack"`
	FileFingerprint      int               `json:"fileFingerprint"`
	Modules              []interface{}     `json:"modules"`
}

// ForgeMods Main Data array entrypoint
type ForgeMods struct {
	Data       []ForgeMod      `json:"data"`
	Pagination ForgePagination `json:"pagination"`
}
