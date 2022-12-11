package forgecli

// FabricMetadataURL is the url for the installers meta data:
const FabricMetadataURL = "https://maven.fabricmc.net/net/fabricmc/fabric-installer/maven-metadata.xml"

// Fabric version struct
type fabricVersion struct {
	Version string `xml:"version"`
}

// fabricVersioning struct which contains most of the data about the versions
type fabricVersioning struct {
	Latest      string          `xml:"latest"`
	Release     string          `xml:"release"`
	Versions    []fabricVersion `xml:"versions"`
	LastUpdated string          `xml:"lastUpdated"`
}

// XMLFabric metadata definition struct
type XMLFabric struct {
	GroupID    string           `xml:"groupId"`
	ArtifactID string           `xml:"artifactId"`
	Versioning fabricVersioning `xml:"versioning"`
}
