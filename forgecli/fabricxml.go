package forgecli

type fabricVersion struct {
	Version string `xml:"version"`
}

type fabricVersions struct {
	Versions []fabricVersion `xml:"versions"`
}

type fabricVersioning struct {
	Latest      string         `xml:"latest"`
	Release     string         `xml:"release"`
	Versions    fabricVersions `xml:"versions"`
	LastUpdated string         `xml:"lastUpdated"`
}

type fabricMetadata struct {
	GroupId    string           `xml:"groupId"`
	ArtifactId string           `xml:"artifactId"`
	Versioning fabricVersioning `xml:"versioning"`
}

// XMLFabric is a simple Fabric release matrix
var XMLFabric struct {
	Metadata fabricMetadata `xml:"metadata"`
}
