package forgecli

// ForgeInstallerURL string formatted url to download the installer jar
const ForgeInstallerURL = "https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-installer.jar"

// ForgeMetadataURL is the url for the installers meta data:
const ForgeMetadataURL = "https://maven.minecraftforge.net/net/minecraftforge/forge/maven-metadata.xml"

// Forge version struct
type forgeVersion struct {
	Version string `xml:"version"`
}

// forgeVersioning struct which contains most of the data about the versions
type forgeVersioning struct {
	Latest      string         `xml:"latest"`
	Release     string         `xml:"release"`
	Versions    []forgeVersion `xml:"versions"`
	LastUpdated string         `xml:"lastUpdated"`
}

// XMLForge metadata definition struct
type XMLForge struct {
	GroupID    string          `xml:"groupId"`
	ArtifactID string          `xml:"artifactId"`
	Versioning forgeVersioning `xml:"versioning"`
}
