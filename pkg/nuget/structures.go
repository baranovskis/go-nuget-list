package nuget

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

type ResponseResources struct {
	Version   string     `json:"version"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Id   string `json:"@id"`
	Type string `json:"@type"`
}

type ResponseQuery struct {
	XMLName xml.Name      `xml:"feed"`
	Data    []PackageData `json:"data" xml:"entry"`
}

type PackageAuthors struct {
	Values []string
}

func (pa *PackageAuthors) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	pa.Values = strings.Split(value, ",")

	for i := range pa.Values {
		pa.Values[i] = strings.TrimSpace(pa.Values[i])
	}

	return nil
}

func (pa *PackageAuthors) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &pa.Values)
}

type PackageTags struct {
	Values []string
}

func (pt *PackageTags) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	pt.Values = strings.Split(value, " ")
	return nil
}

func (pt *PackageTags) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &pt.Values)
}

type PackageData struct {
	Id          string         `json:"id" xml:"properties>Id"`
	Version     string         `json:"version" xml:"properties>NormalizedVersion"`
	Description string         `json:"description" xml:"properties>Description"`
	Title       string         `json:"title" xml:"title"`
	Summary     string         `json:"summary" xml:"summary"`
	Authors     PackageAuthors `json:"authors" xml:"author>name"`
	Tags        PackageTags    `json:"tags" xml:"properties>Tags"`
	LicenseUrl  string         `json:"licenseUrl" xml:"properties>LicenseUrl"`
	ProjectUrl  string         `json:"projectUrl" xml:"properties>ProjectUrl"`
}

type PackageSource struct {
	FileName        string
	SourceName      string
	Path            string
	ProtocolVersion string
}
